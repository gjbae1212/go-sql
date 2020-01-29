package mysql

import (
	"fmt"
	"sync"
	"time"
	"database/sql"

	"github.com/cenkalti/backoff/v4"
	gosql "github.com/gjbae1212/go-sql"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/luna-duclos/instrumentedsql"
	"github.com/luna-duclos/instrumentedsql/opentracing"
)

// Connector is connector for mysql.
type Connector interface {
	gosql.Connector
	DB() (*sqlx.DB, error)
}

type conn struct {
	driverName string
	dsn        string
	db         *sqlx.DB
	tries      int
	backoff    *backoff.ExponentialBackOff
	lock       sync.RWMutex
}

// NewConnector returns the mysql-connector.
func NewConnector(dsn string, tries int) (Connector, error) {
	if dsn == "" || tries <= 0 {
		return nil, fmt.Errorf("%w mysql.NewConnector", gosql.ErrInvalidParam)
	}

	c := &conn{
		driverName: "mysql",
		dsn:        dsn,
		tries:      tries,
		backoff:    backoff.NewExponentialBackOff(),
	}

	return c, nil
}

// NewConnectorWithOpentracing returns the mysql-connector hooked opentracing.
func NewConnectorWithOpentracing(dsn string, tries int) (Connector, error) {
	if dsn == "" || tries <= 0 {
		return nil, fmt.Errorf("%w mysql.NewConnectorWithOpentracing", gosql.ErrInvalidParam)
	}

	driverName := "mysql-opentracing"
	sql.Register(driverName, instrumentedsql.WrapDriver(
		&mysql.MySQLDriver{}, instrumentedsql.WithTracer(opentracing.NewTracer(true))))

	c := &conn{
		driverName: driverName,
		dsn:        dsn,
		tries:      tries,
		backoff:    backoff.NewExponentialBackOff(),
	}

	return c, nil
}

// DSN returns dns string.
func (c *conn) DSN() string {
	return c.dsn
}

// DB returns sqlx.DB object if db was connected successfully.
func (c *conn) DB() (*sqlx.DB, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.db == nil {
		return nil, fmt.Errorf("%w mysql.DB", gosql.ErrNotExistDB)
	}

	return c.db, nil
}

// Connect is to connect to DB
// if a connection fails, retrying to connect as much as a count of retries with exponential interval.
func (c *conn) Connect() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// if DB is already connected, closing a current DB.
	if c.db != nil {
		return nil
	}

	// try to connect
	defer c.backoff.Reset()
	for i := 0; i < c.tries; i++ {
		time.Sleep(c.backoff.NextBackOff())
		db, err := sqlx.Connect(c.driverName, c.dsn)
		if err != nil {
			continue
		}
		c.db = db
		return nil
	}

	return fmt.Errorf("%w mysql.Connect", gosql.ErrFailConnectDB)
}

// Close is to close a connection for DB.
func (c *conn) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.db != nil {
		c.db.Close()
		c.db = nil
	}
}
