package main

import (
	"log"
	"main/controller"
	"main/router"
	csvservice "main/service/csv"
	httpservice "main/service/http"
	"main/usecase"
	"net/http"
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
