# S5H-Proxy

A SOCKS5 to HTTP proxy bridge written in Go. This proxy accepts SOCKS5 connections and forwards them through an HTTP CONNECT proxy.

## How It Works

1. **SOCKS5 Server**: The application starts a SOCKS5 proxy server listening on a configurable port (default: 1080).

2. **Custom Dialer**: When a client connects via SOCKS5 and requests a connection to a target server, the proxy uses a custom dialer (`httpProxyDialer`) instead of direct TCP connections.

3. **HTTP CONNECT Tunneling**: The custom dialer:
   - Connects to the upstream HTTP proxy
   - Sends an HTTP CONNECT request to establish a tunnel to the target server
   - Handles HTTP proxy authentication if credentials are provided in the proxy URL
   - Returns the established connection to the SOCKS5 server

4. **Data Forwarding**: Once the tunnel is established, the SOCKS5 server transparently forwards data between the client and target server through the HTTP proxy tunnel.

## Configuration

The proxy is configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Port for the SOCKS5 server to listen on | `1080` |
| `HTTP_PROXY_URL` | Upstream HTTP proxy URL (supports auth: `http://user:pass@host:port`) | `http://proxy2.hamkor.local:8085` |


## Installation

### Prerequisites

- Go 1.24.0 or higher.

### Build

```bash
go build -o s5h-proxy proxy.go
```

## Usage

### Basic Usage

```bash
./s5h-proxy
```

This starts the proxy with default settings (listening on port 1080, using the default upstream proxy).

### Custom Configuration

```bash
PORT=8080 HTTP_PROXY_URL=http://myproxy.example.com:3128 ./s5h-proxy
```

### With Proxy Authentication

```bash
HTTP_PROXY_URL=http://username:password@proxy.example.com:8080 ./s5h-proxy
```

### Using Docker (Example)

```bash
docker run -e PORT=1080 -e HTTP_PROXY_URL=http://proxy.example.com:8080 -p 1080:1080 s5h-proxy
```

## Troubleshooting

### Connection Refused
- Verify the SOCKS5 port is not in use: `lsof -i :1080`
- Check firewall settings

### Proxy Connection Failed
- Verify upstream HTTP proxy URL is correct
- Test upstream proxy connectivity: `curl -x $HTTP_PROXY_URL https://example.com`
- Check proxy authentication credentials

## License

This project uses the following open-source libraries:
- go-socks5 (Mozilla Public License 2.0)
- go-envconfig (Apache License 2.0)

## Contributing

To contribute to this project:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## Support

For issues and questions, please open an issue on the project repository.
