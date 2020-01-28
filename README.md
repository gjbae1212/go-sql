# go-sql
<p align="left">
<a href="https://hits.seeyoufarm.com"/><img src="https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fgjbae1212%2Fgo-sql"/></a>
<a href="https://goreportcard.com/badge/github.com/gjbae1212/go-sql"><img src="https://goreportcard.com/badge/github.com/gjbae1212/go-sql" alt="Go Report Card"/></a>
<a href="/LICENSE"><img src="https://img.shields.io/badge/license-MIT-GREEN.svg" alt="license" /></a> 
</p>
This project is a connector for SQL databases.

## Getting Started
### Mysql
```go
package main

import (
	gomysql "github.com/gjbae1212/go-sql/mysql"
)

func main() {
	conn, err := gomysql.NewConnector("user:password@/dbname", 2)
	if err != nil {
		panic(err)
	}

	if err := conn.Connect(); err != nil {
		panic(err)
	}

	// db is *sqlx.DB(https://github.com/jmoiron/sqlx)
	db, err := conn.DB()
	if err != nil {
		panic(err)
	}

	_, _ = db.Query("SELECT * FROM sample")
	_, _ = db.Queryx("SELECT * FROM sample")
}
```

## To be Supported
- [ ] Sqlite
- [ ] Postgres
- [ ] BigQuery
- [ ] And so on ... 

## License
This project is licensed under the MIT License
