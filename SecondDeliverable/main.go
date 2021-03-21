package main

import (
	"log"
	"main/router"
	"net/http"
)

func main() {
	r := router.InitRouter()
	log.Fatal(http.ListenAndServe(":9000", r))
}
