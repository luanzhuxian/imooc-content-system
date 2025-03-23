package api

import (
	"imooc-content-system/internal/controllers"
	"imooc-content-system/internal/repository"
	"imooc-content-system/internal/services"

	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	rootPath   = "/api/"
	noAuthPath = "/out/api/"
)

func CmsRouters(r *gin.Engine) {
	cmsApp := services.NewCmsApp()
	authMiddleware := &SessionAuthMiddleware{rdb: cmsApp.GetRedisClient()}

	accountRepository := repository.NewAccountRepository(cmsApp.GetDB())
	accountService := services.NewAccountService(accountRepository)
	accountController := controllers.NewAccountController(accountService)
	root := r.Group(rootPath).Use(authMiddleware.Auth)
	{
		root.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		root.GET("/hello", cmsApp.Hello)

		root.POST("/cms/content/find", cmsApp.ContentFind)
		root.POST("/cms/content/create", cmsApp.ContentCreate)
		root.POST("/cms/content/update", cmsApp.ContentUpdate)
		root.DELETE("/cms/content/delete", cmsApp.ContentDelete)
	}
	noAuth := r.Group(noAuthPath)
	{
		noAuth.POST("/cms/register", cmsApp.Register)
		noAuth.POST("/cms/login", cmsApp.Login)
		noAuth.POST("/cms/register2", accountController.Register)
		noAuth.GET("/cms/user/:user_id", accountController.FindByUserID)
	}

}
