package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"testing"
	"errors"

	"github.com/stretchr/testify/mock"
	//"github.com/stretchr/testify/assert"
	"github.com/stellar/go-stellar-base/horizon"
	. "github.com/smartystreets/goconvey/convey"
)



type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Get(dest interface{}, query string, args ...interface{}) error {
	a := m.Called(dest, query, args[0])
	return a.Error(0)
}

func (m *MockDatabase) Exec(query string) (sql.Result, error) {
	a := m.Called(query)
	return a.Get(0).(sql.Result), a.Error(1)
}

func (m *MockDatabase) getLastProcessed(threadID string) (res string, err error) {
	a := m.Called(threadID)
	return a.String(0), a.Error(1)
}

func (m *MockDatabase) setLastProcessed(threadID string, since string) error {
	a := m.Called(threadID,since)
	return a.Error(0)
}


func (m *MockDatabase) getTip(threadID string) (res string, err error) {
	a := m.Called(threadID)
	return a.String(0), a.Error(1)
}

func (m *MockDatabase) getUser(name string) (user User, err error) {
	a := m.Called(name)
	return a.Get(0).(User), a.Error(1)
}

func (m *MockDatabase) createUser(destination string) (user User, err error) {
	a := m.Called(destination)
	return a.Get(0).(User), a.Error(1)
}

func (m *MockDatabase) getPreset(name string, userID int) (pre Preset, err error) {
	//a := m.Called(name,userID)
	//return a.Get(0).(Preset), a.Error(1)
	return Preset{},nil
}



type MockHorizon struct {
	mock.Mock
}

func (m *MockHorizon)  setup(app *App) {

}

func (m *MockHorizon) LoadAccount(accountID string) (account horizon.Account, err error) {
	a := m.Called(accountID)
	return a.Get(0).(horizon.Account), a.Error(1)
}


type MockBridge struct {
	mock.Mock
}

func (m *MockBridge)  setup(app *App) {

}

func (m *MockBridge) SendTx(sourceUser User, destAddress string, amount float64) error {
	a := m.Called(sourceUser,destAddress,amount)
	return a.Error(0)
}


func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestParsing(t *testing.T) {
	log.Println("in testing2")

	var amount float64
	var assetCode string

	amount = getAmount("16577.26")
	if amount != 16577.26 {
		t.Error("Expected 16577.26, got ", amount)
	}

	amount = getAmount("16577.026")
	if amount != 16577.026 {
		t.Error("Expected 16577.026, got ", amount)
	}

	amount = getAmount("16577.26sfs")
	if amount != 0 {
		t.Error("Expected 0, got ", amount)
	}

	assetCode, amount = getPrefixedAmount("12.5")
	if amount > 0 {
		t.Error("Expected 0, got ", assetCode, amount)
	}

	assetCode, amount = getPrefixedAmount("12.5$")
	if assetCode != "USD" || amount != 12.5 {
		t.Error("Expected 12.5 USD, got ", assetCode, amount)
	}

	assetCode, amount = getPrefixedAmount("*13")
	if assetCode != "XLM" || amount != 13 {
		t.Error("Expected 13 XLM, got ", assetCode, amount)
	}

	assetCode, amount = getPrefixedAmount("€132.05")
	if assetCode != "EUR" || amount != 132.05 {
		t.Error("Expected 132.05 EUR, got ", assetCode, amount)
	}

	dest := getDestination("@joe")
	if dest != "joe" {
		t.Error("Expected joe EUR, got ", dest)
	}

	dest = getDestination("joe")
	if dest != "" {
		t.Error("Expected empty EUR, got ", dest)
	}

	assetCode, amount = standardPresets("日本語")
	if amount != 0 {
		t.Error("Expected 0 , got ", assetCode, amount)
	}

	assetCode, amount = standardPresets("EUR")
	if amount != 1 {
		t.Error("Expected 1 EUR, got ", assetCode, amount)
	}
}


/*
Test:
	Sending default amount
	Sending custom amount
   Sending to someone not yet in the system
   Sending to someone in the system
   Sending to someone with out a GH account
   Sending more than you have
   Sending when you don't have an account
   Sending to tipbot
   Sending to yourself

   Two send commands in the same line

   emptying to a payment address
   emptying to a Stellar account ID
   Emptying a stellar account that doesn't exist
   emptying to an invalid accountID
   emptying to a payment address that isn't found

 */

/*
		mockDatabase.On("Get", &responseRecord, "SELECT ", "bob").Return(nil).Run(func(args mock.Arguments) {
			record := args.Get(0).(*FedRecord)
			record.AccountId = accountId
			record.StellarAddress = username + "*" + app.config.Domain
		})

		mockDatabase.On("getUser", "bob").Return(nil).Run(func(args mock.Arguments) {
			record := args.Get(0).(*FedRecord)
			record.AccountId = accountId
			record.StellarAddress = username + "*" + app.config.Domain
		})
*/

// test that parsing various messages results in the right things being put in the DB and the right things sent over stellar
func TestSendingTip(t *testing.T) {
	mockDatabase := new(MockDatabase)
	mockHorizon := new(MockHorizon)
	mockBridge := new(MockBridge)

	app := App{
		config: Config{
			GITHUB_ACCESS_TOKEN:  "1",
			BOT_GITHUB_NAME:      "tipbot",
			DB_CONNECTION_STRING: "",
			MIN_XML_BALANCE:      20,
			WEBSITE_URL:          "test.com",
			HORIZON_URL:          "test.com/horizon",
			FEDERATION_DOMAIN:    "ttest.com",
			DEFAULT_TIP_AMOUNT:   200,
		},
		database: mockDatabase,
		horizon:  mockHorizon,
		bridge: mockBridge,
	}

	bobUser := User{
		UserID: 1,
		GithubName: "bob",
		AccountID:  "bobAccountID",
		SecretKey:  "bobSecretKey",

		StellarAccount: horizon.Account{
			SubentryCount: 1,
			Balances: [ ]horizon.Balance{
				horizon.Balance{
					Balance: horizon.Amount(5000),
					Asset: horizon.Asset{
						Type: "native",
					},
				},
			},
		},
	}

	oldGuyUser := User{
		UserID: 1,
		GithubName: "OldGuy",
		AccountID:  "ogAccountID",
		SecretKey:  "ogSecretKey",

		StellarAccount: horizon.Account{
			SubentryCount: 1,
			Balances: [ ]horizon.Balance{
				horizon.Balance{
					Balance: horizon.Amount(4000),
					Asset: horizon.Asset{
						Type: "native",
					},
				},
			},
		},
	}

	newGuyUser := User{
		UserID: 3,
		GithubName: "NewGuy",
		AccountID:  "ngAccountID",
		SecretKey:  "ngSecretKey",
	}

	newGuyIssue := GithubIssue{ThreadID: "2", Owner: "NewGuy", Repo: "repo1", IssueNumber: 5}
	oldGuyIssue := GithubIssue{ThreadID: "22", Owner: "OldGuy", Repo: "repo12", IssueNumber: 52}

	mockDatabase.On("getUser", "bob").Return(bobUser,nil)
	mockDatabase.On("getUser", "OldGuy").Return(oldGuyUser,nil)
	mockDatabase.On("getUser", "NewGuy").Return(User{},nil)
	mockDatabase.On("createUser", "NewGuy").Return(newGuyUser,nil)


	mockHorizon.On("LoadAccount",bobUser.AccountID).Return(bobUser.StellarAccount,nil)
	mockHorizon.On("LoadAccount",oldGuyUser.AccountID).Return(oldGuyUser.StellarAccount,nil)
	mockHorizon.On("LoadAccount",newGuyUser.AccountID).Return(newGuyUser.StellarAccount,errors.New("nope"))




	Convey("Tipping default amount", t, func() {
		mockBridge.On("SendTx",bobUser,oldGuyUser.AccountID,200.0).Return(nil)
		message := "hey @tipbot "
		app.parseBody(oldGuyIssue, "bob", message)

		// make sure correct amount was sent
		mockBridge.AssertCalled(t,"SendTx",bobUser,oldGuyUser.AccountID,200.0)
	})

	Convey("Tipping custom amount", t, func() {
		mockBridge.On("SendTx",bobUser,oldGuyUser.AccountID,205.0).Return(nil)
			message := "hey @tipbot 205"
			app.parseBody(oldGuyIssue, "bob", message)

		// make sure correct amount was sent
		mockBridge.AssertCalled(t,"SendTx",bobUser,oldGuyUser.AccountID,205.0)

		})

	Convey("Tipping someone not yet in the system", t, func() {
		message := "hey @tipbot "
		mockBridge.On("SendTx",bobUser,newGuyUser.AccountID,200.0).Return(nil)
		app.parseBody(newGuyIssue, "bob", message)

		mockBridge.AssertCalled(t,"SendTx",bobUser,newGuyUser.AccountID,200.0)

	})

	Convey("Tipping too little to someone not yet in the system", t, func() {
		message := "hey @tipbot 5"
		app.parseBody(newGuyIssue, "bob", message)

	})

	Convey("Tipping more than you have", t, func() {
		message := "hey @tipbot 6000 "

		app.parseBody(oldGuyIssue, "bob", message)

	})

	Convey("Tipping when you don't have an account", t, func() {
		message := "hey @tipbot "

		app.parseBody(oldGuyIssue, "NewGuy", message)

	})

	Convey("Two tip commands in the same line", t, func() {
		message := "hey @tipbot 205 @OldGuy   @tipbot 100"
		mockBridge.On("SendTx",bobUser,oldGuyUser.AccountID,205.0).Return(nil)

		app.parseBody(oldGuyIssue, "bob", message)

	})

	Convey("Tipping yourself", t, func() {
		message := "hey @tipbot "
		app.parseBody(oldGuyIssue, "OldGuy", message)
	})

	Convey("emptying to a payment address", t, func() {
		mockBridge.On("SendTx",bobUser,"bob*stellar.org",float64(bobUser.StellarAccount.GetNativeBalance()-bobUser.getMinBalance())).Return(nil)

		message := "hey @tipbot empty bob*stellar.org "

			app.parseBody(oldGuyIssue, "bob", message)

		})

	Convey("emptying to nothing", t, func() {
		message := "hey @tipbot empty "


		app.parseBody(oldGuyIssue, "bob", message)

	})

	Convey("emptying to a Stellar account ID", t, func() {
		mockBridge.On("SendTx",bobUser,"GBCXF42Q26WFS2KJ5XDM5KGOWR5M4GHR3DBTFBJVRYKRUYJK4DBIH3RX",float64(bobUser.StellarAccount.GetNativeBalance()-bobUser.getMinBalance())).Return(nil)
		message := "@tipbot empty GBCXF42Q26WFS2KJ5XDM5KGOWR5M4GHR3DBTFBJVRYKRUYJK4DBIH3RX"

		app.parseBody(oldGuyIssue, "bob", message)
		//mockBridge.(t,"SendTx",bobUser,oldGuyUser.AccountID,100.0)

	})

	Convey("emptying to an invalid accountID", t, func() {
		// TODO: have this return an error
		mockBridge.On("SendTx",bobUser,"invalidaccountID",float64(bobUser.StellarAccount.GetNativeBalance()-bobUser.getMinBalance())).Return(nil)

		message := "hey @tipbot empty invalidaccountID"

		app.parseBody(oldGuyIssue, "bob", message)

	})

	Convey("emptying to a payment address that isn't found", t, func() {
		message := "hey @tipbot unknown*stellar.org"
		// TODO: have this return an error
		mockBridge.On("SendTx",bobUser,"unknown*stellar.org",float64(bobUser.StellarAccount.GetNativeBalance()-bobUser.getMinBalance())).Return(nil)
		//mockDatabase.On("getPreset", "unknown*stellar.org",1).Return(Preset{},nil)

		app.parseBody(oldGuyIssue, "bob", message)

	})

	Convey("emptying to a payment address that isn't found", t, func() {
		message := "@tipbot hi    out there"
		// TODO: have this return an error
		//mockBridge.On("SendTx",bobUser,"unknown*stellar.org",float64(bobUser.StellarAccount.GetNativeBalance()-bobUser.getMinBalance())).Return(nil)

		app.parseBody(oldGuyIssue, "bob", message)

	})
}
