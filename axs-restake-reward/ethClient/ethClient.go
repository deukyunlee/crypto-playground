package ethClient

import (
	"context"
	"github.com/deukyunlee/crypto-playground/logging"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"net/http"
	"time"
)

const (
	EstimateGasMaxRetryCount  = 5
	EstimateGasRetryDelay     = 1 * time.Minute
	PendingNonceMaxRetryCount = 5
	PendingNonceRetryDelay    = 1 * time.Minute
)

var (
	logger = logging.GetLogger()
)

type headerTransport struct {
	Transport http.RoundTripper
	headers   map[string][]string
}

type ClientManger struct {
	Client *ethclient.Client
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

func (ethCli *ClientManger) GetPendingNonceWithRetry(accountAddress common.Address, ctx context.Context) (nonce uint64, err error) {

	for i := 0; i < PendingNonceMaxRetryCount; i++ {
		nonce, err = ethCli.Client.PendingNonceAt(ctx, accountAddress)
		if err != nil {
			logger.Errorf("Failed to get Pending nonce (attempt %d/%d): %v", i+1, PendingNonceMaxRetryCount, err)

			// Delay before retrying
			time.Sleep(PendingNonceRetryDelay)
		} else {
			return nonce, nil
		}
	}

	return 0, nil
}

func (ethCli *ClientManger) EstimateGasWithRetry(ctx context.Context, msg ethereum.CallMsg) (gas uint64, err error) {

	for i := 0; i < EstimateGasMaxRetryCount; i++ {
		gas, err = ethCli.Client.EstimateGas(ctx, msg)
		if err != nil {

			logger.Errorf("EstimateGas failed (attempt %d/%d): %v", i+1, EstimateGasMaxRetryCount, err)

			// Delay before retrying
			time.Sleep(EstimateGasRetryDelay)
		} else {
			return gas, nil // Success
		}
	}

	return 0, err
}
