package mysql

import (
	"errors"
	"testing"

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
		assert.Equal(t.output, conn.DSN())
	}
}

func TestConn_DB(t *testing.T) {
	assert := assert.New(t)

	conn, _ := NewConnector("allan-test", 1)
	tests := map[string]struct {
		conn Connector
		err  error
	}{
		"fail": {
			conn: conn,
			err:  gosql.ErrNotExistDB,
		},
	}

	for _, t := range tests {
		_, err := conn.DB()
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
		err := conn.Connect()
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
