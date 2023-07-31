package service

import (
	"context"
	"fmt"
	"net"

	"github.com/Arovelti/word_of_wisdom_tcp_server/pkg/env"
	"github.com/Arovelti/word_of_wisdom_tcp_server/pow"
	"github.com/sirupsen/logrus"
)

type Client struct {
	env *env.ClientEnv
	log *logrus.Logger
	w   pow.Worker
}

func NewClient(e *env.ClientEnv, l *logrus.Logger, w pow.Worker) *Client {
	return &Client{
		env: e,
		log: l,
		w:   w,
	}
}

func (c *Client) GetQuote(ctx context.Context) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	d := net.Dialer{}
	conn, err := d.DialContext(ctx, c.env.Transport, c.env.ServerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to start dial tcp")
	}
	defer func() {
		if err := conn.Close(); err != nil {
			c.log.Error("failed to close connection: ", err.Error())
		}
	}()

	if err := WriteMsg(conn, []byte("challenge")); err != nil {
		return nil, fmt.Errorf("write message err: %w", err)
	}

	challenge, err := ReadMsg(conn)
	if err != nil {
		return nil, fmt.Errorf("read message with challenge err: %w", err)
	}

	solution, err := c.w.Work(challenge)
	if err != nil {
		return nil, fmt.Errorf("challenge error: %w", err)
	}

	if err := WriteMsg(conn, solution); err != nil {
		return nil, fmt.Errorf("send solution err: %w", err)
	}

	quote, err := ReadMsg(conn)
	if err != nil {
		return nil, fmt.Errorf("read message with quote err: %w", err)
	}

	return quote, nil
}
