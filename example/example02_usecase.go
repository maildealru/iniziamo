package example

import (
	"context"
	"log"
)

type mailruOAuth2TokenInfo struct {
}

func mailruOAuth2TokenInfoParser(_ string, n int, body []byte) (mailruOAuth2TokenInfo, error) {
	//TODO
	return mailruOAuth2TokenInfo{}, nil
}

func DoMailruO2Login() {
	c := NewMailruOAuth2Client()

	r, err := c.ExchangeAuthCodeToToken(
		context.TODO(),
		MailruOAuth2ClientExchangeAuthCodeToTokenRequest(),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(r)
}
