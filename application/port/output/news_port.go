package output

import (
	"github.com/willianVini-dev/hexagonal/application/domain"
	"github.com/willianVini-dev/hexagonal/configuration/rest_err"
)

type NewsPort interface {
	GetNewsPort(domain.NewsRequestDomain) (*domain.NewsDomain, *rest_err.RestErr)
}
