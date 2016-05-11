package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stellar/go-stellar-base/keypair"
)

type Database struct {
	database *sqlx.DB
}

func setupDB(config Config) Database {

	var database Database
	db, err := sqlx.Connect("postgres", config.DB_CONNECTION_STRING)
	if err != nil {
		panic(err)
	}

	database.database = db

	return database
}

func (db *Database) Get(dest interface{}, query string, args ...interface{}) error {
	return db.database.Get(dest, query, args)
}

func (db *Database) Exec(url string) (sql.Result, error) {
	return db.database.Exec(url)
}

func (db *Database) getLastProcessed(threadID string) (string, error) {

	sql := fmt.Sprintf("SELECT since from processed where thread_id=%s", threadID)
	//log.Println(sql)

	row, err := db.database.Query(sql)
	if err != nil {
		log.Println("select err: ", err)
		return "", err
	}

	if row != nil {
		row.Next()
		var msg string
		row.Scan(&msg)
		return msg, nil
	}

	return "", nil
}

func (db *Database) setLastProcessed(threadID string, since string) error {
	sql := fmt.Sprintf("SELECT since from processed where thread_id=%s", threadID)
	//log.Println(sql)

	rows, err := db.database.Query(sql)
	if err != nil {
		log.Println("select err: ", err)
		return err
	}

	if rows.Next() {
		sql = fmt.Sprintf("UPDATE processed set since='%s' where thread_id=%s", since, threadID)
	} else {
		sql = fmt.Sprintf("INSERT INTO processed (thread_id,since) values (%s,'%s')", threadID, since)
	}

	log.Println(sql)

	_, err = db.database.Exec(sql)
	return err
}

func (db *Database) getTip(threadID string) (string, error) {

	sql := fmt.Sprintf("SELECT since from processed where thread_id=%s", threadID)
	//log.Println(sql)

	row, err := db.database.Query(sql)
	if err != nil {
		log.Println("select err: ", err)
		return "", err
	}

	if row != nil {
		row.Next()
		var msg string
		row.Scan(&msg)
		return msg, nil
	}

	return "", nil
}

func (db *Database) getUser(name string) (User, error) {
	var user User

	err := db.database.Get(&user, "SELECT * from tip_users where github_name=$1", name)

	if (err != nil) && (err.Error() != "sql: no rows in result set") {

		log.Println("select err: ", err)
		return user, err
	}

	return user, nil
}

func (db *Database) getPreset(name string, userID int) (Preset, error) {
	var preset Preset

	err := db.database.Get(&preset, "SELECT * from presets where user_id=$1 AND preset=$2", userID, name)

	if (err != nil) && (err.Error() != "sql: no rows in result set") {

		log.Println("select err: ", err)
		return preset, err
	}

	return preset, nil
}

func (db *Database) createUser(destination string) (User, error) {
	key, err := keypair.Random()
	var user User

	user.GithubName = destination
	user.AccountID = key.Address()
	user.SecretKey = key.Seed()

	sql := fmt.Sprintf("INSERT INTO tip_users (github_name,account_id,secret_key) values ('%s','%s','%s')", user.GithubName, user.AccountID, user.SecretKey)

	//log.Println(sql)

	_, err = db.database.Exec(sql)
	return user, err

}

//////////////////////////////////////
/*

func getCursor() (string,error) {
    database, err := sqlite3.Open(gConf.SQLITE_LOCATION)
    if err != nil {
        return "",err
    }

    defer database.Close()

    row, err := database.Query("SELECT value from Notes where key='cursor'")
    if err != nil {
        fmt.Println("select err: ",err)
        return "",err
    }

    if row != nil {
        defer row.Close()
        row.Next()
        var msg string
        row.Scan(&msg)
        return msg,nil
    }

    return "",nil
}

func storeCursor(cursor string) error {
    sql := fmt.Sprintf("UPDATE Notes set value='%s' where key='start'", cursor)
    return execute(sql)
}

func storePeriod() error {
    return execute("UPDATE Notes set value=CURRENT_TIMESTAMP where key='start'")
}


func getAmountSpentInPeriod() (int64,error) {
    database, err := sqlite3.Open(gConf.SQLITE_LOCATION)
    if err != nil {
        return 0,err
    }

    defer database.Close()

    row, err := database.Query("SELECT sum(amount) from Payments where date > (SELECT value from Notes where key='start')")
    if err != nil {
        fmt.Println("select err: ",err)
        return 0,err
    }

    if row != nil {
        defer row.Close()
        row.Next()
        var msg int64
        row.Scan(&msg)
        return msg,nil
    }

    return 0,nil
}

func storePayment(source string,fullText string,action string,amount int) error {
    sql := fmt.Sprintf("INSERT INTO Payments (source, fullText, action, amount, date) values (%s,%s,%s,%d,CURRENT_TIMESTAMP)", source, fullText, action, amount)
    return execute(sql)
}

func execute(sql string) error {
    database, err := sqlite3.Open(gConf.SQLITE_LOCATION)
    if err != nil {
        return err
    }
    defer database.Close()
    err = database.Exec(sql)
    return err
}
*/
