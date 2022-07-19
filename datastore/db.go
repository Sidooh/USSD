package datastore

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
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

type SessionLog struct {
	Id         string                 `json:"id"`
	SessionId  string                 `json:"session_id"`
	Phone      string                 `json:"phone"`
	Text       string                 `json:"text"`
	Code       string                 `json:"code"`
	Status     string                 `json:"status"`
	Product    int                    `json:"product"`
	ScreenPath map[string]interface{} `json:"screen_path"`
	Vars       map[string]string      `json:"vars"`
	CreatedAt  *time.Time             `json:"created_at,omitempty"`
	UpdatedAt  *time.Time             `json:"updated_at,omitempty"`
}

func Init() {
	fmt.Println("Initializing USSD subsystem database")

	driverName := "mysql"

	env := os.Getenv("APP_ENV")
	if env == "TEST" {
		driverName = "sqlite"
	}

	conn, err := sql.Open(driverName, os.Getenv("DB_DSN"))
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

func FetchSessionLogs() ([]SessionLog, error) {
	var sessions []SessionLog

	rows, err := db.Query(`SELECT * FROM sessions ORDER BY id DESC LIMIT 50`)
	if err != nil {
		return nil, err
	}

	var counter = 0

	//TODO: Move to repo or something so that we can unmarshall to ScreenPath struct
	// Fetch rows
	for rows.Next() {
		session := new(SessionLog)

		var screenPath []byte
		var vars []byte

		// get RawBytes from data
		err = rows.Scan(&session.Id, &session.SessionId, &session.Phone, &session.Text, &session.Code, &session.Status,
			&session.Product, &screenPath, &vars, &session.CreatedAt, &session.UpdatedAt)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		err := json.Unmarshal(screenPath, &session.ScreenPath)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(vars, &session.Vars)
		if err != nil {
			panic(err)
		}

		sessions = append(sessions, *session)

		counter += 1
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return sessions, nil
}
