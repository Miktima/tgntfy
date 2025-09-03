package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

	chats_c.fs.Bool("verbose", false, "verbose chat list updates of telegram bot")

	return chats_c
}

type ChatsCmd struct {
	fs *flag.FlagSet
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
		fmt.Printf("Cannot create new request  %s, error: %v\n", url, err)
	}
	// Send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error with GET request: %v\n", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))

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
