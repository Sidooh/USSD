package datastore

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "modernc.org/sqlite"
	"strings"
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

	env := strings.ToUpper(viper.GetString("APP_ENV"))
	if env == "TEST" {
		driverName = "sqlite"
	}

	conn, err := sql.Open(driverName, viper.GetString("DB_DSN"))
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
	err := db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Closed Database")
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

func FetchSessionLogs(limit int) ([]SessionLog, error) {
	var sessions []SessionLog

	rows, err := db.Query(`SELECT * FROM sessions ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		fmt.Println(err)
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

func FetchSessionLog(id int) (SessionLog, error) {
	var session SessionLog

	row := db.QueryRow(`SELECT * FROM sessions WHERE id = ?`, id)

	var screenPath []byte
	var vars []byte

	// get RawBytes from data
	err := row.Scan(&session.Id, &session.SessionId, &session.Phone, &session.Text, &session.Code, &session.Status,
		&session.Product, &screenPath, &vars, &session.CreatedAt, &session.UpdatedAt)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	err = json.Unmarshal(screenPath, &session.ScreenPath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(vars, &session.Vars)
	if err != nil {
		panic(err)
	}

	if err = row.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return session, nil
}

func ReadTimeSeriesCount() (interface{}, error) {
	type Dataset struct {
		Date  int `json:"date"`
		Count int `json:"count"`
	}

	var datasets []Dataset

	rows, err := db.Query(
		`SELECT DATE_FORMAT(created_at, '%Y%m%d%H') as date, COUNT(id) as count 
				FROM sessions GROUP BY date ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		dataset := new(Dataset)

		if err := rows.Scan(&dataset.Date, &dataset.Count); err != nil {
			log.Fatal(err)
		}

		datasets = append(datasets, *dataset)
	}

	return datasets, nil
}

func ReadSummaries() (interface{}, error) {
	var sessions struct {
		Today int `json:"today"`
		Total int `json:"total"`
	}
	now := time.Now().UTC()
	today := fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day())

	rows, err := db.Query(`SELECT SUM(created_at > ?) as today, COUNT(created_at) as total FROM sessions`, today)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&sessions.Today, &sessions.Total); err != nil {
			log.Fatal(err)
		}
	}

	return sessions, nil
}
