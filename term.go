package main

import "fmt"

type Term struct {
	Code   string
	Option string
}

func ParseTerm(code string) (*Term, error) {
	options := map[string]string{
		"FA23": "5320:::FA23",
	}

	option, ok := options[code]
	if !ok {
		return nil, fmt.Errorf("invalid term code")
	}

	return &Term{
		Code:   code,
		Option: option,
	}, nil
}
