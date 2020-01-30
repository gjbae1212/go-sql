package mysql

import (
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	gosql "github.com/gjbae1212/go-sql"
	"github.com/stretchr/testify/assert"
)

func TestNewConnector(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		dns   string
		tries int
		isErr bool
	}{
		"fail":    {isErr: true},
		"success": {dns: "test", tries: 1},
	}

	for _, t := range tests {
		_, err := NewConnector(t.dns, t.tries)
		assert.Equal(t.isErr, err != nil)
	}
}

func TestNewConnectorWithOpentracing(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		dns   string
		tries int
		isErr bool
	}{
		"fail":    {isErr: true},
		"success": {dns: "test", tries: 1},
	}

	for _, t := range tests {
		_, err := NewConnectorWithOpentracing(t.dns, t.tries)
		assert.Equal(t.isErr, err != nil)
	}
}

func TestConn_DriverName(t *testing.T) {
	assert := assert.New(t)

	conn1, _ := NewConnector("allan-test", 1)
	conn2, _ := NewConnectorWithOpentracing("allan-test", 1)
	tests := map[string]struct {
		conn   Connector
		output string
	}{
		"default":     {conn: conn1, output: defaultDriver},
		"opentracing": {conn: conn2, output: opentracingDriver},
	}

	for _, t := range tests {
		assert.Equal(t.output, t.conn.DriverName())
	}
}

func TestConn_DSN(t *testing.T) {
	assert := assert.New(t)

	conn, _ := NewConnector("allan-test", 1)
	tests := map[string]struct {
		conn   Connector
		output string
	}{
		"success": {conn: conn, output: "allan-test"},
	}

	for _, t := range tests {
		assert.Equal(t.output, t.conn.DSN())
	}
}

func TestConn_DB(t *testing.T) {
	assert := assert.New(t)

	conn1, _ := NewConnector("allan-test", 1)
	conn2, _ := NewConnector("allan-test", 1)
	db, _, _ := sqlmock.New()
	conn2.(*conn).db = db

	tests := map[string]struct {
		conn Connector
		err  error
	}{
		"fail": {
			conn: conn1,
			err:  gosql.ErrNotExistDB,
		},
		"success": {
			conn: conn2,
			err:  nil,
		},
	}

	for _, t := range tests {
		_, err := t.conn.DB()
		assert.True(errors.Is(err, t.err))
	}
}

func TestConn_Connect(t *testing.T) {
	assert := assert.New(t)

	conn, _ := NewConnector("empty-test", 1)
	tests := map[string]struct {
		conn   Connector
		output error
	}{
		"fail": {conn: conn, output: gosql.ErrFailConnectDB},
	}

	for _, t := range tests {
		err := t.conn.Connect()
		assert.True(errors.Is(err, t.output))
	}
}

func TestConn_Close(t *testing.T) {
	assert := assert.New(t)

	conn, _ := NewConnector("empty-test", 1)
	tests := map[string]struct {
		conn Connector
	}{
		"success": {conn: conn},
	}

	for _, t := range tests {
		t.conn.Close()
		_ = assert
	}
}
