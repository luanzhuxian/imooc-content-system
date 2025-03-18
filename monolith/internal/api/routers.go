package api

import (
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
	session := &SessionAuth{rdb: cmsApp.GetRedisClient()}
	root := r.Group(rootPath).Use(session.Auth)
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
	}

}
