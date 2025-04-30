package valueobjects

import (
	"regexp"

	"github.com/gblcarvalho/go-expert-lab-otel/internal/utils"
)

type CEP struct {
	value string
}

func NewCEP(cepStr string) (CEP, error) {
	if match, err := regexp.MatchString(`^\d{8}$`, cepStr); !match || err != nil {
		return CEP{}, utils.ErrInvalidCEP
	}
	return CEP{value: cepStr}, nil
}

func (c *CEP) Value() string {
	return c.value
}
