package main

import (
	"github.com/g0ng0n-dev/webservicito/database"
	"github.com/g0ng0n-dev/webservicito/product"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

const apiBasePath = "/api"
func main() {
	database.SetupDatabase()
	product.SetupRoute(apiBasePath)
	http.ListenAndServe(":5000", nil)
}