package feishu

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/fastwego/feishu"
	"github.com/fastwego/feishu/apis/bot/group_manage"
	"github.com/fastwego/feishu/apis/message"
)

const messageType = "text"

type (
	App struct {
		app *feishu.App
	}
	content struct {
		Text string `json:"text"`
	}

	Message struct {
		UserId  string  `json:"user_id"`
		ChatId  string  `json:"chat_id"`
		Email   string  `json:"email"`
		MsgType string  `json:"msg_type"`
		Content content `json:"content"`
	}
)

func NewApp(appid, secret string) *App {
	return &App{
		app: feishu.NewApp(feishu.AppConfig{
			AppId:     appid,
			AppSecret: secret,
		}),
	}
}

func (app *App) ListChatGroup(title string) {
	v, err := group_manage.ChatList(app.app, url.Values{})
	if err != nil {
		log.Fatalf("list chat error: %v", err)
	}
	m := map[string]interface{}{}
	_ = json.Unmarshal(v, &m)
	b, _ := json.MarshalIndent(m, "", "    ")
	log.Printf("chat result is:\n%s", string(b))

}

func (app *App) Send(receiver, receiverEmail, chatId, text string) error {
	payload, err := json.Marshal(Message{
		UserId:  receiver,
		Email:   receiverEmail,
		MsgType: messageType,
		ChatId:  chatId,
		Content: content{
			Text: text,
		},
	})
	if err != nil {
		return err
	}

	_, err = message.Send(app.app, payload)
	return err
}
