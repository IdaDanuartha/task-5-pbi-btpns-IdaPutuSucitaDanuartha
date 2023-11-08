package main

import (
	"btpn-syariah-final-project/database"
	"btpn-syariah-final-project/router"
)

func init() {
	database.LoadEnvVariables()
	database.ConnectToDb()
	database.SyncDatabase()
}

func main() {
	
	r := router.SetupRouter()
	
	r.Run()
}

