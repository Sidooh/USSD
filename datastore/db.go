package datastore

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

var (
	db *sql.DB
)

type Session struct {
	Id         string
	SessionId  string
	Phone      string
	Text       string
	Code       string
	Status     string
	Product    int
	ScreenPath string
	Vars       []byte
}

func Init() {
	fmt.Println("Initializing USSD subsystem database")

	conn, err := sql.Open("mysql", os.Getenv("DB_DSN"))
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	err = conn.Ping()
	if err != nil {
		panic(err) // proper error handling instead of panic in your app
	}

	_, err = conn.Exec(`create table if not exists sessions
	(
	    id           bigint unsigned auto_increment primary key,
	    
	    session_id   varchar(191) not null unique,
	    phone        varchar(15)  not null,
	    text         varchar(191) not null,
	    service_code varchar(15)  not null,
	    
	    status 		 varchar(32)  not null,
	    product		 varchar(8)   not null,
	    screen_path  text 		  not null,
	    vars 		 text 		  not null,
	    
	    created_at   timestamp    null,
	    updated_at   timestamp    null
	)`)
	if err != nil {
		panic(err)
	}

	db = conn
}

func Close() {
	fmt.Println("Closing USSD subsystem database")

	db.Close()
}

func UnmarshalFromDatabase(sessionId string, session *Session) error {
	stmtOut, err := db.Prepare(`SELECT session_id, status, product, screen_path, vars 
										FROM sessions WHERE session_id = ?`)
	if err != nil {
		return err
	}
	defer stmtOut.Close()

	err = stmtOut.QueryRow(sessionId).Scan(&session.SessionId, &session.Status, &session.Product, &session.ScreenPath, &session.Vars)
	if err != nil {
		return err
	}

	return nil
}

func MarshalToDatabase(session Session) error {
	// Check if exists
	dbSession := new(Session)
	err := UnmarshalFromDatabase(session.SessionId, dbSession)
	if err != nil {
		//	Insert
		stmtIns, err := db.Prepare(`INSERT INTO sessions(session_id, phone, text, service_code, status, product, screen_path, vars, created_at, updated_at) 
										VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )`)
		if err != nil {
			return err
		}
		defer stmtIns.Close()

		_, err = stmtIns.Exec(session.SessionId, session.Phone, session.Text, session.Code, session.Status, session.Product, session.ScreenPath, session.Vars, time.Now(), time.Now())
		if err != nil {
			return err
		}

		return nil
	}

	// Update
	stmtUpd, err := db.Prepare(`UPDATE sessions 
										SET text = ?, status = ?, product = ?, screen_path = ?, vars = ?, updated_at = ?
										WHERE session_id = ?`)
	if err != nil {
		return err
	}
	defer stmtUpd.Close()

	_, err = stmtUpd.Exec(session.Text, session.Status, session.Product, session.ScreenPath, session.Vars, time.Now(), session.SessionId)
	if err != nil {
		return err
	}

	return nil
}
