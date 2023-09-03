package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"simple-douyin-backend/dal/db/dao"
	"simple-douyin-backend/dal/db/gorm_gen"
	"simple-douyin-backend/pkg/constants"
	"time"
)

var DB *gorm.DB
var DSN string = getDSN()

// Init init DB
func Init() {
	var err error
	DB, err = gorm.Open(mysql.Open(DSN),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}
	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	// Automatically update datable tables' structure
	DB.AutoMigrate(&dao.Message{})
	DB.AutoMigrate(&dao.Relation{})
	DB.AutoMigrate(&dao.UserDetail{})
	// Setup dao
	dao.Init(DB)
	// Setup gorm generated code
	gorm_gen.SetDefault(DB)
}

func getDSN() string {
	host, hostok := os.LookupEnv("MYSQL_HOST")
	if !hostok {
		host = constants.MySQLDefaultHost
	}
	port, portok := os.LookupEnv("MYSQL_PORT")
	if !portok {
		port = constants.MySQLDefaultPort
	}
	user, userok := os.LookupEnv("MYSQL_USER")
	if !userok {
		user = constants.MySQLDefaultUser
	}
	pwd, pwdok := os.LookupEnv("MYSQL_PASSWORD")
	if !pwdok {
		pwd = constants.MySQLDefaultPwd
	}
	dbname, dbnameok := os.LookupEnv("MYSQL_DATABASE")
	if !dbnameok {
		dbname = constants.MySQLDefaultDBName
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pwd, host, port, dbname)
	return dsn
}
