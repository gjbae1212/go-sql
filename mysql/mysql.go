package mysql

import (
	"fmt"
	"sync"
	"time"
	"database/sql"

	"github.com/cenkalti/backoff/v4"
	gosql "github.com/gjbae1212/go-sql"
	"github.com/go-sql-driver/mysql"
	"github.com/luna-duclos/instrumentedsql"
	"github.com/luna-duclos/instrumentedsql/opentracing"
)

const (
	defaultDriver     = "mysql"
	opentracingDriver = "mysql-opentracing"
)

// Connector is connector for mysql.
type Connector interface {
	gosql.Connector
}

type conn struct {
	driverName string
	dsn        string
	db         *sql.DB
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
		driverName: defaultDriver,
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

	c := &conn{
		driverName: opentracingDriver,
		dsn:        dsn,
		tries:      tries,
		backoff:    backoff.NewExponentialBackOff(),
	}

	return c, nil
}

// DriverName returns driver name.
func (c *conn) DriverName() string {
	return c.driverName
}

// DSN returns dns string.
func (c *conn) DSN() string {
	return c.dsn
}

// DB returns *gorm.DB object if db was connected successfully.
func (c *conn) DB() (*sql.DB, error) {
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
		db, err := sql.Open(c.driverName, c.dsn)
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

func init() {
	sql.Register(opentracingDriver, instrumentedsql.WrapDriver(
		&mysql.MySQLDriver{}, instrumentedsql.WithTracer(opentracing.NewTracer(true))))
}
