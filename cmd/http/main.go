// @title REPLACE_APP_NAME API
// @version 1.0
// @description REPLACE_APP_NAME API
// @BasePath /example
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"backend-pipeline-security/internal/boot"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := boot.HTTP(); err != nil {
		log.Println("[HTTP] failed to boot http server due to " + err.Error())
	}
}
