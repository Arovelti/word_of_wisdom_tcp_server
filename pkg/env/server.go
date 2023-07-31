package env

import (
	"time"

	"github.com/sirupsen/logrus"
)

var (
	defaultServerAddress   = ":8080"
	defaultServerKeepAlive = time.Duration(time.Second * 15)
	defaultServerDeadline  = time.Duration(time.Second * 10)
)

type ServerEnv struct {
	Address    string
	Transport  string
	KeepAlive  time.Duration
	Deadline   time.Duration
	Difficulty int
}

func (srv *ServerEnv) NewServerEnv(address, transport string, keepAlive, deadline time.Duration, difficulty int) *ServerEnv {
	if address == "" {
		address = defaultServerAddress
	}

	if transport == "" {
		transport = defaultTransport
	}

	if keepAlive == 0 {
		keepAlive = defaultServerKeepAlive
	}

	if deadline == 0 {
		deadline = defaultServerDeadline
	}

	if difficulty == 0 {
		difficulty = defaultDifficulty
	}

	return &ServerEnv{
		Address:    address,
		Transport:  transport,
		KeepAlive:  keepAlive,
		Deadline:   deadline,
		Difficulty: difficulty,
	}
}

func InitServerEnv(l *logrus.Logger) *ServerEnv {
	if err := LoadEnvFromFile("../.env/client.env"); err != nil {
		l.Errorf("can't load env from file: %v", err)
	}

	address := MustString("ADDRESS")
	transport := MustString("TRANSPORT")
	keepAlive := MustInt("KEEP_ALIVE")
	deadline := MustInt("DEADLINE")
	difficulty := MustInt("DIFFICULTY")

	keepAliveDuration := time.Duration(time.Second * time.Duration(keepAlive))
	deadlineDuration := time.Duration(time.Second * time.Duration(deadline))

	s := ServerEnv{}
	return s.NewServerEnv(address, transport, time.Duration(keepAliveDuration), time.Duration(deadlineDuration), difficulty)

}
