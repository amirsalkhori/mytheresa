package app

import (
	"fmt"
	"log"
	"mytheresa/configs"
	mysql "mytheresa/internal/infra/db/mysql/config"
)

func StartApplication() {
	cfg := configs.GetConfig()

	mysqlRepo, err := mysql.NewMySQLRepository(&cfg)
	if err != nil {
		log.Fatal("Error connecting to MySQL:", err)
	}

	fmt.Println("DB", mysqlRepo.DB)
}
