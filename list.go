package main

import (
	"fmt"
)

var ErrNameNotExists = fmt.Errorf("name not exists")

type List map[int]string

func (l List) GetName(spaces int) (string, error) {
	val, exists := l[spaces]
	if !exists {
		return "", ErrNameNotExists
	}

	return val, nil
}
