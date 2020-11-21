package main

import (
	"github.com/g0ng0n-dev/webservicito/product"
	"net/http"
)

const apiBasePath = "/api"
func main() {

	product.SetupRoute(apiBasePath)
	http.ListenAndServe(":5000", nil)
}