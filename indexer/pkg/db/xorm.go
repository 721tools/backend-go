package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var xormEngine *xorm.Engine

func NewDBEngine(dsn string) {
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		panic(fmt.Errorf("new xorm failed, err is %v", err))
	}
	err = engine.Ping()
	if err != nil {
		panic(fmt.Errorf("ping failed"))
	}
	// engine.ShowSQL(true)
	// engine.Logger().SetLevel(log.LOG_DEBUG)
	xormEngine = engine
}

func GetDBEngine() *xorm.Engine {
	return xormEngine
}

func NewSession() *xorm.Session {
	return xormEngine.NewSession()
}
