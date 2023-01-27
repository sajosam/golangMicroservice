package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"example.com/mainmicro/userdb"

	resty "github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Inventory struct {
	Id int `json:"ID"`
	ProductName string `json:"product_name"`
	ProductQuantity int `json:"product_quantity"`
}

type Order struct {
	gorm.Model
	Id int `json:"ID"`
	Product string `json:"product"`
	Quantity int `json:"quantity"`
	User_name string `json:"user_name"`
}


var usrhandlerobj userdb.UsrHandler
func main() {
	// usrhandlerobj.Connection("host.docker.internal","postgres","root","forgolang","5433")
	usrhandlerobj.Connection("localhost","postgres","root","forgolang","5433")

	router:=mux.NewRouter()
	router.HandleFunc("/health", userdb.HealthCheck).Methods("GET")
	router.HandleFunc("/user", usrhandlerobj.GetUser).Methods("GET")
	router.HandleFunc("/adduser", usrhandlerobj.AddUser).Methods("POST")
	router.HandleFunc("/delUser/{id}", usrhandlerobj.DelUser).Methods("DELETE")
	router.HandleFunc("/ordHome", OrdHome).Methods("GET")
	router.HandleFunc("/invHome", InvHome).Methods("GET")
	router.HandleFunc("/proHome", ProHome).Methods("GET")
	router.HandleFunc("/pro/{id}", SinglePro).Methods("GET")
	router.HandleFunc("/delInv/{id}", DelInv).Methods("DELETE")
	router.HandleFunc("/addInv", addInv).Methods("POST")
	router.HandleFunc("/addPro", addPro).Methods("POST")
	router.HandleFunc("/addOrd", addOrd).Methods("POST")


	http.Handle("/", router)
	http.ListenAndServe(":8400", router)

	dbinstance,_ := usrhandlerobj.DB.DB()
	defer dbinstance.Close()
}

func InvHome(w http.ResponseWriter, r *http.Request) {
	client := resty.New()
	resp, err := client.R().Get("http://localhost:8200/inventory")
	// print the values in the response
	if err != nil {
		panic(err)
	}
	w.Write([]byte(resp.Body()))
}

func OrdHome(w http.ResponseWriter, r *http.Request) {
	client := resty.New()
	resp, err := client.R().Get("http://localhost:8300/order")
	// print the values in the response
	if err != nil {
		panic(err)
	}
	w.Write([]byte(resp.Body()))
}

func ProHome(w http.ResponseWriter, r *http.Request) {
	client := resty.New()
	resp, err := client.R().Get("http://localhost:8100/product")
	// print the values in the response
	if err != nil {
		panic(err)
	}
	w.Write([]byte(resp.Body()))
}

func SinglePro(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	client := resty.New()
	resp, err := client.R().Get("http://localhost:8100/product/"+vars["id"])
	// print the values in the response
	if err != nil {
		panic(err)
	}
	w.Write([]byte(resp.Body()))
}

func DelInv(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	client := resty.New()
	resp, err := client.R().Delete("http://localhost:8200/delinventory/"+vars["id"])
	// print the values in the response
	if err != nil {
		panic(err)
	}
	w.Write([]byte(resp.Body()))
}

func addInv(w http.ResponseWriter, r *http.Request) {
	client := resty.New()
	w.Header().Set("Content-Type", "application/json")
	// var inventory Inventory
	d,_:=ioutil.ReadAll(r.Body)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(string(d)).
		Post("http://localhost:8200/addinventory")
	// print the values in the response
	if err != nil {
		panic(err)
	}
	w.Write([]byte(resp.Body()))


}

func addPro(w http.ResponseWriter, r *http.Request) {
	client := resty.New()
	d,_:=ioutil.ReadAll(r.Body)
	// resp, err := client.R().Post("http://localhost:8100/addproduct")
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(string(d)).
		Post("http://localhost:8100/addproduct")
	// print the values in the response
	if err != nil {
		panic(err)
	}
	w.Write([]byte(resp.Body()))
}

func addOrd(w http.ResponseWriter, r *http.Request) {

	
	// get data from body
	w.Header().Set("Content-Type", "application/json")

	// get data from json
	var order Order
	d,_:=ioutil.ReadAll(r.Body)
	json.Unmarshal(d, &order)
	pr_name:=order.Product

	// get product quantity from inventory table

	var inventory []Inventory
	usrhandlerobj.DB.Where("product_name = ?", pr_name).Find(&inventory)
	inv_qty:=inventory[0].ProductQuantity
	ord_qty:=order.Quantity

	fmt.Println(inv_qty)
	fmt.Println(ord_qty)
	fmt.Println(inventory[0].ProductQuantity)
	fmt.Println(inventory)
	if ord_qty > inv_qty {
		w.Write([]byte("Order quantity is greater than inventory quantity"))
	} else {
		// post data to order table
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(string(d)).
			Post("http://localhost:8300/addorder")
		// print the values in the response
		if err != nil {
			panic(err)
		}
		w.Write([]byte(resp.Body()))
		// update inventory table
		inv_qty=inv_qty-ord_qty
		inventory[0].ProductQuantity=inv_qty
		usrhandlerobj.DB.Save(&inventory[0])
	}
}




// 	usrhandlerobj.DB.Find(&inventory)
// 	fmt.Println(inventory)
// 	// json.NewEncoder(w).Encode(inventory)

// }