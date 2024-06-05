package main

import (
	"net/http"

	"github.com/narie-monarie/go-mpesa/config"
)

func main() {
	conf := config.NewConfig("passKey", "consumerKey", "consumerSecret")

	http.HandleFunc("/stkpush", conf.GetSTKPUSH)
	http.ListenAndServe(":8080", nil)
}
