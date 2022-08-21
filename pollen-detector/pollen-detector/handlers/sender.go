package handlers

import (
	"context"
	"fmt"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

// byteToInt64 converts a []byte into an int64
func byteToInt64(b []byte) (int64, error) {
	s := string(b)
	ss := strings.TrimSpace(s)
	i, err := strconv.ParseInt(ss, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// getSecret retrieves the secret from openfaas and makes it available for use.
func getSecret(secretName string) ([]byte, error) {
	secret, err := ioutil.ReadFile(fmt.Sprintf("/var/openfaas/secrets/%s", secretName))
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func Sender(d *Pollen) error {
	token, err := getSecret("tg-token")
	if err != nil {
		log.Fatalln("no telegram token found")
	}
	chatId, err := getSecret("chat-id")
	if err != nil {
		log.Fatalln("no telegram chat_id found")
	}

	cid, err := byteToInt64(chatId)
	if err != nil {
		return err
	}

	if d.Count == "Low" {
		log.Printf("no message sent. pollen count %q", d.Count)
		return nil
	}
	svc, err := telegram.New(strings.TrimSpace(string(token)))
	if err != nil {
		return err
	}
	svc.SetParseMode(telegram.ModeMarkdown)

	// chat-id which the messages get sent to.
	svc.AddReceivers(cid)

	notify.UseServices(svc)

	err = svc.Send(
		context.Background(),
		d.Header,
		fmt.Sprintf("\n\n%s: *%s*", d.Date, d.Count),
	)
	if err != nil {
		return err
	}

	return nil
}
