package product

import (
	"encoding/json"
	"fmt"
	"github.com/g0ng0n-dev/webservicito/cors"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const productsBasePath = "products"

func SetupRoute(apiBasePath string){
	handleProducts := http.HandlerFunc(productsHandler)
	handleProduct := http.HandlerFunc(productHandler)
	reportHandler := http.HandlerFunc(handleProductReport)
	http.Handle("/websocket", websocket.Handler(productSocket))
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productsBasePath), cors.Middleware(handleProducts))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productsBasePath), cors.Middleware(handleProduct))
	http.Handle(fmt.Sprintf("%s/%s/reports", apiBasePath, productsBasePath), cors.Middleware(reportHandler))
}

func productsHandler(w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case http.MethodGet:
		productList, err := getProductList()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		productsJson, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJson)
	case http.MethodPost:
		// add a new product to the list
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if newProduct.ProductID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = insertProduct(newProduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	case http.MethodOptions:
		return
	}
}

func productHandler(w http.ResponseWriter, r *http.Request){
	urlPathSegments := strings.Split(r.URL.Path, "products/")
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments) -1 ])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	product, err := getProduct(productID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		productJSON, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productJSON)
	case http.MethodPut:
		var updatedProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updatedProduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if updatedProduct.ProductID != productID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		updateProduct(updatedProduct)
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		removeProduct(productID)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}