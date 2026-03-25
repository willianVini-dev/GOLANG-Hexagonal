package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/willianVini-dev/hexagonal/adapter/input/controller"
	"github.com/willianVini-dev/hexagonal/adapter/output/news_http"
	"github.com/willianVini-dev/hexagonal/application/service"
)

func InitRoutes(r *gin.Engine) {

	newsApiClient := news_http.NewNewsApiClient()
	newsService := service.NewNewsService(newsApiClient)
	newsController := controller.NewNewsController(newsService)

	r.GET("/news", newsController.GetNews)

}
