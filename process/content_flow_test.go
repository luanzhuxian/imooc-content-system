package process

//
//import (
//	"gorm.io/driver/mysql"
//	"gorm.io/gorm"
//	"testing"
//)
//
//func TestExecContentFlow(t *testing.T) {
//	err := ExecContentFlow(func() *gorm.DB {
//		mysqlDB, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/?charset=utf8mb4&parseTime=True&loc=Local"))
//		if err != nil {
//			panic(err)
//		}
//		db, err := mysqlDB.DB()
//		if err != nil {
//			panic(err)
//		}
//		db.SetMaxOpenConns(4)
//		db.SetMaxIdleConns(2)
//		return mysqlDB
//	}())
//	if err != nil {
//		panic(err)
//	}
//}
