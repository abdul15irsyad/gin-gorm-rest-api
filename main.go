package main

import (
	"fmt"
	"gin-gorm-rest-api/controllers"
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/middlewares"
	"gin-gorm-rest-api/routes"
	"gin-gorm-rest-api/services"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading `.env` file: " + err.Error())
	}
	database.InitDatabase()

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
	// routes.AuthRoutes(router)
	// routes.UserRoutes(router)
	// routes.RoleRoutes(router)
	// routes.FileRoutes(router)

	// services init
	jwtService := services.NewJwtService()
	userService := services.NewUserService(database.DB)
	roleService := services.NewRoleService(database.DB)
	fileService := services.NewFileService(database.DB)
	// controllers init
	rootController := controllers.NewRootController()
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(jwtService, userService)
	authUserController := controllers.NewAuthUserController(userService)
	roleController := controllers.NewRoleController(roleService)
	fileController := controllers.NewFileController(fileService)
	// middlewares init
	authMiddleware := middlewares.NewAuthMiddleware(jwtService, userService)
	// routes init
	authRoute := routes.NewAuthRoute(authMiddleware, authController, authUserController)
	authRoute.Init(router)
	rootRoute := routes.NewRootRoute(rootController)
	rootRoute.Init(router)
	userRoute := routes.NewUserRoute(authMiddleware, userController)
	userRoute.Init(router)
	roleRoute := routes.NewRoleRoute(authMiddleware, roleController)
	roleRoute.Init(router)
	fileRoute := routes.NewFileRoute(authMiddleware, fileController)
	fileRoute.Init(router)

	// listen on port
	Port := os.Getenv("PORT")
	log.Println("environment=" + Env)
	router.Run("localhost:" + fmt.Sprint(Port))
}
