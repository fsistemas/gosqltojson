package database

import (
	"franciscoperez.dev/gosqltojson/internal/config"
	"testing"
)

func TestNewDBConnOK(t *testing.T) {
	configFile := config.ConfigFile{
		Connections: map[string]string{
			"testdata_db": "sqlite+../../testdata/testdata.db",
		},
	}

	_, err := NewDBConn(&configFile, "testdata_db")

	if err != nil {
		t.Log("Expect to get a sqlite connection but got an error", err)
		t.Fail()
	}
}

func TestNewDBConnInvalidConnectionString(t *testing.T) {
	configFile := config.ConfigFile{
		Connections: map[string]string{
			"testdata_db": "../testdata/testdata.db",
		},
	}

	_, err := NewDBConn(&configFile, "testdata_db")

	if err == nil {
		t.Log("Expect error caused by invalid connection string")
		t.Fail()
	}
}

func TestNewDBConnInvalidConnectionType(t *testing.T) {
	configFile := config.ConfigFile{
		Connections: map[string]string{
			"testdata_db": "unknown+../testdata/testdata.db",
		},
	}

	_, err := NewDBConn(&configFile, "testdata_db")

	if err == nil {
		t.Log("Expect error because unknown connection type")
		t.Fail()
	}
}
