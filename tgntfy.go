package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func ReadKeyAPI() (string, error) {

	// Path
	path, _ := os.Executable()
	path = path[:strings.LastIndex(path, "/")+1]

	var byteValue []byte

	// Read file with API key
	if _, err := os.Stat(path + "/key.txt"); err == nil {
		byteValue, err = os.ReadFile(path + "/key.txt")
		// if we os.ReadFile returns an error then handle it
		if err != nil {
			return "", err
		}
	} else {
		return "", err
	}
	return string(byteValue), nil
}

func RunChatsCmd() *ChatsCmd {
	chats_c := &ChatsCmd{
		fs: flag.NewFlagSet("chats", flag.ContinueOnError),
	}

	chats_c.fs.BoolVar(&chats_c.verbose, "verbose", false, "verbose chat list updates of telegram bot")

	return chats_c
}

type ChatsCmd struct {
	fs *flag.FlagSet

	verbose bool
}

func (g *ChatsCmd) Name() string {
	return g.fs.Name()
}

func (g *ChatsCmd) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *ChatsCmd) Run() error {

	// Read API key from file
	keyAPI, err := ReadKeyAPI()
	if err != nil {
		return err
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Send message to Telegram
	client := &http.Client{}

	url := "https://api.telegram.org/bot" + keyAPI + "/getUpdates"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	//fmt.Println(string(body))

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	switch g.verbose {
	case false:
		if messages, ok := result["result"].([]interface{}); ok {
			for _, msgi := range messages {
				if msg, ok := msgi.(map[string]interface{}); ok {
					if m, ok := msg["message"].(map[string]interface{}); ok {
						if from, ok := m["from"].(map[string]interface{}); ok {
							fmt.Println("From:")
							if id, ok := from["id"].(float64); ok {
								fmt.Println("     ID:", int(id))
							}
							if _, ok := from["first_name"].(string); ok {
								fmt.Println("     First Name:", from["first_name"])
							}
							if _, ok := from["last_name"].(string); ok {
								fmt.Println("     Last Name:", from["last_name"])
							}
							if _, ok := from["username"].(string); ok {
								fmt.Println("     Username:", from["username"])
							}
							fmt.Println("--------------------------------------")
						}
					}
				}
			}
		}
	case true:
		if messages, ok := result["result"].([]interface{}); ok {
			for _, msgi := range messages {
				if msg, ok := msgi.(map[string]interface{}); ok {
					if m, ok := msg["message"].(map[string]interface{}); ok {
						if from, ok := m["from"].(map[string]interface{}); ok {
							fmt.Println("From:")
							if id, ok := from["id"].(float64); ok {
								fmt.Println("     ID: ", int(id))
							}
							if _, ok := from["first_name"].(string); ok {
								fmt.Println("     First Name: ", from["first_name"])
							}
							if _, ok := from["last_name"].(string); ok {
								fmt.Println("     Last Name: ", from["last_name"])
							}
							if _, ok := from["username"].(string); ok {
								fmt.Println("     Username: ", from["username"])
							}
							if _, ok := from["is_bot"].(bool); ok {
								fmt.Println("     Is bot: ", from["is_bot"])
							}
							if _, ok := from["language_code"].(string); ok {
								fmt.Println("     Language: ", from["language_code"])
							}
						}
						if chat, ok := m["chat"].(map[string]interface{}); ok {
							fmt.Println("Chat:")
							if _, ok := chat["username"].(string); ok {
								fmt.Println("     Username: ", chat["username"])
							}
							if _, ok := chat["type"].(string); ok {
								fmt.Println("     Type: ", chat["type"])
							}
						}
						if _, ok := m["text"].(string); ok {
							fmt.Println("     Text: ", m["text"])
						}
						if t, ok := m["date"].(float64); ok {
							timestamp := time.Unix(int64(t), 0)
							fmt.Println("     Date: ", timestamp.Local())
						}
						fmt.Println("--------------------------------------")
					}
				}
			}
		}
	}
	return nil
}

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}

func root(args []string) error {
	if len(args) < 1 {
		return errors.New("You must pass a sub-command")
	}

	cmds := []Runner{
		RunChatsCmd(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			return cmd.Run()
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subcommand)
}

func main() {
	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
