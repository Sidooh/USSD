package datastore

import (
	"os"
	"testing"
)

func initializeDB() {
	err := os.Setenv("APP_ENV", "TEST")
	if err != nil {
		return
	}

	Init()
}

func TestInit(t *testing.T) {
	if db != nil {
		t.Errorf("Init() = %v; want nil", db)
	}

	initializeDB()

	row := db.QueryRow("SELECT * FROM sessions")
	if row.Err() != nil {
		t.Errorf("Init() = %v; want nil", row.Err())
	}

}

func TestUnmarshalFromDatabase(t *testing.T) {
	initializeDB()

	session := Session{}
	err := UnmarshalFromDatabase("a", &session)
	if err == nil {
		t.Errorf("UnmarshalFromDatabase() = %v; want err", session)
	}

	_, err = db.Exec(`INSERT INTO 
								sessions(session_id, phone, text, service_code, status, product, screen_path, vars) 
								VALUES ('a', '1', '1', 'a', 'a', '1', '{}', '{}')`)
	if err != nil {
		t.Errorf("UnmarshalFromDatabase() = %v; expect insert", err)
	}

	err = UnmarshalFromDatabase("a", &session)
	if err != nil {
		t.Errorf("UnmarshalFromDatabase() = %v; want data", err)
	}

	if session.SessionId != "a" {
		t.Errorf("UnmarshalFromDatabase() = %v; want 'a'", session.SessionId)
	}
}

func TestMarshalToDatabase(t *testing.T) {
	initializeDB()

	session := Session{
		SessionId: "a",
		Vars:      []byte{},
	}
	err := MarshalToDatabase(session)
	if err != nil {
		t.Errorf("MarshalToDatabase() = %v; want nil", err)
	}

	sess := Session{}
	err = db.QueryRow("SELECT session_id FROM sessions").Scan(&sess.SessionId)
	if err != nil {
		t.Errorf("MarshalToDatabase() = %v; want nil", err)
	}
}

func TestFetchSessionLogs(t *testing.T) {
	initializeDB()

	logs, err := FetchSessionLogs()
	if err != nil {
		t.Errorf("FetchSessionLogs() = %v; want nil", err)
	}
	if len(logs) > 0 {
		t.Errorf("FetchSessionLogs() = %v; want []", logs)
	}

	_, _ = db.Exec(`INSERT INTO 
								sessions(id, session_id, phone, text, service_code, status, product, screen_path, vars) 
								VALUES (1, 'a', '1', '1', 'a', 'a', '1', '{}', '{}'), 
								       (2, 'b', '1', '1', 'a', 'a', '1', '{}', '{}')`)

	logs, err = FetchSessionLogs()
	if err != nil {
		t.Errorf("FetchSessionLogs() = %v; want nil", err)
	}
	if len(logs) != 2 {
		t.Errorf("FetchSessionLogs() = %v; want 2", len(logs))
	}

	if logs[0].SessionId != "b" && logs[1].SessionId != "a" {
		t.Errorf("FetchSessionLogs() = %v; wrong data", len(logs))
	}
}
