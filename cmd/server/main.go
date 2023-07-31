package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Arovelti/word_of_wisdom_tcp_server/pkg/env"
	"github.com/Arovelti/word_of_wisdom_tcp_server/pkg/logger"
	"github.com/Arovelti/word_of_wisdom_tcp_server/pow"
	"github.com/Arovelti/word_of_wisdom_tcp_server/repository"
	"github.com/Arovelti/word_of_wisdom_tcp_server/service"
)

func main() {
	l := logger.Init()
	e := env.InitServerEnv(l)
	repo := repository.New()
	p := pow.Init(l, uint8(e.Difficulty))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-ctx.Done()
		cancel()
	}()

	server := service.NewServer(e, l, repo, p)	
	if err := server.Run(ctx); err != nil {
		l.Fatal("server run error: ", err)
	}

	l.Info("server started")
}
