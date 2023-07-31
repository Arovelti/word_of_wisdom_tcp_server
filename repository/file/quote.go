package file

import (
	_ "embed"

	"bufio"
	"bytes"
	"math/rand"
	"strings"
)

//go:embed quotes.txt
var quotes []byte

type quoteRepo struct {
	quotes []string
}

func NewQuote() *quoteRepo {
	qr := new(quoteRepo)

	reader := bytes.NewReader(quotes)
	s := bufio.NewScanner(reader)

	for s.Scan() {
		if q := strings.TrimSpace(s.Text()); q != "" {
			qr.quotes = append(qr.quotes, q)
		}
	}

	return qr
}

func (q *quoteRepo) GetQuote() (string, error) {
	return q.quotes[rand.Intn(len(q.quotes))], nil
}