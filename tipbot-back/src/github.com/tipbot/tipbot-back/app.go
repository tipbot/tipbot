package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/github"
	_ "github.com/lib/pq"
	"github.com/stellar/go-stellar-base/horizon"
	"golang.org/x/oauth2"
)

type DatabaseI interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Exec(sql string) (sql.Result, error)
	getLastProcessed(threadID string) (string, error)
	setLastProcessed(threadID string, since string) error
	getTip(threadID string) (string, error)
	getUser(name string) (User, error)
	createUser(destination string) (User, error)
	getPreset(name string, userID int) (Preset, error)
}

type BridgeI interface {
	setup(app *App)
	SendTx(sourceUser User, destAddress string, amount float64) error
}

type HorizonI interface {
	setup(app *App)
	LoadAccount(accountID string) (account horizon.Account, err error)
}

type App struct {
	config       Config
	database     DatabaseI
	githubClient *github.Client
	bridge       BridgeI
	horizon      HorizonI
}

// NewApp constructs an new App instance from the provided config.
func NewApp(config Config) (*App, error) {
	database := setupDB(config)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GITHUB_ACCESS_TOKEN},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	result := &App{config: config,
		database: &database,
		githubClient: client,
		bridge: &Bridge{},
		horizon: &Horizon{}}

	result.bridge.setup(result)
	result.horizon.setup(result)
/*
	log.Println("Loading Account")
	_, err := result.horizon.LoadAccount("GATH47CTFPFJSPOIYR3GYDPBMQZPBKCAIBXZJSJNEBJWAIV7AUKW4LUD") // TEMP

	if( err != nil){
		log.Println(err)
	}


	log.Fatal("Later") // TEMP
	*/
	return result, nil
}

func (self *App) Run() {
	log.Println("App run...")

	for true {
		self.fetchMentions()
		time.Sleep(time.Second * 5)
	}
}

func (self *App) fetchMentions() {

	// get all our unread notifications
	notifications, _, err := self.githubClient.Activity.ListNotifications(nil)
	if err != nil {
		log.Fatal(err)
	}

	//log.Println("Response: ", response)

	for _, notice := range notifications {
		//log.Println(github.Stringify(notice))

		if *notice.Reason == "mention" &&
			*notice.Subject.Type == "Issue" {

			since, err := self.database.getLastProcessed(*notice.ID)
			if err != nil {
				log.Fatal(err)
			}

			url := strings.Replace(*notice.Subject.URL, "/", " ", -1)
			var owner string
			var repo string
			var issueID int

			fmt.Sscanf(url, "https:  api.github.com repos %s %s issues %d", &owner, &repo, &issueID)
			issue := GithubIssue{ThreadID: *notice.ID, Owner: owner, Repo: repo, IssueNumber: issueID}

			self.findMentions(issue, since)
		}

		// mark notice as read
		self.githubClient.Activity.MarkThreadRead(*notice.ID)
	}
}

func (self *App) findMentions(issue GithubIssue, since string) {

	hubIssue, _, err := self.githubClient.Issues.Get(issue.Owner, issue.Repo, issue.IssueNumber)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(github.Stringify(issue),since)

	if len(since) == 0 {
		// This thread hasn't already been processed so we need to do the root
		self.parseBody(issue, *hubIssue.User.Login, *hubIssue.Body)
		self.database.setLastProcessed(issue.ThreadID, "2016-01-01T00:00:00Z")
	}


	if *hubIssue.Comments > 0 {
		var opt *github.IssueListCommentsOptions
		var tsince time.Time
		if len(since) > 0 {
			tsince, err = time.Parse(time.RFC3339, since)
			if err != nil {
				log.Fatal(err)
			}

			log.Println("Only showing comments since: ", since)
			opt = &github.IssueListCommentsOptions{Since: tsince}
		}

		comments, _, err := self.githubClient.Issues.ListComments(issue.Owner, issue.Repo, issue.IssueNumber, opt)
		if err != nil {
			log.Fatal(err)
		}

		var parsed bool
		var lastProcessed github.Timestamp
		for _, comment := range comments {

			if( (*comment.User.Login != self.config.BOT_GITHUB_NAME) &&
			    (*comment.UpdatedAt).After(tsince)){
				self.parseBody(issue, *comment.User.Login, *comment.Body)
				lastProcessed.Time = *comment.UpdatedAt
				parsed=true
			}

		}

		if(parsed) {
			lastProcessed.Add(time.Second)
			self.database.setLastProcessed(issue.ThreadID, lastProcessed.Format(time.RFC3339))
		}

	}
}

// super gross
func (self *App) parseBody(issue GithubIssue, sourceName string, body string) {
	log.Println(body)
	index := strings.Index(body, self.config.BOT_GITHUB_NAME)
	if index >= 0 { // this comment mentions tipbot

		// does source have an account set up already
		sourceUser, err := self.database.getUser(sourceName)
		if err != nil {
			log.Fatal(err)
		}

		// what asset is the source trying to send
		// how much is the source trying to send
		// does the source have enough to send?

		if sourceUser.UserID > 0 { // User has been made

			// strip just the line with the comment
			rest := body[index:]
			index = strings.Index(rest, "\n")
			if index > 0 {
				rest = rest[:index]
			}

			tokens := strings.Fields(rest)
			log.Println(tokens)

			var assetCode string = "XLM"
			var amount float64 = self.config.DEFAULT_TIP_AMOUNT
			var destination string = issue.Owner

			if len(tokens) > 1 {
				// first token is either Amount, destination, preset, or nothing
				if tokens[1] == "empty" {
					// second token must be the destination
					if len(tokens) > 2 {
						self.emptyAccount(issue, sourceUser, tokens[2])
						return
					} else {
						self.postReply(issue, self.MALFORMED_COMMAND(sourceUser.GithubName))
						return
					}
				} else {
					a := getAmount(tokens[1])
					if a > 0 { // first token is an amount
						amount = a
						// next token is either verb, destination or nothing
						if len(tokens) > 2 {
							d := getDestination(tokens[2])
							if len(d) > 0 {
								destination = d
							} else {
								ac, ignoredAmount := self.translatePreset(tokens[2], sourceUser.UserID)
								if ignoredAmount != 0 {
									assetCode = ac
									// next token is either destination or nothing
									if len(tokens) > 3 {
										d = getDestination(tokens[2])
										if len(d) > 0 {
											destination = d
										}
									}
								}
							}
						}

					} else {
						ac, a := getPrefixedAmount(tokens[1])
						if a > 0 { // first token is an amount+currency
							assetCode = ac
							amount = a
							// next token is either destination or nothing
							if len(tokens) > 2 {
								d := getDestination(tokens[2])
								if len(d) > 0 {
									destination = d
								}
							}
						} else {
							d := getDestination(tokens[1])
							if len(d) > 0 { // first token is a destination
								destination = d
							} else {

								ac, a = self.translatePreset(tokens[1], sourceUser.UserID)
								if a > 0 { // first token is a preset
									assetCode = ac
									amount = a

									// next token is either destination or nothing
									if len(tokens) > 2 {
										d = getDestination(tokens[2])
										if len(d) > 0 {
											destination = d
										}
									}
								}
							}
						}
					}

				}
			}

			log.Println("Parsed line: ",assetCode, amount, destination)
			self.sendTip(issue, sourceUser, assetCode, amount, destination)

		} else {
			// haven't made an account yet
			self.postReply(issue, self.SOURCE_ACCOUNT_NOT_MADE(sourceName))
		}
	}
}

func (self *App) translatePreset(token string, userID int) (string, float64) {

	preset, err := self.database.getPreset(token, userID)
	if err != nil {
		log.Fatal(err)
	}
	if preset.UserID > 0 {
		return preset.AssetCode, preset.Amount
	}
	return standardPresets(token)
}

func (self *App) postReply(issue GithubIssue, text string) {
	var comment github.IssueComment
	comment.Body = &text

	log.Println("postReply: " + text)
	if(self.githubClient != nil){
		self.githubClient.Issues.CreateComment(issue.Owner, issue.Repo, issue.IssueNumber, &comment)
	}

}

// send all but the minbalance to the destination
// For now just deal with XLM
func (self *App) emptyAccount(issue GithubIssue, sourceUser User, destination string) {

	account, err := self.horizon.LoadAccount(sourceUser.AccountID)
	if err != nil {
		log.Fatal(err)
	}
	sourceUser.StellarAccount = account
	amount := sourceUser.StellarAccount.GetNativeBalance() - sourceUser.getMinBalance()
	if amount < 1 {
		self.postReply(issue, self.ALREADY_EMPTY(sourceUser.GithubName))
		return
	}

	err = self.bridge.SendTx(sourceUser, destination, float64(amount))
	if err != nil {
		self.handleSendTxError(issue, sourceUser.GithubName,err)
		return
	}

	self.postReply(issue, self.EMPTY_SUCCESS(sourceUser.GithubName))
}

func (self *App) handleSendTxError(issue GithubIssue, sourceName string,err error) {

	switch err.Error() {
	//case "INVALID_ACCOUNT_ID:
	//	self.postReply(issue, self.INVALID_ACCOUNT(sourceName))
	case "404":
		self.postReply(issue, self.ADDRESS_NOT_FOUND(sourceName))
	default:
		self.postReply(issue, self.SOMETHING_WRONG(sourceName))
	}

}

func (self *App) sendTip(issue GithubIssue, sourceUser User, assetCode string, amount float64, destination string) {

	// for now only allow XLM to be sent around
	if assetCode != "XLM" {
		//self.postReply(issue, self.ONLY_XLM())
		return
	}

	if sourceUser.GithubName == destination {
		//self.postReply(issue, self.CANT_SEND_TO_SELF(sourceUser.GithubName))
		return
	}

	stellar, err := self.horizon.LoadAccount(sourceUser.AccountID)
	if err != nil {
		log.Print("Error Loading Account: ",sourceUser.AccountID, "  ", err)
		return
	}
	sourceUser.StellarAccount = stellar
	if sourceUser.StellarAccount.GetNativeBalance() == 0 {
		self.postReply(issue, self.SOURCE_ACCOUNT_NOT_MADE(sourceUser.GithubName))
		return
	}

	if float64(sourceUser.StellarAccount.GetNativeBalance()-sourceUser.getMinBalance()) < amount {
		self.postReply(issue, self.NOT_ENOUGH_BALANCE(sourceUser.GithubName))
		return
	}

	destUser, err := self.database.getUser(destination)
	if err != nil {
		log.Fatal(err)
	}
	if destUser.UserID == 0 { // destination doesn't exist yet
		if amount < self.config.MIN_XML_BALANCE {
			self.postReply(issue, self.TIP_NOT_HIGH_ENOUGH(sourceUser.GithubName,destination))
			return
		}

		destUser, err = self.database.createUser(destination)
		if err != nil {
			log.Print(err)
			return
		}

	}

	//
	// Create Tx
	err = self.bridge.SendTx(sourceUser, destUser.AccountID, amount)
	if err != nil {
		self.handleSendTxError(issue,sourceUser.GithubName,err)
		return
	}
	self.postReply(issue, self.REPORT_TIP(sourceUser.GithubName, destUser.GithubName))

}
