package context

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

type Context struct {
	Db		*sql.DB
}

var DefaultContext *Context

func init() {
	db, err := sql.Open("mysql", "root:@/mystories?charset=utf8")
	if err != nil {
		log.Fatalf("init db failed:%v\n", err)
		return
	}
	DefaultContext = &Context{
		Db: db,
	}
}
