package alert

import (
	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/parse"
	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/convert"
	"github.com/rancher/rancher/pkg/notifiers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func NotifierCollectionFormatter(apiContext *types.APIContext, collection *types.GenericCollection) {
	collection.AddAction(apiContext, "send")
}

func NotifierFormatter(apiContext *types.APIContext, resource *types.RawResource) {
	resource.AddAction(apiContext, "send")
}

func (h *Handler) NotifierActionHandler(actionName string, action *types.Action, apiContext *types.APIContext) error {
	switch actionName {
	case "send":
		return h.testNotifier(actionName, action, apiContext)
	}

	return httperror.NewAPIError(httperror.InvalidAction, "invalid action: "+actionName)

}

func (h *Handler) testNotifier(actionName string, action *types.Action, apiContext *types.APIContext) error {
	actionInput, err := parse.ReadBody(apiContext.Request)
	if err != nil {
		return err
	}

	msg := ""
	msgInf, exist := actionInput["message"]
	if exist {
		message, ok := msgInf.(string)
		if ok {
			msg = message
		}
	}

	if apiContext.ID != "" {
		parts := strings.Split(apiContext.ID, ":")
		ns := parts[0]
		id := parts[1]
		notifier, err := h.Notifiers.GetNamespaced(ns, id, metav1.GetOptions{})
		if err != nil {
			return err
		}
		return notifiers.SendMessage(notifier, "", msg)
	} else {

		slackConfigInterface, exist := actionInput["slackConfig"]
		if exist {
			slackConfig := convert.ToMapInterface(slackConfigInterface)
			url, ok := slackConfig["url"].(string)
			if ok {
				channel := convert.ToString(slackConfig["defaultRecipient"])
				return notifiers.TestSlack(url, channel, msg)
			}
		}

		smtpConfigInterface, exist := actionInput["smtpConfig"]
		if exist {
			smtpConfig := convert.ToMapInterface(smtpConfigInterface)
			host, ok := smtpConfig["host"].(string)
			if ok {
				port, _ := convert.ToNumber(smtpConfig["port"])
				password := convert.ToString(smtpConfig["password"])
				username := convert.ToString(smtpConfig["username"])
				sender := convert.ToString(smtpConfig["sender"])
				receiver := convert.ToString(smtpConfig["defaultRecipient"])
				tls := convert.ToBool(smtpConfig["tls"])
				return notifiers.TestEmail(host, password, username, int(port), tls, msg, receiver, sender)
			}
		}

		webhookConfigInterface, exist := actionInput["webhookConfig"]
		if exist {
			webhookConfig := convert.ToMapInterface(webhookConfigInterface)
			url, ok := webhookConfig["url"].(string)
			if ok {
				return notifiers.TestWebhook(url, msg)
			}
		}

		pagerdutyConfigInterface, exist := actionInput["pagerdutyConfig"]
		if exist {
			pagerdutyConfig := convert.ToMapInterface(pagerdutyConfigInterface)
			key, ok := pagerdutyConfig["serviceKey"].(string)
			if ok {
				return notifiers.TestPagerduty(key, msg)
			}
		}

		return httperror.NewAPIError(httperror.ErrorCode{Status: 400}, "Notifier not configured")
	}

	return nil
}
