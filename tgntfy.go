package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
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

func Chats(verbose, text bool) error {

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
	if !verbose && !text {
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
	} else if verbose {
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
	} else if text {
		if messages, ok := result["result"].([]interface{}); ok {
			for _, msgi := range messages {
				if msg, ok := msgi.(map[string]interface{}); ok {
					if m, ok := msg["message"].(map[string]interface{}); ok {
						if from, ok := m["from"].(map[string]interface{}); ok {
							fmt.Println("Message:")
							if id, ok := from["id"].(float64); ok {
								fmt.Println("     ID: ", int(id))
							}
						}
						if _, ok := m["text"].(string); ok {
							fmt.Println("     Text: ", m["text"])
						}
						fmt.Println("--------------------------------------")
					}
				}
			}
		}
	}
	return nil
}

func sendTlgrm(idchats []string, message string) error {

	const bold string = "2d5fef6c87f16217aa1b50e1ebc89720"
	const br string = "8b0f0ea73162b7552dda3c149b6c045d"
	const italic string = "3c4abf728013d2a85b9fd0e44dbc2353"

	// Read API key from file
	keyAPI, err := ReadKeyAPI()
	if err != nil {
		return err
	}
	reBold := regexp.MustCompile(`<[\/]?b>`)
	reBr := regexp.MustCompile(`<br>`)
	reItalic := regexp.MustCompile(`<[\/]?i>`)
	retag := regexp.MustCompile(`<.*>?`)
	resmb := regexp.MustCompile(`([_\*\[\]\(\)~\>\#\+\-\=\|\{\}\.!])`)

	message = reBold.ReplaceAllString(message, bold)
	message = reBr.ReplaceAllString(message, br)
	message = reItalic.ReplaceAllString(message, italic)
	message = retag.ReplaceAllString(message, "")

	message = resmb.ReplaceAllString(message, "\\$1")

	message = strings.ReplaceAll(message, bold, "*")
	message = strings.ReplaceAll(message, br, "\n")
	message = strings.ReplaceAll(message, italic, "_")

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Send message to Telegram
	client := &http.Client{}

	url := "https://api.telegram.org/bot" + keyAPI + "/sendMessage"

	for _, tgid := range idchats {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		q := req.URL.Query()
		q.Add("parse_mode", "MarkdownV2")
		q.Add("chat_id", tgid)
		q.Add("disable_web_page_preview", "1")
		q.Add("text", message)
		req.URL.RawQuery = q.Encode()
		// Отправляем запрос
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode > 299 {
			fmt.Println("Message with was not sent")
			fmt.Println("message: ", message)
			fmt.Println("Error: ", err)
			return err
		}
		defer resp.Body.Close()
	}
	return nil
}

func main() {
	// Chats sub-command parameters
	chatsCmd := flag.NewFlagSet("chats", flag.ExitOnError)
	chatsVerbose := chatsCmd.Bool("verbose", false, "verbose chat list updates of telegram bot")
	chatsText := chatsCmd.Bool("text", false, "Onle text and id chat of telegram bot")

	// Send sub-command parameters
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendIds := sendCmd.String("ids", "", "Chat IDs comma separated")
	sendMsg := sendCmd.String("text", "", "Message text")

	if len(os.Args) < 2 {
		fmt.Println("expected 'chats' or 'send' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "chats":
		chatsCmd.Parse(os.Args[2:])
		if *chatsVerbose && *chatsText {
			fmt.Println("Parameters 'verbose' and 'text' cannot be true at the same time")
			os.Exit(1)
		}
		err := Chats(*chatsVerbose, *chatsText)
		if err != nil {
			fmt.Println(err)
		}
	case "send":
		sendCmd.Parse(os.Args[2:])
		idchats := strings.Split(*sendIds, ",")
		err := sendTlgrm(idchats, *sendMsg)
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Println("expected 'chats' or 'send' subcommands")
		os.Exit(1)
	}
}
