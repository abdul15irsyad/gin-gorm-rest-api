package main

import (
	"fmt"
	"gin-gorm-rest-api/configs"
	"gin-gorm-rest-api/handlers"
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
		AllowOrigins:     []string{"http://localhost:3002"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowWildcard:    true,
		AllowCredentials: true,
	}))
	router.SetTrustedProxies([]string{})
	router.Static("/assets", "./assets")

	// init routes
	allRoutes, logService, logMiddleware := InitRoutes(router)
	logService.Init()
	router.Use(logMiddleware.Log)
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

func InitRoutes(router *gin.Engine) ([]Route, *services.LogService, *middlewares.LogMiddleware) {
	container := dig.New()
	if err := ProvideMultiple(container, []interface{}{
		// configs
		configs.NewDatabaseConfig,
		// middlewares
		middlewares.NewAuthMiddleware,
		middlewares.NewLogMiddleware,
		// services
		services.NewFileService,
		services.NewJwtService,
		services.NewMailService,
		services.NewRoleService,
		services.NewUserService,
		services.NewLogService,
		// handlers
		handlers.NewAuthUserHandler,
		handlers.NewAuthHandler,
		handlers.NewFileHandler,
		handlers.NewRoleHandler,
		handlers.NewRootHandler,
		handlers.NewUserHandler,
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
	var logService *services.LogService
	var logMiddleware *middlewares.LogMiddleware

	if err := container.Invoke(func(
		authRoute *routes.AuthRoute,
		fileRoute *routes.FileRoute,
		roleRoute *routes.RoleRoute, rootRoute *routes.RootRoute, userRoute *routes.UserRoute, ls *services.LogService, lm *middlewares.LogMiddleware) {
		allRoutes = []Route{authRoute, fileRoute, roleRoute, rootRoute, userRoute}
		logService = ls
		logMiddleware = lm
	}); err != nil {
		panic(err)
	}

	return allRoutes, logService, logMiddleware
}
