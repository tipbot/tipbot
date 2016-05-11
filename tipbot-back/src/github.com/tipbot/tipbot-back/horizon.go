package main

import (
	"github.com/stellar/go-stellar-base/horizon"
	"log"
)

type Horizon struct {
	app     *App
	horizon horizon.Client
}

func (self *Horizon) setup(a *App) {
	self.app = a
	self.horizon.URL = a.config.HORIZON_URL
}

func (self *Horizon) LoadAccount(accountID string) (horizon.Account, error) {

	log.Println("LoadAccount: ",accountID)
	return self.horizon.LoadAccount(accountID)
}
