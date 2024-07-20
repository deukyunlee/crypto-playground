package ethClient

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"net/http"
	"time"
)

type headerTransport struct {
	Transport http.RoundTripper
	headers   map[string][]string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, values := range t.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return t.Transport.RoundTrip(req)
}

func GetEthClient() (*ethclient.Client, context.Context) {
	headers := map[string][]string{
		"Content-Type": {"application/json"},
		"User-Agent":   {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36"},
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	httpClient.Transport = &headerTransport{
		Transport: transport,
		headers:   headers,
	}

	ctx := context.Background()

	client, err := rpc.DialHTTPWithClient("https://api.roninchain.com/rpc", httpClient)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	ethClient := ethclient.NewClient(client)
	return ethClient, ctx
}
