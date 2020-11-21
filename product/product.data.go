package product

import (
	"fmt"
	"github.com/g0ng0n-dev/webservicito/database"
	"sort"
	"sync"
)

/* We are using a mutex because our Webservices are multithreaded and maps and go are inherently not thread safe.
 * which means we need to wrap our map using a mutex to avoid 2 threads from writing and reading the map at the same time.
*/
var productMap = struct {
	sync.RWMutex
	m map[int]Product

}{m: make(map[int]Product)}

func getProduct(productID int) *Product {
	productMap.RLock()
	defer productMap.RUnlock()
	if product, ok := productMap.m[productID]; ok {
		return &product
	}
	return nil
}

func removeProduct(productID int){
	productMap.Lock()
	defer productMap.Unlock()
	delete(productMap.m, productID)
}

func getProductList() ([]Product, error) {
	results, err := database.DbConn.Query(`SELECT productId, 
	manufacturer, 
	sku, 
	upc, 
	pricePerUnit, 
	quantityOnHand,
	productName
	FROM inventorydb.products;`)

	if err != nil {
		fmt.Println("Error On DB", err)
		return nil, err
	}

	defer results.Close()
	products := make([]Product, 0)

	for results.Next() {
		var product Product
		results.Scan(&product.ProductID,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName)
		products = append(products, product)
	}

	return products, nil
}

func getProductsIds() []int {
	productMap.RLock()
	productIds := []int{}
	for key := range productMap.m {
		productIds = append(productIds, key)
	}
	productMap.RUnlock()
	sort.Ints(productIds)
	return productIds
}

func getNextProductID() int {
	productIDs := getProductsIds()
	return productIDs[len(productIDs)-1] + 1
}

func addOrUpdateFunc(product Product) (int, error) {
	// if the product id is set, update, otherwise add
	addOrUpdateID := -1
	if product.ProductID > 0 {
		oldProduct := getProduct(product.ProductID)
		// if it exist, replace it, otherwise return error
		if oldProduct == nil {
			return 0, fmt.Errorf("product id [%d] doesn't exist ", product.ProductID)
		}
		addOrUpdateID = product.ProductID
	} else {
		addOrUpdateID = getNextProductID()
		product.ProductID = addOrUpdateID
	}

	productMap.Lock()
	productMap.m[addOrUpdateID] = product
	productMap.Unlock()
	return addOrUpdateID, nil
}