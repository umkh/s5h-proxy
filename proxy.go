package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/armon/go-socks5"
	"github.com/sethvargo/go-envconfig"
)

type config struct {
	Port         uint   `env:"PORT, default=1080"`
	HTTPProxyURL string `env:"HTTP_PROXY_URL, default=http://proxy2.hamkor.local:8085"`
}

func httpProxyDialer(proxyURL *url.URL, timeout time.Duration) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		dialer := &net.Dialer{
			Timeout:   timeout,
			KeepAlive: 30 * time.Second,
		}

		proxyConn, err := dialer.DialContext(ctx, "tcp", proxyURL.Host)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to proxy: %w", err)
		}

		connectReq := &http.Request{
			Method: "CONNECT",
			URL:    &url.URL{Opaque: addr},
			Host:   addr,
			Header: make(http.Header),
		}

		if proxyURL.User != nil {
			if password, ok := proxyURL.User.Password(); ok {
				connectReq.SetBasicAuth(proxyURL.User.Username(), password)
			}
		}

		if err := connectReq.Write(proxyConn); err != nil {
			proxyConn.Close()
			return nil, fmt.Errorf("failed to send CONNECT: %w", err)
		}

		br := bufio.NewReader(proxyConn)
		resp, err := http.ReadResponse(br, connectReq)
		if err != nil {
			proxyConn.Close()
			return nil, fmt.Errorf("failed to read proxy response: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			proxyConn.Close()
			return nil, fmt.Errorf("proxy returned status %d: %s", resp.StatusCode, resp.Status)
		}

		return proxyConn, nil
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var cfg config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		panic(err)
	}

	upstreamURL, err := url.Parse(cfg.HTTPProxyURL)
	if err != nil {
		fmt.Printf("Error parsing upstream proxy URL: %v\n", err)
		log.Fatal(err)
	}

	fmt.Printf("Upstream proxy configured: %s\n", upstreamURL.Redacted())

	conf := &socks5.Config{
		Dial: httpProxyDialer(upstreamURL, 30*time.Second),
	}

	server, err := socks5.New(conf)
	if err != nil {
		fmt.Printf("Error creating SOCKS5 server: %v\n", err)
		log.Fatal(err)
	}

	fmt.Printf("SOCKS5 proxy started and listening on port :%d\n", cfg.Port)
	fmt.Println("Use this address to connect")
	go func() {
		if err := server.ListenAndServe("tcp", fmt.Sprintf(":%d", cfg.Port)); err != nil {
			fmt.Printf("Error starting SOCKS5 server: %v\n", err)
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	fmt.Println("Shutdown")
}
