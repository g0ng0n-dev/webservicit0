package main

import (
	"github.com/g0ng0n-dev/webservicito/database"
	"log"
	"github.com/g0ng0n-dev/webservicito/product"
	"net/http"
	_"github.com/go-sql-driver/mysql"
)

const apiBasePath = "/api"
func main() {
	database.SetupDatabase()
	product.SetupRoute(apiBasePath)
	http.ListenAndServe(":5000", nil)
}