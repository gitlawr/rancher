package bitbucketserver

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/mrjones/oauth"
	"github.com/rancher/rancher/pkg/randomtoken"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/types"
	"github.com/rancher/types/apis/project.cattle.io/v3"

	"github.com/rancher/norman/api/access"
	"github.com/rancher/norman/types/convert"
	"github.com/rancher/rancher/pkg/pipeline/remote"
	"github.com/rancher/rancher/pkg/pipeline/remote/model"
	"github.com/rancher/rancher/pkg/pipeline/utils"
	"github.com/rancher/rancher/pkg/ref"
	"github.com/rancher/types/client/project/v3"
	"github.com/satori/go.uuid"
	"io/ioutil"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"strings"
)

const (
	bitbucketDefaultHostName = "https://bitbucket.org"
	actionDisable            = "disable"
	actionTestAndApply       = "testAndApply"
	actionGenerateKeys       = "generateKeys"
	actionRequestLogin       = "requestLogin"
	actionLogin              = "login"
	projectNameField         = "projectName"
)

func (b *BitbucketServerProvider) Formatter(apiContext *types.APIContext, resource *types.RawResource) {
	if convert.ToBool(resource.Values["enabled"]) {
		resource.AddAction(apiContext, actionDisable)
	}

	resource.AddAction(apiContext, actionGenerateKeys)
	resource.AddAction(apiContext, actionRequestLogin)
	resource.AddAction(apiContext, actionTestAndApply)
}

func (b *BitbucketServerProvider) ActionHandler(actionName string, action *types.Action, request *types.APIContext) error {
	if actionName == actionTestAndApply {
		return b.testAndApply(actionName, action, request)
	} else if actionName == actionDisable {
		return b.disableAction(request)
	} else if actionName == actionGenerateKeys {
		return b.generateKeys(actionName, action, request)
	} else if actionName == actionRequestLogin {
		return b.requestLogin(actionName, action, request)
	}

	return httperror.NewAPIError(httperror.ActionNotAvailable, "")
}

func (b *BitbucketServerProvider) providerFormatter(apiContext *types.APIContext, resource *types.RawResource) {
	resource.AddAction(apiContext, actionLogin)
}

func (b *BitbucketServerProvider) providerActionHandler(actionName string, action *types.Action, request *types.APIContext) error {
	if actionName == actionLogin {
		return b.login(request)
	}

	return httperror.NewAPIError(httperror.ActionNotAvailable, "")
}

func (b *BitbucketServerProvider) generateKeys(actionName string, action *types.Action, apiContext *types.APIContext) error {
	ns, _ := ref.Parse(apiContext.ID)
	pConfig, err := b.GetProviderConfig(ns)
	if err != nil {
		return err
	}
	storedBitbucketPipelineConfig, ok := pConfig.(*v3.BitbucketServerPipelineConfig)
	if !ok {
		return fmt.Errorf("Failed to get bitbucket server provider config")
	}
	toUpdate := storedBitbucketPipelineConfig.DeepCopy()
	token, err := randomtoken.Generate()
	if err != nil {
		return err
	}
	toUpdate.ConsumerKey = token
	//FIXME
	pub, private, _ := generateKeyPair()
	toUpdate.PrivateKey = private
	toUpdate.PublicKey = pub
	if _, err = b.SourceCodeProviderConfigs.ObjectClient().Update(toUpdate.Name, toUpdate); err != nil {
		return err
	}

	return nil
}

func generateKeyPair() (string, string, error) {
	privateKey := `-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDNUw3p7+ZGX2mPKV27PUSyZWEcL+/ut78A1sfzuEv31DtHYFei
XV29/H/SUzqstr9OzdzGrwxXVBKPQXk53wAQ7eOWacUjQVqeL9HsC1+2c06pOuup
brUsAUe1o2n0juYwI7Olgynw2h5hfGwnlh27vJAURHjk3lbEWLnXVRuQwwIDAQAB
AoGBAI4aw3B7dtaRxo8sxBCI8Pi/LZzCmL6RMYK1JCJMFVfq7TQTO9PF5tFM5nJ8
5AkRWgqCdCCWmmX+a/H2EJ669mGDxeBe1lbXEBGW0dqTqIU7Rgs4ZdYoM3d5aIZg
I7HhUb11VqZlSL2q5wZPeSGf3Bxj8kLvljzWrkA7YoIWw+yRAkEA50GNXsAG4hrM
EmXB2jAVo8KTpUu9D+IE0CIXVlKmfGIAzBHSwf8GLizOeJhwTt87J2Lo9fW5HmfR
17wMHjemVwJBAONLMjTN1d7pQslXbqUp7xr0eQSLsS/1kmw9UCOyl/JHsDMjMj1r
VPBMuqfPgrdEe8YSgIJZ96aKY4sxCGwC7XUCQQCsYXzT6Cg5WuhLvnZmAfnffCc6
y94+fKhBzWe//RQFG7ikZZTI7yTYPqYZ1ufAoz4g+eXVkjlPpOwS+CXAUJM5AkEA
iC5rnFufQnl7vGqYLnkbe4jyYRjZRqTZ3+Q0ec7tXwo4tcrmtQnz0C4Iv7aC2Q89
IYXAXVlOGghcb+8m3qA6aQJAFSvcj/+cZPKPbMQqjbgCjyb3eBimMBo/cSY7Ogfp
Xd+ZSRUTPQ/LnfwS+bEuj7xYBD3y+pWOLqsIkRbfzziLNg==
-----END RSA PRIVATE KEY-----`
	publicKey := `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDNUw3p7+ZGX2mPKV27PUSyZWEc
L+/ut78A1sfzuEv31DtHYFeiXV29/H/SUzqstr9OzdzGrwxXVBKPQXk53wAQ7eOW
acUjQVqeL9HsC1+2c06pOuupbrUsAUe1o2n0juYwI7Olgynw2h5hfGwnlh27vJAU
RHjk3lbEWLnXVRuQwwIDAQAB
-----END PUBLIC KEY-----`
	return publicKey, privateKey, nil
}

func (b *BitbucketServerProvider) requestLogin(actionName string, action *types.Action, apiContext *types.APIContext) error {
	input := &v3.BitbucketServerRequestLoginInput{}
	if err := json.NewDecoder(apiContext.Request.Body).Decode(input); err != nil {
		return httperror.NewAPIError(httperror.InvalidBodyContent,
			fmt.Sprintf("Failed to parse body: %v", err))
	}

	ns, _ := ref.Parse(apiContext.ID)
	pConfig, err := b.GetProviderConfig(ns)
	if err != nil {
		return err
	}
	storedBitbucketPipelineConfig, ok := pConfig.(*v3.BitbucketServerPipelineConfig)
	if !ok {
		return fmt.Errorf("Failed to get bitbucket server provider config")
	}
	consumerKey := storedBitbucketPipelineConfig.ConsumerKey
	rsaKey := storedBitbucketPipelineConfig.PrivateKey
	host := ""
	if input.TLS {
		host = "https://" + input.Hostname
	} else {
		host = "http://" + input.Hostname
	}
	consumer, err := getOauthConsumer(consumerKey, rsaKey, host)
	if err != nil {
		return err
	}
	_, loginURL, err := consumer.GetRequestTokenAndUrl(input.RedirectURL)
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"loginUrl": loginURL,
		"type":     "bitbucketServerRequestLoginOutput",
	}
	apiContext.WriteResponse(http.StatusOK, data)
	return nil
}

func (b *BitbucketServerProvider) testAndApply(actionName string, action *types.Action, apiContext *types.APIContext) error {
	applyInput := &v3.BitbucketServerApplyInput{}

	if err := json.NewDecoder(apiContext.Request.Body).Decode(applyInput); err != nil {
		return httperror.NewAPIError(httperror.InvalidBodyContent,
			fmt.Sprintf("Failed to parse body: %v", err))
	}

	ns, _ := ref.Parse(apiContext.ID)
	pConfig, err := b.GetProviderConfig(ns)
	if err != nil {
		return err
	}
	storedBitbucketPipelineConfig, ok := pConfig.(*v3.BitbucketServerPipelineConfig)
	if !ok {
		return fmt.Errorf("Failed to get bitbucket server provider config")
	}
	toUpdate := storedBitbucketPipelineConfig.DeepCopy()
	toUpdate.Hostname = applyInput.Hostname
	toUpdate.TLS = applyInput.TLS
	toUpdate.UserName = applyInput.UserName
	toUpdate.Password = applyInput.Password
	toUpdate.RedirectURL = applyInput.RedirectURL

	code := fmt.Sprintf("%s:%s", applyInput.OAuthToken, applyInput.OAuthVerifier)
	//oauth and add user
	userName := apiContext.Request.Header.Get("Impersonate-User")
	sourceCodeCredential, err := b.authAddAccount(userName, code, toUpdate)
	if err != nil {
		return err
	}
	if _, err = b.refreshReposByCredentialAndConfig(sourceCodeCredential, toUpdate); err != nil {
		return err
	}

	toUpdate.Enabled = true
	//update bitbucket pipeline config
	if _, err = b.SourceCodeProviderConfigs.ObjectClient().Update(toUpdate.Name, toUpdate); err != nil {
		return err
	}

	apiContext.WriteResponse(http.StatusOK, nil)
	return nil
}

func (b *BitbucketServerProvider) login(apiContext *types.APIContext) error {
	loginInput := v3.BitbucketServerLoginInput{}
	requestBytes, err := ioutil.ReadAll(apiContext.Request.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(requestBytes, &loginInput); err != nil {
		return err
	}

	ns, _ := ref.Parse(apiContext.ID)
	pConfig, err := b.GetProviderConfig(ns)
	if err != nil {
		return err
	}
	config, ok := pConfig.(*v3.BitbucketServerPipelineConfig)
	if !ok {
		return fmt.Errorf("Failed to get bitbucket provider config")
	}
	if !config.Enabled {
		return errors.New("bitbucket oauth app is not configured")
	}

	//oauth and add user
	userName := apiContext.Request.Header.Get("Impersonate-User")
	code := fmt.Sprintf("%s:%s", loginInput.OAuthToken, loginInput.OAuthVerifier)
	account, err := b.authAddAccount(userName, code, config)
	if err != nil {
		return err
	}
	data := map[string]interface{}{}
	if err := access.ByID(apiContext, apiContext.Version, client.SourceCodeCredentialType, account.Name, &data); err != nil {
		return err
	}

	if _, err := b.refreshReposByCredentialAndConfig(account, config); err != nil {
		return err
	}

	apiContext.WriteResponse(http.StatusOK, data)
	return nil
}

func (b *BitbucketServerProvider) authAddAccount(userID string, code string, config *v3.BitbucketServerPipelineConfig) (*v3.SourceCodeCredential, error) {
	if userID == "" {
		return nil, errors.New("unauth")
	}

	remote, err := remote.New(config)
	if err != nil {
		return nil, err
	}
	account, err := remote.Login(code)
	if err != nil {
		return nil, err
	}
	account.Name = strings.ToLower(fmt.Sprintf("%s-%s-%s", config.Namespace, model.BitbucketServerType, account.Spec.LoginName))
	account.Namespace = userID
	account.Spec.UserName = userID
	account.Spec.ProjectName = config.ProjectName
	_, err = b.SourceCodeCredentials.Create(account)
	if apierror.IsAlreadyExists(err) {
		exist, err := b.SourceCodeCredentialLister.Get(userID, account.Name)
		if err != nil {
			return nil, err
		}
		account.ResourceVersion = exist.ResourceVersion
		return b.SourceCodeCredentials.Update(account)
	} else if err != nil {
		return nil, err
	}
	return account, nil
}

func (b *BitbucketServerProvider) disableAction(request *types.APIContext) error {
	ns, _ := ref.Parse(request.ID)
	o, err := b.SourceCodeProviderConfigs.ObjectClient().UnstructuredClient().GetNamespaced(ns, model.BitbucketServerType, metav1.GetOptions{})
	if err != nil {
		return err
	}
	u, _ := o.(runtime.Unstructured)
	config := u.UnstructuredContent()
	if convert.ToBool(config[client.SourceCodeProviderConfigFieldEnabled]) {
		config[client.SourceCodeProviderConfigFieldEnabled] = false
		if _, err := b.SourceCodeProviderConfigs.ObjectClient().Update(model.BitbucketServerType, o); err != nil {
			return err
		}
		if t := convert.ToString(config[projectNameField]); t != "" {
			return b.cleanup(t)
		}
	}

	return nil
}

func (b *BitbucketServerProvider) cleanup(projectID string) error {
	pipelines, err := b.PipelineIndexer.ByIndex(utils.PipelineByProjectIndex, projectID)
	if err != nil {
		return err
	}
	for _, p := range pipelines {
		pipeline, _ := p.(*v3.Pipeline)
		if err := b.Pipelines.DeleteNamespaced(pipeline.Namespace, pipeline.Name, &metav1.DeleteOptions{}); err != nil {
			return err
		}
	}

	pipelineExecutions, err := b.PipelineExecutionIndexer.ByIndex(utils.PipelineExecutionByProjectIndex, projectID)
	if err != nil {
		return err
	}
	for _, e := range pipelineExecutions {
		execution, _ := e.(*v3.PipelineExecution)
		if err := b.PipelineExecutions.DeleteNamespaced(execution.Namespace, execution.Name, &metav1.DeleteOptions{}); err != nil {
			return err
		}
	}

	crdKey := utils.ProjectNameAndSourceCodeTypeKey(projectID, model.BitbucketServerType)
	credentials, err := b.SourceCodeCredentialIndexer.ByIndex(utils.SourceCodeCredentialByProjectAndTypeIndex, crdKey)
	if err != nil {
		return err
	}
	for _, c := range credentials {
		credential, _ := c.(*v3.SourceCodeCredential)
		if err := b.SourceCodeCredentials.DeleteNamespaced(credential.Namespace, credential.Name, &metav1.DeleteOptions{}); err != nil {
			return err
		}
	}

	repoKey := utils.ProjectNameAndSourceCodeTypeKey(projectID, model.BitbucketServerType)
	repositories, err := b.SourceCodeRepositoryIndexer.ByIndex(utils.SourceCodeRepositoryByProjectAndTypeIndex, repoKey)
	if err != nil {
		return err
	}
	for _, r := range repositories {
		repo, _ := r.(*v3.SourceCodeRepository)
		if err := b.SourceCodeRepositories.DeleteNamespaced(repo.Namespace, repo.Name, &metav1.DeleteOptions{}); err != nil {
			return err
		}
	}

	return nil
}

func (b *BitbucketServerProvider) refreshReposByCredentialAndConfig(credential *v3.SourceCodeCredential, config *v3.BitbucketServerPipelineConfig) ([]v3.SourceCodeRepository, error) {
	namespace := credential.Namespace
	credentialID := ref.Ref(credential)
	remote, err := remote.New(config)
	if err != nil {
		return nil, err
	}
	repos, err := remote.Repos(credential)
	if err != nil {
		return nil, err
	}

	//remove old repos
	repositories, err := b.SourceCodeRepositoryIndexer.ByIndex(utils.SourceCodeRepositoryByCredentialIndex, credentialID)
	if err != nil {
		return nil, err
	}
	for _, r := range repositories {
		repo, _ := r.(*v3.SourceCodeRepository)
		if err := b.SourceCodeRepositories.DeleteNamespaced(namespace, repo.Name, &metav1.DeleteOptions{}); err != nil {
			return nil, err
		}
	}

	//store new repos
	for _, repo := range repos {
		if !repo.Spec.Permissions.Admin {
			//store only admin repos
			continue
		}
		repo.Spec.SourceCodeCredentialName = credentialID
		repo.Spec.UserName = credential.Spec.UserName
		repo.Spec.SourceCodeType = credential.Spec.SourceCodeType
		repo.Name = uuid.NewV4().String()
		repo.Namespace = namespace
		repo.Spec.ProjectName = credential.Spec.ProjectName
		if _, err := b.SourceCodeRepositories.Create(&repo); err != nil {
			return nil, err
		}
	}

	return repos, nil
}

func getOauthConsumer(consumerKey, rsaKey, hostURL string) (*oauth.Consumer, error) {
	keyBytes := []byte(rsaKey)
	block, _ := pem.Decode(keyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	bitbucketOauthConsumer := oauth.NewRSAConsumer(
		consumerKey,
		privateKey,
		oauth.ServiceProvider{
			RequestTokenUrl:   fmt.Sprintf("%s/plugins/servlet/oauth/request-token", hostURL),
			AuthorizeTokenUrl: fmt.Sprintf("%s/plugins/servlet/oauth/authorize", hostURL),
			AccessTokenUrl:    fmt.Sprintf("%s/plugins/servlet/oauth/access-token", hostURL),
			HttpMethod:        http.MethodPost,
		})
	return bitbucketOauthConsumer, nil
}
