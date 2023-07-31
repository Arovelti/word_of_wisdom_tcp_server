package repository

import "github.com/Arovelti/word_of_wisdom_tcp_server/repository/file"

type Quotes interface {
	GetQuote() (string, error)
}

type Repository struct {
	Quotes
}

func New() Repository {
	return Repository{
		Quotes: file.NewQuote(),
	}
}
