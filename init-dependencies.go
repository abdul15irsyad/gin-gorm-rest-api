package main

import (
	"gin-gorm-rest-api/config"
	"gin-gorm-rest-api/controllers"
	"gin-gorm-rest-api/middlewares"
	"gin-gorm-rest-api/routes"
	"gin-gorm-rest-api/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func InitDependencies(router *gin.Engine) {
	container := dig.New()
	// config
	if err := container.Provide(config.NewDatabaseConfig); err != nil {
		panic(err)
	}
	// middlewares
	if err := container.Provide(middlewares.NewAuthMiddleware); err != nil {
		panic(err)
	}
	// services
	if err := container.Provide(services.NewFileService); err != nil {
		panic(err)
	}
	if err := container.Provide(services.NewJwtService); err != nil {
		panic(err)
	}
	if err := container.Provide(services.NewRoleService); err != nil {
		panic(err)
	}
	if err := container.Provide(services.NewUserService); err != nil {
		panic(err)
	}
	// controllers
	if err := container.Provide(controllers.NewAuthUserController); err != nil {
		panic(err)
	}
	if err := container.Provide(controllers.NewAuthController); err != nil {
		panic(err)
	}
	if err := container.Provide(controllers.NewFileController); err != nil {
		panic(err)
	}
	if err := container.Provide(controllers.NewRoleController); err != nil {
		panic(err)
	}
	if err := container.Provide(controllers.NewRootController); err != nil {
		panic(err)
	}
	if err := container.Provide(controllers.NewUserController); err != nil {
		panic(err)
	}
	// routes
	if err := container.Provide(routes.NewAuthRoute); err != nil {
		panic(err)
	}
	if err := container.Provide(routes.NewFileRoute); err != nil {
		panic(err)
	}
	if err := container.Provide(routes.NewRoleRoute); err != nil {
		panic(err)
	}
	if err := container.Provide(routes.NewRootRoute); err != nil {
		panic(err)
	}
	if err := container.Provide(routes.NewUserRoute); err != nil {
		panic(err)
	}

	if err := container.Invoke(func(authRoute *routes.AuthRoute, fileRoute *routes.FileRoute, roleRoute *routes.RoleRoute, rootRoute *routes.RootRoute, userRoute *routes.UserRoute) {
		authRoute.Init(router)
		fileRoute.Init(router)
		roleRoute.Init(router)
		rootRoute.Init(router)
		userRoute.Init(router)
	}); err != nil {
		panic(err)
	}
}
