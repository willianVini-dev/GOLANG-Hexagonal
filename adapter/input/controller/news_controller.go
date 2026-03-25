package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/willianVini-dev/hexagonal/adapter/input/model/request"
	"github.com/willianVini-dev/hexagonal/application/domain"
	"github.com/willianVini-dev/hexagonal/application/port/input"
	"github.com/willianVini-dev/hexagonal/configuration/logger"
	"github.com/willianVini-dev/hexagonal/configuration/validation"
)

type newsController struct {
	newsUseCase input.NewsUseCase
}

func NewNewsController(newsUseCase input.NewsUseCase) *newsController {
	return &newsController{
		newsUseCase: newsUseCase,
	}
}

func (nc *newsController) GetNews(c *gin.Context) {

	logger.Info("Getting news...")
	newRequest := request.NewsRequest{}

	if err := c.ShouldBindQuery(&newRequest); err != nil {

		logger.Error("Error trying to bind query: ", err)
		errRest := validation.ValidateUseError(err)
		c.JSON(errRest.Code, errRest)
		return

	}

	newsDomain := domain.NewsRequestDomain{
		Subject: newRequest.Subject,
		From:    newRequest.From.Format("2006-01-02"),
	}
	newsResponseDomain, err := nc.newsUseCase.GetNewsService(newsDomain)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, newsResponseDomain)
}
