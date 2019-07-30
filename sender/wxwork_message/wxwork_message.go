package wxwork_message

import (
	"context"
	"github.com/lanceryou/wxwork"
)

type WxWork struct {
	message *wxwork.WxWorkMessage
	appName string
	targets string
}

func (w *WxWork) Send(message string) {
	if message == "" {
		return
	}
	w.message.SendMessage(context.TODO(), w.appName, w.targets, message)
}

func (w *WxWork) String() string {
	return "wxwork_message"
}

func NewWxWork(message *wxwork.WxWorkMessage, appName string, targets string) *WxWork {
	return &WxWork{
		message: message,
		appName: appName,
		targets: targets,
	}
}
