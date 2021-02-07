package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
)

var db *gorm.DB
var err error

// Product is a representation of a product
type Product struct {
	ID    int             `json:"id"`
	Code  string          `json:"code"`
	Name  string          `json:"name"`
	Price decimal.Decimal `json:"price" sql:"type:decimal(16,2)"`
}

// Result is an array of product
type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	db, err = gorm.Open("mysql", "root:1223334444@/go_rest_api_crud?charset=utf8&parseTime=True")

	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection estabilished")
	}

	db.AutoMigrate((&Product{}))
	handleRequests()

}

func handleRequests() {
	log.Println("Start the development server at http://127.0.0.1:9999")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/api/products", createProduct).Methods("POST")
	myRouter.HandleFunc("/api/products", getProducts).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", getProductByID).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", updateProductByID).Methods("PUT")

	log.Fatal(http.ListenAndServe(":9999", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome Oppa!")
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "ini create product")

	// baca seluruh isi bodynya
	payloads, _ := ioutil.ReadAll(r.Body)

	// deklarasi var utk nampung data product yg mau di insert
	var product Product
	json.Unmarshal(payloads, &product)

	// insert data
	db.Create(&product)

	// respon nya
	res := Result{Code: 200, Data: product, Message: "Succcess create product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

func getProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{}

	db.Find(&products)

	res := Result{Code: 200, Data: products, Message: "Success get products"}
	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

func getProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Product
	db.First(&product, productID)

	res := Result{Code: 200, Data: product, Message: "Success get products"}
	resultOne, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultOne)
}

func updateProductByID(w http.ResponseWriter, r *http.Request) {
	// get ID product nya
	vars := mux.Vars(r)
	productID := vars["id"]

	// baca smua data body nya
	payloads, _ := ioutil.ReadAll(r.Body)

	//tampung data product yg mau di update
	var productUpdate Product
	json.Unmarshal(payloads, &productUpdate)

	// get data product by ID S
	var product Product
	db.First(&product, productID)

	// update datanya
	db.Model(&product).Updates(productUpdate)

	res := Result{Code: 200, Data: product, Message: "Succcess update product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}
