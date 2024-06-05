package main

import (
	"net/http"

	"github.com/monari-onsongo/go-mm/config"
)

func main() {
	conf := config.NewConfig("passKey", "consumerKey", "consumerSecret")

	http.HandleFunc("/stkpush", conf.GetSTKPUSH)
	http.ListenAndServe(":8080", nil)
}
