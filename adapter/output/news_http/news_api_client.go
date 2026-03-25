package news_http

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/copier"
	"github.com/willianVini-dev/hexagonal/adapter/output/model/response"
	"github.com/willianVini-dev/hexagonal/application/domain"
	"github.com/willianVini-dev/hexagonal/configuration/env"
	"github.com/willianVini-dev/hexagonal/configuration/rest_err"
)

type newsApiClient struct {
}

func NewNewsApiClient() *newsApiClient {
	client = resty.New().SetBaseURL("https://newsapi.org/v2")
	return &newsApiClient{}
}

var (
	client *resty.Client
)

func (nc *newsApiClient) GetNewsPort(newsDomain domain.NewsRequestDomain) (*domain.NewsDomain, *rest_err.RestErr) {

	newsResponse := &response.NewsClientResponse{}

	fmt.Printf("news domain: %+v\n", newsDomain)
	_, err := client.R().
		SetQueryParams(map[string]string{
			"q":      newsDomain.Subject,
			"from":   newsDomain.From,
			"apiKey": env.GetNewsTokenApi(),
		}).
		SetResult(newsResponse).
		Get("/everything")

	fmt.Printf("news response: %+v\n", newsResponse)
	fmt.Printf("news env: %+v\n", env.GetNewsTokenApi())

	if err != nil {
		return nil, rest_err.NewInternalServerError("error when trying to get news")
	}

	newsDomainResponse := &domain.NewsDomain{}
	copier.Copy(newsDomainResponse, newsResponse)

	return newsDomainResponse, nil
}
