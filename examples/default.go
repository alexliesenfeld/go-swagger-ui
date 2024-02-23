package main

import (
	"fmt"
	swaggerui "github.com/alexliesenfeld/go-swagger-ui"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", swaggerui.NewHandler(
		swaggerui.WithHTMLTitle("My Example Petstore API"),
		swaggerui.WithSpecURL("https://petstore.swagger.io/v2/swagger.json"),
		swaggerui.WithDocExpansion(swaggerui.DocExpansionFull),
	))

	fmt.Println("Starting Swagger UI on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
