package main

type Preset struct {
	UserID    int    `db:"user_id"`
	Preset    string  `db:"preset"`
	AssetCode string   `db:"asset_code"`
	Amount    float64  `db:"amount"`
}

type GithubIssue struct {
	ThreadID    string
	Owner       string
	Repo        string
	IssueNumber int
}
