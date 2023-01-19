package main

import (
	"net/http"

	"golang/user/userdb"

	"github.com/gorilla/mux"
)

func main() {
	var usrhandlerobj userdb.UsrHandler

	usrhandlerobj.Connection("localhost","postgres","root","forgolang","5433")
	
	router:=mux.NewRouter()
	router.HandleFunc("/health", userdb.HealthCheck).Methods("GET")
	router.HandleFunc("/user", usrhandlerobj.GetUser).Methods("GET")
	router.HandleFunc("/adduser", usrhandlerobj.AddUser).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8400", router)

	dbinstance,_ := usrhandlerobj.DB.DB()
	defer dbinstance.Close()



}