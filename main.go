package main

import (
	"fmt"
	"gin-gorm-rest-api/handlers"
	"gin-gorm-rest-api/lib"
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
		panic("error loading `.env` file: " + err.Error())
	}

	Env := os.Getenv("ENV")
	gin.SetMode(gin.ReleaseMode)
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
	allRoutes, allMiddlewares := InitRoutes(router)
	router.Use(allMiddlewares.logMiddleware.Handler)
	router.Use(allMiddlewares.errorMiddleware.Handler)
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

func ProvideMultiple(container *dig.Container, constructors []any) error {
	for _, constructor := range constructors {
		if err := container.Provide(constructor); err != nil {
			return err
		}
	}
	return nil
}

type AllMiddleware struct {
	logMiddleware   *middlewares.LogMiddleware
	errorMiddleware *middlewares.ErrorMiddleware
}

func InitRoutes(router *gin.Engine) ([]Route, AllMiddleware) {
	container := dig.New()
	if err := ProvideMultiple(container, []any{
		// lib
		lib.NewLogger,
		lib.NewDatabase,
		// middlewares
		middlewares.NewAuthMiddleware,
		middlewares.NewLogMiddleware,
		middlewares.NewErrorMiddleware,
		// services
		services.NewFileService,
		services.NewJwtService,
		services.NewMailService,
		services.NewRoleService,
		services.NewUserService,
		services.NewTokenService,
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
	var allMiddlewares AllMiddleware

	if err := container.Invoke(func(
		authRoute *routes.AuthRoute,
		fileRoute *routes.FileRoute,
		roleRoute *routes.RoleRoute,
		rootRoute *routes.RootRoute,
		userRoute *routes.UserRoute,
		em *middlewares.ErrorMiddleware,
		lm *middlewares.LogMiddleware) {
		allRoutes = []Route{authRoute, fileRoute, roleRoute, rootRoute, userRoute}
		allMiddlewares.logMiddleware = lm
		allMiddlewares.errorMiddleware = em
	}); err != nil {
		panic(err)
	}

	return allRoutes, allMiddlewares
}
