package bitbucketserver

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/mrjones/oauth"
	"github.com/pkg/errors"
	"github.com/rancher/norman/httperror"
	"github.com/rancher/rancher/pkg/pipeline/remote/model"
	"github.com/rancher/rancher/pkg/pipeline/utils"
	"github.com/rancher/rancher/pkg/ref"
	"github.com/rancher/rancher/pkg/settings"
	"github.com/rancher/types/apis/project.cattle.io/v3"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	maxPerPage        = "100"
	requestTokenURL   = "%s/plugins/servlet/oauth/request-token"
	authorizeTokenURL = "%s/plugins/servlet/oauth/authorize"
	accessTokenURL    = "%s/plugins/servlet/oauth/access-token"
)

type client struct {
	BaseURL     string
	ConsumerKey string
	PrivateKey  string
	RedirectURL string
}

func New(config *v3.BitbucketServerPipelineConfig) (model.Remote, error) {
	if config == nil {
		return nil, errors.New("empty bitbucket server config")
	}
	bsClient := &client{
		ConsumerKey: config.ConsumerKey,
		PrivateKey:  config.PrivateKey,
		RedirectURL: config.RedirectURL,
	}
	if config.TLS {
		bsClient.BaseURL = "https://" + config.Hostname
	} else {
		bsClient.BaseURL = "http://" + config.Hostname
	}
	return bsClient, nil
}

func (c *client) Type() string {
	return model.BitbucketServerType
}

func (c *client) CanLogin() bool {
	return true
}

func (c *client) CanRepos() bool {
	return true
}

func (c *client) CanHook() bool {
	return true
}

func (c *client) getOauthConsumer() (*oauth.Consumer, error) {
	keyBytes := []byte(c.PrivateKey)
	block, _ := pem.Decode(keyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	bitbucketOauthConsumer := oauth.NewRSAConsumer(
		c.ConsumerKey,
		privateKey,
		oauth.ServiceProvider{
			RequestTokenUrl:   fmt.Sprintf(requestTokenURL, c.BaseURL),
			AuthorizeTokenUrl: fmt.Sprintf(authorizeTokenURL, c.BaseURL),
			AccessTokenUrl:    fmt.Sprintf(accessTokenURL, c.BaseURL),
			HttpMethod:        http.MethodPost,
		})
	return bitbucketOauthConsumer, nil
}

func (c *client) Login(code string) (*v3.SourceCodeCredential, error) {
	splits := strings.SplitN(code, ":", 2)
	if len(splits) < 2 {
		return nil, errors.New("invalid code")
	}
	OAuthToken := splits[0]
	OAuthVerifier := splits[1]
	consumer, err := c.getOauthConsumer()
	if err != nil {
		return nil, err
	}
	requestToken := &oauth.RequestToken{
		Token: OAuthToken,
	}
	token, err := consumer.AuthorizeToken(requestToken, OAuthVerifier)
	if err != nil {
		return nil, err
	}

	user, err := c.getUser(token.Token)
	if err != nil {
		return nil, err
	}
	cred := convertUser(user)
	cred.Spec.AccessToken = token.Token
	return cred, nil
}

func (c *client) Repos(account *v3.SourceCodeCredential) ([]v3.SourceCodeRepository, error) {
	if account == nil {
		return nil, fmt.Errorf("empty account")
	}
	url := c.BaseURL + "/rest/api/1.0/repos?permission=REPO_ADMIN"
	hasNext := true
	var repos []Repository
	for hasNext {
		b, err := c.getFromBitbucket(url, account.Spec.AccessToken)
		if err != nil {
			return nil, err
		}
		var pageRepos PaginatedRepositories
		if err := json.Unmarshal(b, &pageRepos); err != nil {
			return nil, err
		}
		hasNext = !pageRepos.IsLastPage
		url = fmt.Sprintf("%s/rest/api/1.0/repos?permission=REPO_ADMIN&start=%d", c.BaseURL, pageRepos.NextPageStart)
		repos = append(repos, pageRepos.Values...)
	}

	return convertRepos(repos), nil
}

func (c *client) CreateHook(pipeline *v3.Pipeline, accessToken string) (string, error) {
	user, repo, err := getUserRepoFromURL(pipeline.Spec.RepositoryURL)
	if err != nil {
		return "", err
	}
	hookURL := fmt.Sprintf("%s/hooks?pipelineId=%s", settings.ServerURL.Get(), ref.Ref(pipeline))
	hook := Hook{
		Name:   "pipeline webhook",
		URL:    hookURL,
		Active: true,
		Configuration: HookConfiguration{
			Secret: pipeline.Status.Token,
		},
		Events: []string{
			"repo:refs_changed",
			"pr:opened",
			"pr:modified",
		},
	}

	url := fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos/%s/webhooks", c.BaseURL, user, repo)
	b, err := json.Marshal(hook)
	if err != nil {
		return "", err
	}
	reader := bytes.NewReader(b)
	resp, err := c.doRequestToBitbucket(http.MethodPost, url, accessToken, nil, reader)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(resp, &hook)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(hook.ID), nil
}

func (c *client) DeleteHook(pipeline *v3.Pipeline, accessToken string) error {
	user, repo, err := getUserRepoFromURL(pipeline.Spec.RepositoryURL)
	if err != nil {
		return err
	}

	hook, err := c.getHook(pipeline, accessToken)
	if err != nil {
		return err
	}
	if hook != nil {
		url := fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos/%s/webhooks/%d", c.BaseURL, user, repo, hook.ID)
		_, err := c.doRequestToBitbucket(http.MethodDelete, url, accessToken, nil, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *client) getHook(pipeline *v3.Pipeline, accessToken string) (*Hook, error) {
	user, repo, err := getUserRepoFromURL(pipeline.Spec.RepositoryURL)
	if err != nil {
		return nil, err
	}

	var hooks PaginatedHooks
	var result *Hook
	url := fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos/%s/webhooks", c.BaseURL, user, repo)

	b, err := c.getFromBitbucket(url, accessToken)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &hooks); err != nil {
		return nil, err
	}
	for _, hook := range hooks.Values {
		if strings.HasSuffix(hook.URL, fmt.Sprintf("hooks?pipelineId=%s", ref.Ref(pipeline))) {
			result = &hook
			break
		}
	}
	return result, nil
}

func (c *client) getFileFromRepo(filename string, owner string, repo string, branch string, accessToken string) ([]byte, error) {
	url := fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos/%s/raw/%s?at=%s", c.BaseURL, owner, repo, filename, branch)
	return c.getFromBitbucket(url, accessToken)
}

func (c *client) GetPipelineFileInRepo(repoURL string, branch string, accessToken string) ([]byte, error) {
	owner, repo, err := getUserRepoFromURL(repoURL)
	if err != nil {
		return nil, err
	}
	content, err := c.getFileFromRepo(utils.PipelineFileYaml, owner, repo, branch, accessToken)
	if err != nil {
		//look for both suffix
		content, err = c.getFileFromRepo(utils.PipelineFileYml, owner, repo, branch, accessToken)
	}
	if err != nil {
		logrus.Debugf("error GetPipelineFileInRepo - %v", err)
		return nil, nil
	}
	logrus.Infof("I got pipeline content:%s", string(content))
	return content, nil
}

func (c *client) SetPipelineFileInRepo(repoURL string, branch string, accessToken string, content []byte) error {
	owner, repo, err := getUserRepoFromURL(repoURL)
	if err != nil {
		return err
	}

	currentContent, err := c.getFileFromRepo(utils.PipelineFileYml, owner, repo, branch, accessToken)
	currentFileName := utils.PipelineFileYml
	if err != nil {
		if httpErr, ok := err.(*httperror.APIError); !ok || httpErr.Code.Status != http.StatusNotFound {
			return err
		}
		//look for both suffix
		currentContent, err = c.getFileFromRepo(utils.PipelineFileYaml, owner, repo, branch, accessToken)
		if err != nil {
			if httpErr, ok := err.(*httperror.APIError); !ok || httpErr.Code.Status != http.StatusNotFound {
				return err
			}
		} else {
			currentFileName = utils.PipelineFileYaml
		}
	}

	apiurl := fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos/%s/browse/%s", c.BaseURL, owner, repo, currentFileName)
	message := "Create .rancher-pipeline.yml file"
	if currentContent != nil {
		//update pipeline file
		message = fmt.Sprintf("Update %s file", currentFileName)
	}

	data := url.Values{}
	data.Set("message", message)
	data.Set("branch", branch)
	data.Set("content", string(content))
	//TODO
	//data.Set("sourceCommitId","")
	data.Encode()
	reader := strings.NewReader(data.Encode())
	header := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	_, err = c.doRequestToBitbucket(http.MethodPut, apiurl, accessToken, header, reader)

	return err
}

func (c *client) GetBranches(repoURL string, accessToken string) ([]string, error) {
	owner, repo, err := getUserRepoFromURL(repoURL)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos/%s/branches", c.BaseURL, owner, repo)

	b, err := c.getFromBitbucket(url, accessToken)
	if err != nil {
		return nil, err
	}
	var branches PaginatedBranches
	if err := json.Unmarshal(b, &branches); err != nil {
		return nil, err
	}
	result := []string{}
	for _, b := range branches.Values {
		if b.Type != "BRANCH" {
			continue
		}
		result = append(result, b.DisplayID)
	}
	return result, nil
}

func (c *client) GetHeadInfo(repoURL string, branch string, accessToken string) (*model.BuildInfo, error) {
	owner, repo, err := getUserRepoFromURL(repoURL)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos/%s/branches", c.BaseURL, owner, repo)

	b, err := c.getFromBitbucket(url, accessToken)
	if err != nil {
		return nil, err
	}
	var branches PaginatedBranches
	if err := json.Unmarshal(b, &branches); err != nil {
		return nil, err
	}
	headCommit := ""
	for _, b := range branches.Values {
		if b.DisplayID == branch {
			headCommit = b.LatestCommit
			break
		}
	}
	if headCommit == "" {
		return nil, fmt.Errorf("cannot find head commit of branch '%s'", branch)
	}
	url = fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos/%s/commits/%s", c.BaseURL, owner, repo, headCommit)
	b, err = c.getFromBitbucket(url, accessToken)
	if err != nil {
		return nil, err
	}
	var commit Commit
	if err := json.Unmarshal(b, &commit); err != nil {
		return nil, err
	}
	info := &model.BuildInfo{}
	info.Commit = commit.ID
	info.Ref = "refs/head/" + branch
	info.Branch = branch
	info.Message = commit.Message
	info.HTMLLink = fmt.Sprintf("%s/projects/%s/repos/%s/commits/%s", c.BaseURL, owner, repo, headCommit)
	info.AvatarURL = fmt.Sprintf("%s/users/%s/avatar.png", c.BaseURL, commit.Author.Name)
	info.Author = commit.Author.Name

	return info, nil
}

func convertUser(bitbucketUser *User) *v3.SourceCodeCredential {

	if bitbucketUser == nil {
		return nil
	}
	cred := &v3.SourceCodeCredential{}
	cred.Spec.SourceCodeType = model.BitbucketServerType

	cred.Spec.AvatarURL = bitbucketUser.Links.Avatar.Href
	cred.Spec.HTMLURL = bitbucketUser.Links.Html.Href
	cred.Spec.LoginName = bitbucketUser.Name
	cred.Spec.DisplayName = bitbucketUser.DisplayName

	return cred

}

func (c *client) getUser(accessToken string) (*User, error) {
	url := c.BaseURL + "/plugins/servlet/applinks/whoami"
	b, err := c.getFromBitbucket(url, accessToken)
	if err != nil {
		return nil, err
	}
	username := string(b)
	url = fmt.Sprintf("%s/rest/api/1.0/users/%s", c.BaseURL, username)
	b, err = c.getFromBitbucket(url, accessToken)
	if err != nil {
		return nil, err
	}
	user := &User{}
	if err := json.Unmarshal(b, user); err != nil {
		return nil, err
	}
	return user, nil
}

func convertRepos(repos []Repository) []v3.SourceCodeRepository {
	result := []v3.SourceCodeRepository{}
	for _, repo := range repos {
		r := v3.SourceCodeRepository{}
		for _, link := range repo.Links.Clone {
			if strings.HasPrefix(link.Name, "http") {
				r.Spec.URL = link.Href
				break
			}
		}
		r.Spec.Permissions.Admin = true
		r.Spec.Permissions.Pull = true
		r.Spec.Permissions.Push = true
		result = append(result, r)
	}
	return result
}

func (c *client) getFromBitbucket(url string, accessToken string) ([]byte, error) {
	return c.doRequestToBitbucket(http.MethodGet, url, accessToken, nil, nil)
}

func (c *client) doRequestToBitbucket(method string, url string, accessToken string, header map[string]string, body io.Reader) ([]byte, error) {
	logrus.Infof("requesting to url:%s", url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	consumer, err := c.getOauthConsumer()
	if err != nil {
		return nil, err
	}
	var token oauth.AccessToken
	token.Token = accessToken
	client, err := consumer.MakeHttpClient(&token)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	if method == http.MethodGet {
		q.Set("limit", maxPerPage)
	}
	req.URL.RawQuery = q.Encode()
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var body bytes.Buffer
		io.Copy(&body, resp.Body)
		return nil, httperror.NewAPIErrorLong(resp.StatusCode, "", body.String())
	}
	r, err := ioutil.ReadAll(resp.Body)
	logrus.Infof("response is:%s", string(r))
	return r, err
}

func getUserRepoFromURL(repoURL string) (string, string, error) {
	reg := regexp.MustCompile(".*/([^/]*?)/([^/]*?).git")
	match := reg.FindStringSubmatch(repoURL)
	if len(match) != 3 {
		return "", "", fmt.Errorf("error getting user/repo from gitrepoUrl:%v", repoURL)
	}
	return match[1], match[2], nil
}
