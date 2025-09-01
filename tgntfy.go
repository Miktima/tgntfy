package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	var keyAPI string

	// Ключи для командной строки
	botUpdate := flag.Bool("update", false, "list telegram bot updates")

	flag.Parse()

	// Определяем путь
	path, _ := os.Executable()
	path = path[:strings.LastIndex(path, "/")+1]

	// Читаем файл с настройками telegram
	if _, err := os.Stat(path + "/key.txt"); err == nil {
		byteValue, err := os.ReadFile(path + "/key.txt")
		// if we os.ReadFile returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}
		keyAPI = string(byteValue)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Send message to Telegram
	client := &http.Client{}

	if *botUpdate {
		url := "https://api.telegram.org/bot" + keyAPI + "/getUpdates"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Cannot create new request  %s, error: %v\n", url, err)
		}
		// Отправляем запрос
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error with GET request: %v\n", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		fmt.Print(body)
	}

}
