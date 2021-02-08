package main

import (
	"corpus"
	"net/http"
	"pkg"
)

func main() {
	pkg.InitLog("corpus")
	pkg.Logi("poem server started")

	http.HandleFunc("/authors", pkg.Cors(corpus.ListAuthors))
	http.HandleFunc("/poems/seged", pkg.Cors(corpus.ListSegedPoems))

	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		panic(err.Error())
	}
}
