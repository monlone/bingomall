package main

import (
	"bingomall/helpers"
	"bingomall/routers"
	"bingomall/system"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"net/http"
	"os"
	"time"
)

// @title goodCorn API 文档
// @version 1.0
// @description 这是 goodCorn 应用 swagger 的示例
// @contact.name kris
// @host 0.0.0.0:8088
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		helper.ErrorLogger.Errorln("Error loading .env file：", err)
	}
	runMode := os.Getenv("RUN_MODE")
	ginConfig := system.GetGinConfig()
	gin.SetMode(runMode)
	router := gin.New()
	//配置跨域
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:9528"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Gin-Access-Token", "X-CSRF-Token", "Authorization", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		AllowWildcard:    true,
		MaxAge:           12 * time.Hour,
	}))

	router.HandleMethodNotAllowed = ginConfig.HandleMethodNotAllowed
	router.Use(system.Logger(helper.AccessLogger), gin.Recovery())
	router.Static("/page", "view")
	router.MaxMultipartMemory = ginConfig.MaxMultipartMemory
	routers.RegisterApiRoutes(router)
	routers.RegisterAppRoutes(router)
	routers.RegisterOpenRoutes(router)
	routers.RegisterAuthRoutes(router)
	serverConfig := system.GetServerConfig()
	server := &http.Server{
		Addr:           serverConfig.Addr,
		IdleTimeout:    serverConfig.IdleTimeout * time.Second,
		ReadTimeout:    serverConfig.ReadTimeout * time.Second,
		WriteTimeout:   serverConfig.WriteTimeout * time.Second,
		MaxHeaderBytes: serverConfig.MaxHeaderBytes,
		Handler:        router,
	}
	_ = server.ListenAndServe()
}

func init() {
	// 先读取服务配置文件
	err := system.LoadServerConfig("conf/server-config.yml")
	if err != nil {
		helper.ErrorLogger.Errorln("读取服务配置错误：", err)
		os.Exit(3)
	}
}
