package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"stargazers/feishu"
	"stargazers/gh"

	"github.com/tal-tech/go-zero/core/conf"
)

var configFile = flag.String("f", "config.yaml", "the config file")

type (
	Feishu struct {
		AppId          string `json:"appId"`
		AppSecret      string `json:"appSecret"`
		Receiver       string `json:"receiver,optional"`
		ReceiverEmail  string `json:"receiver_email,optional"`
		ReceiverChatId string `json:"receiver_chat_id,optional`
	}

	Config struct {
		Token    string        `json:"token"`
		Repo     string        `json:"repo"`
		Interval time.Duration `json:"interval,default=1m"`
		Feishu   Feishu        `json:"feishu"`
	}
)

func main() {
	flag.Parse()

	var c Config
	conf.MustLoad(*configFile, &c)

	app := feishu.NewApp(c.Feishu.AppId, c.Feishu.AppSecret)

	if len(os.Args) > 1 && os.Args[1] == "list_chat" {
		app.ListChatGroup(os.Args[1])
	} else if len(os.Args) > 1 && os.Args[1] == "monitor" {
		mon := gh.NewMonitor(c.Repo, c.Token, c.Interval, func(text string) error {
			return app.Send(c.Feishu.Receiver, c.Feishu.ReceiverEmail, c.Feishu.ReceiverChatId, text)
		})
		log.Fatal(mon.Start())
	} else {
		fmt.Printf(`usage:
		stargazers <command>

The commands are:
		list_chat:		list the chats including the robot.
		monitor:			monitor the repo stargazers.
`)
	}
}
