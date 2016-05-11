package main

import (
	"log"
	"net/http"
	"strconv"
	"net/url"
)

type Bridge struct {
	app *App
}

func (self *Bridge) setup(a *App) {
	self.app = a
}

func (self *Bridge) SendTx(sourceUser User, destAddress string, amount float64) error {

	urlData := url.Values{}
	urlData.Set("source", sourceUser.SecretKey)
	urlData.Set("destination", destAddress)
	urlData.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))

	resp, err := http.PostForm(self.app.config.BRIDGE_URL + "/payment",urlData)

	if err != nil {
		log.Print(err)
		// TODO handle returns from bridge server and generate our own error types
		return err
	}

	defer resp.Body.Close()

	log.Print(resp)

	return nil
}
