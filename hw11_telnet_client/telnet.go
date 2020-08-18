package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	Send() error
	Receive() error
	Close() error
}

type Client struct {
	addr    string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{addr: address, timeout: timeout, in: in, out: out}
}

func (c *Client) Connect() error {
	var err error
	if c.conn, err = net.DialTimeout("tcp", c.addr, c.timeout); err != nil {
		return err
	}
	return nil
}

func (c *Client) Send() error {
	return passData(c.in, c.conn)
}

func (c *Client) Receive() error {
	return passData(c.conn, c.out)
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func passData(in io.Reader, out io.Writer) error {
	buf := make([]byte, 1024)
	n, err := in.Read(buf)
	if err != nil {
		return err
	}

	if _, err = out.Write(buf[:n]); err != nil {
		return err
	}
	return nil
}
