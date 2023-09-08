package db 

import (
	"database/sql"
	//"os"
	"fmt"
	
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Connect() error{
	cfg := mysql.Config{
		User:   "username",//os.Getenv("DBUSER01"),
		Passwd: "password",//os.Getenv("DBPASS01"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "vidseeker",
		AllowNativePasswords: true,
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
		//log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return pingErr
		//log.Fatal(pingErr)
	}
	return nil
}

func Close() {
	if db != nil {
		db.Close()
		fmt.Println("Closed")
	}
}

func AddUser(username string, password string) error {
	result, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		return fmt.Errorf("adduser: %v", err)
	}
	_, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("user: %v", err)
	}
	return nil
}

func AddYoutuber(name string, youtuberid string, playlistid string, user_id int) error {
	result, err := db.Exec("INSERT INTO youtuber (channelid, channelname, playlistid, user_id) VALUES (?, ?, ?, ?)", youtuberid, name, playlistid, user_id)
	if err != nil {
		return fmt.Errorf("addYoutuber: %v", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("youtuber: %v", err)
	}
	return nil
}

func RemoveYoutuber(name string, user_id int) error {
	_, err := db.Exec("DELETE FROM youtuber WHERE channelname = ? AND user_id = ?", name, user_id)
	if err != nil {
		return fmt.Errorf("RemoveYoutuber: %v", err)
	}

	return nil
}

func YoutuberById(id int, user_id int) (string, error) {
	var youtuber string
	err := db.QueryRow("SELECT channelname FROM youtuber WHERE id = ? AND user_id = ?", id, user_id).Scan(&youtuber)
	if err != nil {
		return "", err
	}
	return youtuber, nil
}

func UserById(username string, password string) (string, string, int, error) {
	var dbuser, dbpass string
	var user_id int
	err := db.QueryRow("SELECT username,password,id FROM users WHERE username = ? AND password = ?", username, password).Scan(&dbuser, &dbpass, &user_id)
	if err != nil {
		return "", "", -1, err
	}
	return dbuser, dbpass, user_id, nil
}
