package main

import (
	"gateway/dispatcher"
	"gateway/model"
	"net/http"
	"pkg"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	pkg.InitLog("gateway")

	model.InitUserDb()

	pkg.Logi("gateway server started")
	// dispense requests
	http.HandleFunc("/user/", pkg.Cors(dispatcher.DispenseUser))
	http.HandleFunc("/corpus/", pkg.Cors(dispatcher.DispenseCorpus))

	err := http.ListenAndServe(":13000", nil)

	if err != nil {
		panic(err.Error())
	}
}
