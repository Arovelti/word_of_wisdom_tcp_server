package main

import (
	"context"
	"encoding/hex"
	"os"
	"os/signal"
	"syscall"

	"github.com/Arovelti/word_of_wisdom_tcp_server/pkg/env"
	"github.com/Arovelti/word_of_wisdom_tcp_server/pkg/logger"
	"github.com/Arovelti/word_of_wisdom_tcp_server/pow"
	"github.com/Arovelti/word_of_wisdom_tcp_server/service"
	"github.com/sirupsen/logrus"
)

func main() {
	l := logger.Init()
	e := env.InitClientEnv(l)
	w := pow.Init(l, uint8(e.Difficulty))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-ctx.Done()
		cancel()
	}()

	client := service.NewClient(e, l, w)
	for i := 0; i < 100; i++ {
		if ctx.Err() != nil {
			break
		}

		q, err := client.GetQuote(ctx)
		if err != nil {
			l.Error("unable to get quote", logrus.WithError(err))
			continue
		} else {
			l.WithField("quote", q).Info("Logging byte slice as []byte")
			l.WithField("hexData", hex.EncodeToString(q)).Info("Logging byte slice as hexadecimal string")
		}
	}

	l.Info("Successfully done")
}
