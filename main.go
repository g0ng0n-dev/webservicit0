package main

import (
	"github.com/g0ng0n-dev/webservicito/database"
	"github.com/g0ng0n-dev/webservicito/product"
	"github.com/g0ng0n-dev/webservicito/receipt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

const apiBasePath = "/api"
func main() {
	database.SetupDatabase()
	receipt.SetupRoute(apiBasePath)
	product.SetupRoute(apiBasePath)
	http.ListenAndServe(":5000", nil)
}