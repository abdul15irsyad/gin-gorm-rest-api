package main

import (
	"fmt"
	"gin-gorm-rest-api/config"
	"gin-gorm-rest-api/controllers"
	"gin-gorm-rest-api/middlewares"
	"gin-gorm-rest-api/routes"
	"gin-gorm-rest-api/services"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/dig"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading `.env` file: " + err.Error())
	}

	Env := os.Getenv("ENV")
	if Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))
	router.SetTrustedProxies([]string{})
	router.Static("/assets", "./assets")

	// init routes
	allRoutes := InitRoutes(router)
	for _, route := range allRoutes {
		route.Init(router)
	}

	// listen on port
	Port := os.Getenv("PORT")
	log.Println("environment=" + Env)
	router.Run("localhost:" + fmt.Sprint(Port))
}

type Route interface {
	Init(r *gin.Engine)
}

func ProvideMultiple(container *dig.Container, constructors []interface{}) error {
	for _, constructor := range constructors {
		if err := container.Provide(constructor); err != nil {
			return err
		}
	}
	return nil
}

func InitRoutes(router *gin.Engine) []Route {
	container := dig.New()
	if err := ProvideMultiple(container, []interface{}{
		// config
		config.NewDatabaseConfig,
		// middlewares
		middlewares.NewAuthMiddleware,
		// services
		services.NewFileService,
		services.NewJwtService,
		services.NewRoleService,
		services.NewUserService,
		// controllers
		controllers.NewAuthUserController,
		controllers.NewAuthController,
		controllers.NewFileController,
		controllers.NewRoleController,
		controllers.NewRootController,
		controllers.NewUserController,
		// routes
		routes.NewAuthRoute,
		routes.NewFileRoute,
		routes.NewRoleRoute,
		routes.NewRootRoute,
		routes.NewUserRoute,
	}); err != nil {
		panic(err)
	}

	allRoutes := []Route{}

	if err := container.Invoke(func(authRoute *routes.AuthRoute, fileRoute *routes.FileRoute, roleRoute *routes.RoleRoute, rootRoute *routes.RootRoute, userRoute *routes.UserRoute) {
		allRoutes = []Route{authRoute, fileRoute, roleRoute, rootRoute, userRoute}
	}); err != nil {
		panic(err)
	}

	return allRoutes
}
