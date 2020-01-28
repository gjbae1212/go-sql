package gosql

import (
	"fmt"
)

var (
	ErrInvalidParam  = fmt.Errorf("[err][gosql] invalid param")
	ErrNotExistDB    = fmt.Errorf("[err][gosql] not exist db")
	ErrFailConnectDB = fmt.Errorf("[err][gosql] fail connect db")
)

// Connector is common interface for SQL like database.
// This interface is to support methods like connect, close and so on.
type Connector interface {
	DSN() string
	Connect() error
	Close()
}
