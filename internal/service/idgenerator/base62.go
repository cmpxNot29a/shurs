package idgenerator

import (
	"github.com/cmpxNot29a/shurs/internal/helper"
)

// Base62Generator реализует domain/service.IDGenerator.
type Base62Generator struct {
	length int
}

// NewBase62Generator создает генератор Base62.
func NewBase62Generator(length int) *Base62Generator {
	return &Base62Generator{length: length}
}

// Generate генерирует случайный Base62 ID.
func (g *Base62Generator) Generate() (string, error) {
	randomBytes, err := helper.GenerateRandomBase62(g.length)
	if err != nil {
		return "", err
	}
	return string(randomBytes), nil
}
