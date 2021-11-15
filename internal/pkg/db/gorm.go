package db

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

var MasterDB *gorm.DB

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/lin-cms-test?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.account"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.addr"),
		viper.GetString("mysql.port"),
	)
	var err error
	MasterDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   "lin_",
		},
		// 打印所有执行的SQL
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("%v", err)
		return nil, err
	}
	return MasterDB, nil
}
