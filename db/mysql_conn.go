package db


import (
	"github.com/xormplus/xorm"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

var db *xorm.Engine


func init(){
	var err error
	/*
	TCP using default port (3306) on localhost:
	user:password@tcp/dbname?charset=utf8mb4,utf8
	 */
	db, err = xorm.NewMySQL("mysql","gotest:111111a@tcp(localhost:3306)/test?charset=utf8")

	if err!=nil{
		log.Fatal("connect mysql fail !",err)
	}else{
		db.ShowSQL(true)
		db.RegisterSqlMap(xorm.Xml("../dao/sqlmap", ".xml"))
		log.Printf("connect mysql successful!")
	}
}


func GetConn()*xorm.Engine{
	return db
}

func CloseDb(){
	err := db.Close()
	if err!= nil{
		log.Printf("close mysql conn fail")
	}
}
