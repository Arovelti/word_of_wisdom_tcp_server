package pow

import (
	"github.com/Arovelti/word_of_wisdom_tcp_server/pow/hashcash"
	"github.com/sirupsen/logrus"
)

type PoW interface {
	Prover
	Worker
}

type Prover interface {
	GenerateBytesBuffer(difficulty uint8) ([]byte, error)
	Prove(data, nonce []byte) error
}

type Worker interface {
	Work(data []byte) ([]byte, error)
}

func Init(l *logrus.Logger, difficulty uint8) PoW {
	p, err := hashcash.New(l, difficulty)
	if err != nil {
		l.Fatal("error while initiating proof of work")
	}

	return p
}
