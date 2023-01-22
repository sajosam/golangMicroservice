package main

import (
	"net/http"

	"example.com/mainmicro/userdb"

	resty "github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
)

func main() {
	var usrhandlerobj userdb.UsrHandler
	usrhandlerobj.Connection("host.docker.internal","postgres","root","forgolang","5433")
	router:=mux.NewRouter()
	router.HandleFunc("/health", userdb.HealthCheck).Methods("GET")
	router.HandleFunc("/user", usrhandlerobj.GetUser).Methods("GET")
	router.HandleFunc("/adduser", usrhandlerobj.AddUser).Methods("POST")
	router.HandleFunc("/ordHome", OrdHome).Methods("GET")
	router.HandleFunc("/invHome", InvHome).Methods("GET")
	router.HandleFunc("/proHome", ProHome).Methods("GET")
	router.HandleFunc("/pro/{id}", SinglePro).Methods("GET")
	router.HandleFunc("/delInv/{id}", DelInv).Methods("DELETE")
	router.HandleFunc("/addInv", addInv).Methods("POST")
	router.HandleFunc("/addPro", addPro).Methods("POST")


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
	resp, err := client.R().Post("http://localhost:8200/addinventory")
	// print the values in the response
	if err != nil {
		panic(err)
	}
	w.Write([]byte(resp.Body()))
}

func addPro(w http.ResponseWriter, r *http.Request) {
	client := resty.New()
	resp, err := client.R().Post("http://localhost:8100/addproduct")
	// print the values in the response
	if err != nil {
		panic(err)
	}
	w.Write([]byte(resp.Body()))
}