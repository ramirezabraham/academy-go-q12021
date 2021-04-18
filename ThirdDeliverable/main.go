package main

import (
	"log"
	"net/http"

	"main/controller"
	"main/router"
	"main/usecase"

	csvservice "main/service/csv"
	httpservice "main/service/http"
)

func main() {

	csvService := csvservice.New()
	httpService := httpservice.New()
	useCase := usecase.New(csvService, httpService)
	mainController := controller.New(useCase)

	mainRouter := router.New(mainController)
	r := mainRouter.InitRouter()
	log.Fatal(http.ListenAndServe(":9000", r))
}
