package database

import (
	"ci-cd-golang/config"
	"ci-cd-golang/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var openConnectionDB *gorm.DB
var err error

func PostgresConnection() (*gorm.DB, error) {
	myDSN := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Vientiane",
		config.GetEnv("postgres.host", "localhost"),
		config.GetEnv("postgres.user", "postgres"),
		config.GetEnv("postgres.password", "123456"),
		config.GetEnv("postgres.database", "naga"),
		config.GetEnv("postgres.port", "5432"),
	)
	//dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	fmt.Println("CONNECTING_TO_POSTGRES_DB")
	openConnectionDB, err = gorm.Open(postgres.Open(myDSN), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),

		NowFunc: func() time.Time {
			ti, _ := time.LoadLocation("Asia/Vientiane")
			return time.Now().In(ti)
		},
	})
	//disable logger mode
	//openConnectionDB.Logger.LogMode(logger.Silent)
	//DryRun: false,
	if err != nil {

		log.Fatal("ERROR_PING_POSTGRES", err)
		return nil, err
	}
	fmt.Println("POSTGRES_CONNECTED")
	migrate(openConnectionDB)
	return openConnectionDB, nil
}

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Migrate Success")
}
