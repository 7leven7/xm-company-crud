package routes

import (
	"github.com/7leven7/xm-company-crud/app/controllers"
	middlewares "github.com/7leven7/xm-company-crud/app/middleware"
	"github.com/gin-gonic/gin"
)

func CompanyRoutes(router *gin.Engine) {
	router.GET("/company/:id", controllers.GetCompanyByID)

	authGroup := router.Group("/company")

	authGroup.Use(middlewares.JWTAuthMiddleware())

	authGroup.POST("", controllers.CreateCompany)
	authGroup.PATCH("/:id", controllers.UpdateCompany)
	authGroup.DELETE("/:id", controllers.DeleteCompany)
}

func UserRoutes(router *gin.Engine) {
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
}
