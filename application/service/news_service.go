package service

import (
	"fmt"

	"github.com/willianVini-dev/hexagonal/application/domain"
	"github.com/willianVini-dev/hexagonal/application/port/output"
	"github.com/willianVini-dev/hexagonal/configuration/logger"
	"github.com/willianVini-dev/hexagonal/configuration/rest_err"
)

type newsService struct {
	newsPort output.NewsPort
}

func NewNewsService(newsPort output.NewsPort) *newsService {
	return &newsService{
		newsPort: newsPort,
	}
}

func (ns *newsService) GetNewsService(newRequest domain.NewsRequestDomain) (*domain.NewsDomain, *rest_err.RestErr) {

	logger.Info(
		fmt.Sprintf(
			"init get news service, subject=%s, from=%s", newRequest.Subject, newRequest.From))

	newsDomainResponse, err := ns.newsPort.GetNewsPort(newRequest)
	return newsDomainResponse, err
}
