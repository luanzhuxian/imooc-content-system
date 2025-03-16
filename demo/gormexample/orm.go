package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Account struct {
	ID       int64     `gorm:"column:id;primary_key"`
	UserID   string    `gorm:"column:user_id"`
	Password string    `gorm:"column:password"`
	Nickname string    `gorm:"column:nickname"`
	Ct       time.Time `gorm:"column:created_at"`
	Ut       time.Time `gorm:"column:updated_at"`
}

// 重写表名 可做动态分表
func (a Account) TableName() string {
	table := "account"
	return table
}

func main() {
	db := connDB()
	var accounts []Account
	// 查询所有，不指定条件，保存到accounts切片中
	if err := db.Find(&accounts).Error; err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(accounts)

	var account Account
	if err := db.Where("id = ?", 2).First(&account).Error; err != nil {
		fmt.Println(err)
		return
	}
	// 分页查询
	// db.Where().Offset().Limit().Find()

	// db.Find()
	// db.Delete()
	// db.Update()

	fmt.Println(account)
}

func connDB() *gorm.DB {
	mysqlDB, err := gorm.Open(mysql.Open("root:root123@tcp(localhost:3306)/cms_account?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	db, err := mysqlDB.DB()
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(2)
	mysqlDB = mysqlDB.Debug() // 打印sql语句
	return mysqlDB
}
