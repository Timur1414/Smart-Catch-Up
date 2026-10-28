package domain

import "errors"

var (
	ErrNothingInTable = errors.New("no rows in result set")
	ErrWrongLLMAnswer = errors.New("wrong LLM answer")
)
