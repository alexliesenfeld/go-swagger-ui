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
		swaggerui.WithTryItOutEnabled(true),
		swaggerui.WithDisplayRequestDuration(true),
		swaggerui.WithDocExpansion(swaggerui.DocExpansionFull),
		swaggerui.WithPersistAuthorization(true),
		swaggerui.WithFilter(true, ""),
		swaggerui.WithSpecURLs("oas-example", []swaggerui.SpecURL{
			{Name: "petstore", URL: "https://petstore.swagger.io/v2/swagger.json"},
			{Name: "oas-example", URL: "https://raw.githubusercontent.com/OAI/OpenAPI-Specification/main/examples/v3.0/api-with-examples.yaml"},
			{Name: "oas-link-example", URL: "https://raw.githubusercontent.com/OAI/OpenAPI-Specification/main/examples/v3.0/link-example.json"},
			{Name: "oas-upsto-example", URL: "https://raw.githubusercontent.com/OAI/OpenAPI-Specification/main/examples/v3.0/uspto.json"},
		}),
	))

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
