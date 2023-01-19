package main

import (
	"net/http"

	"golang/inventory/inventorydb"

	"github.com/gorilla/mux"
)

func main() {
	var invhandlerobj inventorydb.InvHandler

	invhandlerobj.Connection("localhost","postgres","root","forgolang","5433")
	
	router:=mux.NewRouter()
	router.HandleFunc("/health", inventorydb.HealthCheck).Methods("GET")
	router.HandleFunc("/inventory", invhandlerobj.GetInventory).Methods("GET")
	router.HandleFunc("/addinventory", invhandlerobj.AddInventory).Methods("POST")
	router.HandleFunc("/delinventory/{id}", invhandlerobj.DelInventory).Methods("DELETE")

	http.Handle("/", router)
	http.ListenAndServe(":8200", router)

	dbinstance,_ := invhandlerobj.DB.DB()
	defer dbinstance.Close()



}