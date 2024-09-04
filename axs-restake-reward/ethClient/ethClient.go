package ethClient

import (
	"context"
	"github.com/deukyunlee/crypto-playground/logging"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"net/http"
	"time"
)

var (
	logger = logging.GetLogger()
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

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	httpClient.Transport = &headerTransport{
		Transport: transport,
	}
	ctx := context.Background()

	client, err := rpc.DialOptions(ctx, "https://api.roninchain.com/rpc", rpc.WithHTTPClient(httpClient))

	if err != nil {
		logger.Errorf("err: %s\n", err)
		return nil, ctx
	}

	ethClient := ethclient.NewClient(client)
	return ethClient, ctx
}
