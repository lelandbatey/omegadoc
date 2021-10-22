package postprocess

import (
	"github.com/lelandbatey/omegadoc/domain"
)

var pprocessors []domain.Postprocessor

func RegisterPostprocessor(p domain.Postprocessor) {
	pprocessors = append(pprocessors, p)
}

func GetPostprocessors() []domain.Postprocessor {
	tmp := make([]domain.Postprocessor, len(pprocessors))
	copy(tmp, pprocessors)
	return tmp
}
