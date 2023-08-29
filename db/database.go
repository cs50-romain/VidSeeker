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

func AddYoutuber(name string, youtuberid string, playlistid string) error {
	result, err := db.Exec("INSERT INTO youtuber (channelid, channelname, playlistid) VALUES (?, ?, ?)", youtuberid, name, playlistid)
	if err != nil {
		return fmt.Errorf("addYoutuber: %v", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("youtuber: %v", err)
	}
	return nil
}

func RemoveYoutuber(name string) error {
	_, err := db.Exec("DELETE FROM youtuber WHERE channelname='" + name + "';")	
	if err != nil {
		return fmt.Errorf("RemoveYoutuber: %v", err)
	}

	return nil
}

func YoutuberById(id int) (string, error) {
	var youtuber string
	err := db.QueryRow("SELECT channelname FROM youtuber WHERE id = ?", id).Scan(&youtuber)
	if err != nil {
		return "", err
	}
	return youtuber, nil
}
