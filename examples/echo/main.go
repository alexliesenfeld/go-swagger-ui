package main

import (
	"fmt"
	swaggerui "github.com/alexliesenfeld/go-swagger-ui"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// This is required to serve Swagger UI files
	e.GET("/swagger-ui/*", echo.WrapHandler(swaggerui.NewHandler(
		swaggerui.WithHTMLTitle("My Example Petstore API"),
		swaggerui.WithDocExpansion(swaggerui.DocExpansionFull),
		swaggerui.WithSpecURL("https://petstore.swagger.io/v2/swagger.json"),
	)))

	fmt.Println("Starting Swagger UI on http://localhost:1323/swagger-ui/")
	e.Logger.Fatal(e.Start(":1323"))
}
