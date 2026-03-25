package input

import (
	"github.com/willianVini-dev/hexagonal/application/domain"
	"github.com/willianVini-dev/hexagonal/configuration/rest_err"
)

type NewsUseCase interface {
	GetNewsService(domain.NewsRequestDomain) (*domain.NewsDomain, *rest_err.RestErr)
}
