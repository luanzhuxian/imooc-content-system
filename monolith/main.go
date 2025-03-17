package main

import (
	"fmt"
	"imooc-content-system/internal/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	api.CmsRouters(r)

	err := r.Run()
	if err != nil {
		fmt.Printf("error = %v\n", err)
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
