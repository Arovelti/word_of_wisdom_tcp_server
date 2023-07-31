package env

import (
	"time"

	"github.com/sirupsen/logrus"
)

var (
	defaultClientServerAddress = ":8080"
	defaultClientKeepAlive     = time.Duration(time.Second * 15)
	defaultDifficulty          = 12
	defaultTransport           = "tcp"
)

type ClientEnv struct {
	ServerAddress string
	Transport     string
	KeepAlive     time.Duration
	Difficulty    uint8
}

func (cli *ClientEnv) NewClientEnv(serverAddress string, transport string, keepAlive time.Duration, difficulty int) *ClientEnv {
	if serverAddress == "" {
		serverAddress = defaultClientServerAddress
	}

	if transport == "" {
		transport = defaultTransport
	}

	if keepAlive == 0 {
		keepAlive = defaultClientKeepAlive
	}

	if difficulty == 0 {
		difficulty = defaultDifficulty
	}

	return &ClientEnv{
		ServerAddress: serverAddress,
		Transport:     transport,
		KeepAlive:     keepAlive,
		Difficulty:    uint8(difficulty),
	}
}

func InitClientEnv(l *logrus.Logger) *ClientEnv {
	if err := LoadEnvFromFile("../.env/client.env"); err != nil {
		l.Errorf("can't load env from file: %v", err)
	}

	serverAddress := MustString("SERVER_ADDRESS")
	transport := MustString("TRANSPORT")
	keepAlive := MustInt("KEEP_ALIVE")
	difficulty := MustInt("DIFFICULTY")

	c := ClientEnv{}
	client := c.NewClientEnv(serverAddress, transport, time.Duration(time.Duration(time.Second*time.Duration(keepAlive))), difficulty)

	return client
}
