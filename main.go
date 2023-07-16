package main

import (
	"github.com/7leven7/xm-company-crud/app/configs"
	"github.com/7leven7/xm-company-crud/app/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	err := configs.PGConnection()
	if err != nil {
		panic(err)
	}

	configs.Migrate()

	router := gin.Default()

	routes.CompanyRoutes(router)
	routes.UserRoutes(router)

	err = router.Run()
	if err != nil {
		panic(err)
	}
}
