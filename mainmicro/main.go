package main

import (
	"net/http"

	"example.com/mainmicro/userdb"

	resty "github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
)

func main() {
	var usrhandlerobj userdb.UsrHandler
	usrhandlerobj.Connection("localhost","postgres","root","forgolang","5433")
	router:=mux.NewRouter()
	router.HandleFunc("/health", userdb.HealthCheck).Methods("GET")
	router.HandleFunc("/user", usrhandlerobj.GetUser).Methods("GET")
	router.HandleFunc("/adduser", usrhandlerobj.AddUser).Methods("POST")
	router.HandleFunc("/invHome", InvHome).Methods("GET")

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
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(resp.Body())
	w.Write([]byte(resp.Body()))
	// w.Write([]byte(resp.String()))
	// w.Write(resp.Body())
}