package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Arovelti/word_of_wisdom_tcp_server/pkg/env"
	"github.com/Arovelti/word_of_wisdom_tcp_server/pow"
	"github.com/Arovelti/word_of_wisdom_tcp_server/repository"
	"github.com/sirupsen/logrus"
)

type Server struct {
	env     *env.ServerEnv
	l        *logrus.Logger
	p        pow.Prover
	repo     repository.Repository
	listener net.Listener
	wg       sync.WaitGroup
	cancel   context.CancelFunc
}

func NewServer(conf *env.ServerEnv, l *logrus.Logger, repo repository.Repository, p pow.Prover) *Server {
	return &Server{
		env: conf,
		l:    l,
		p:    p,
		repo: repo,
	}
}

func (s *Server) handleConnection(conn net.Conn) error {
	defer func() {
		if err := conn.Close(); err != nil {
			s.l.Error("failed to close connection", err)
		}
	}()

	if err := conn.SetDeadline(time.Now().Add(s.env.Deadline)); err != nil {
		return fmt.Errorf("failed to set deadline timeout: %w", err)
	}

	if _, err := ReadMsg(conn); err != nil {
		return fmt.Errorf("")
	}

	challenge, err := s.p.GenerateBytesBuffer(uint8(s.env.Difficulty))
	if err != nil {
		return fmt.Errorf("failed to generate bytes for prove: %w", err)
	}
	if err := WriteMsg(conn, challenge); err != nil {
		return fmt.Errorf("send challenge err: %w", err)
	}

	solution, err := ReadMsg(conn)
	if err != nil {
		return fmt.Errorf("receive proof err: %w", err)
	}

	if err = s.p.Prove(challenge, solution); err != nil {
		return fmt.Errorf("invalid solution: %w", err)
	}

	quote, err := s.repo.GetQuote()
	if err != nil {
		return fmt.Errorf("get random quote err: %w", err)
	}

	if err = WriteMsg(conn, []byte(quote)); err != nil {
		return fmt.Errorf("send quote err: %w", err)
	}

	return nil
}

func (s *Server) serve(ctx context.Context) {
	defer s.wg.Done()

	go func() {
		<-ctx.Done()
		err := s.listener.Close()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			s.l.Error("failed to close listener: ", err)
		}
	}()

	for {
		conn, err := s.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			s.l.Debug("listener closed")
			return
		} else if err != nil {
			s.l.Error("failed to accept connection: ", err)
			continue
		}

		s.wg.Add(1)
		go func(conn net.Conn) {
			defer s.wg.Done()

			if err := s.handleConnection(conn); err != nil {
				s.l.Error("handle connection error: ", err)
			}
		}(conn)
	}
}

func (s *Server) Run(ctx context.Context) (err error) {
	ctx, s.cancel = context.WithCancel(ctx)
	defer s.cancel()

	lc := net.ListenConfig{
		KeepAlive: s.env.KeepAlive,
	}
	if s.listener, err = lc.Listen(ctx, s.env.Transport, s.env.Address); err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.l.Info("starting server at: ", s.listener.Addr().String())

	s.wg.Add(1)
	go s.serve(ctx)
	s.wg.Wait()

	s.l.Info("server stopped")

	return nil
}

func (s *Server) Stop() {
	s.cancel()
}
