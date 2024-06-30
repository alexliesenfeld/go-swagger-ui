package main

import (
	"fmt"
	swaggerui "github.com/alexliesenfeld/go-swagger-ui"
	"log"
	"net/http"
)

func main() {

	swaggerUIHandler := swaggerui.NewHandler(
		swaggerui.WithSpecURL("https://petstore.swagger.io/v2/swagger.json"),

		// This is only required if you want to serve Swagger UI on a path without a trailing slash,
		// as well, e.g., "/swagger-ui" and not only on "/swagger-ui/" (notice the trailing slash).
		// If the handler receives a request to "/swagger-ui" it will respond just fine, but the browser
		// will be confused and request assets from path "/", which will not work.
		// If we set the base path to "/swagger", it will allow the handler to receive
		// requests on path "/swagger-ui" (without a trailing slash).and "/swagger-ui/"
		// (with trailing slash).
		// You do not need to set a base path at all if you don't mind a trailing slash.
		swaggerui.WithBasePath("/swagger"),
	)

	// This is required to serve Swagger UI files (e.g., HTML, CSS, JavaScript, etc.)
	http.HandleFunc("/swagger/", swaggerUIHandler)

	// This is additionally required if you want Swagger UI to be available
	// on path "/swagger" as well, not only on "/swagger/" or "/swagger/index.html"
	http.HandleFunc("/swagger", swaggerUIHandler)

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
