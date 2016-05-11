package main

import (
	"github.com/stellar/go-stellar-base/horizon"
)

type User struct {
	UserID     int     `db:"user_id"`
	GithubName string  `db:"github_name"`
	AccountID  string  `db:"account_id"`
	SecretKey  string  `db:"secret_key"`

	StellarAccount horizon.Account
}

func (self *User) getMinBalance() float64 {
	return float64(self.StellarAccount.SubentryCount) * 10
}
