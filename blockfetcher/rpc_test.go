package blockfetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	pbbstream "github.com/streamingfast/bstream/pb/sf/bstream/v1"
	"github.com/streamingfast/eth-go"
	"github.com/streamingfast/eth-go/rpc"
	pbeth "github.com/streamingfast/firehose-ethereum/types/pb/sf/ethereum/type/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBlockFetcherFetchPBEthReturnsErrorOnNilRPCBlock(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(`{"jsonrpc":"2.0","id":"0x1","result":null}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	fetcher := NewBlockFetcher(0, 0, 1, func(*rpc.Block, map[string]*rpc.TransactionReceipt, map[string][]eth.Log, *zap.Logger) (*pbeth.Block, map[string]bool) {
		t.Fatal("toEthBlock should not be called when the RPC block is nil")
		return nil, nil
	}, zap.NewNop())
	fetcher.latest = 1

	block, err := fetcher.FetchPBEth(context.Background(), rpc.NewClient(server.URL), 1)
	require.Nil(t, block)
	require.EqualError(t, err, "fetching block 1: rpc returned nil block")
}

func TestBlockFetcherFetchReturnsErrorInsteadOfPanickingOnNilConvertedBlock(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(`{"jsonrpc":"2.0","id":"0x1","result":{"number":"0x1","hash":"0x1111111111111111111111111111111111111111111111111111111111111111","transactions":[]}}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	fetcher := NewBlockFetcher(0, 0, 1, func(*rpc.Block, map[string]*rpc.TransactionReceipt, map[string][]eth.Log, *zap.Logger) (*pbeth.Block, map[string]bool) {
		return nil, nil
	}, zap.NewNop())
	fetcher.latest = 1

	var (
		block *pbbstream.Block
		err   error
	)
	require.NotPanics(t, func() {
		block, err = fetcher.Fetch(context.Background(), rpc.NewClient(server.URL), 1)
	})
	require.Nil(t, block)
	require.EqualError(t, err, `converting block 1 "0x1111111111111111111111111111111111111111111111111111111111111111": converter returned nil protobuf block`)
}
