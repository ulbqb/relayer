package lbm

import (
	"context"
	"strings"
	"time"

	ics23 "github.com/confio/ics23/go"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	commitmenttypes "github.com/cosmos/ibc-go/v5/modules/core/23-commitment/types"
	abci "github.com/line/ostracon/abci/types"
	provtypes "github.com/line/ostracon/light/provider"
	prov "github.com/line/ostracon/light/provider/http"
	crypto "github.com/line/ostracon/proto/ostracon/crypto"
	rpcclient "github.com/line/ostracon/rpc/client"
	rpchttp "github.com/line/ostracon/rpc/client/http"
	libclient "github.com/line/ostracon/rpc/jsonrpc/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChainClient struct {
	LightProvider provtypes.Provider
	RPCClient     rpcclient.Client
	Config        *LBMProviderConfig
}

func NewChainClient(ccc *LBMProviderConfig) (*ChainClient, error) {
	cc := &ChainClient{
		Config: ccc,
	}
	if err := cc.Init(); err != nil {
		return nil, err
	}
	return cc, nil
}

func (cc *ChainClient) Init() error {
	timeout, _ := time.ParseDuration(cc.Config.Timeout)
	rpcClient, err := NewRPCClient(cc.Config.RPCAddr, timeout)
	if err != nil {
		return err
	}

	lightprovider, err := prov.New(cc.Config.ChainID, cc.Config.RPCAddr)
	if err != nil {
		return err
	}

	cc.RPCClient = rpcClient
	cc.LightProvider = lightprovider

	return nil
}

func NewRPCClient(addr string, timeout time.Duration) (*rpchttp.HTTP, error) {
	httpClient, err := libclient.DefaultHTTPClient(addr)
	if err != nil {
		return nil, err
	}
	httpClient.Timeout = timeout
	rpcClient, err := rpchttp.NewWithClient(addr, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}
	return rpcClient, nil
}

func (cc *ChainClient) QueryABCI(ctx context.Context, req abci.RequestQuery) (abci.ResponseQuery, error) {
	opts := rpcclient.ABCIQueryOptions{
		Height: req.Height,
		Prove:  req.Prove,
	}
	result, err := cc.RPCClient.ABCIQueryWithOptions(ctx, req.Path, req.Data, opts)
	if err != nil {
		return abci.ResponseQuery{}, err
	}

	if !result.Response.IsOK() {
		return abci.ResponseQuery{}, sdkErrorToGRPCError(result.Response)
	}

	// data from trusted node or subspace query doesn't need verification
	if !opts.Prove || !isQueryStoreWithProof(req.Path) {
		return result.Response, nil
	}

	return result.Response, nil
}

func sdkErrorToGRPCError(resp abci.ResponseQuery) error {
	switch resp.Code {
	case sdkerrors.ErrInvalidRequest.ABCICode():
		return status.Error(codes.InvalidArgument, resp.Log)
	case sdkerrors.ErrUnauthorized.ABCICode():
		return status.Error(codes.Unauthenticated, resp.Log)
	case sdkerrors.ErrKeyNotFound.ABCICode():
		return status.Error(codes.NotFound, resp.Log)
	default:
		return status.Error(codes.Unknown, resp.Log)
	}
}

// isQueryStoreWithProof expects a format like /<queryType>/<storeName>/<subpath>
// queryType must be "store" and subpath must be "key" to require a proof.
func isQueryStoreWithProof(path string) bool {
	if !strings.HasPrefix(path, "/") {
		return false
	}

	paths := strings.SplitN(path[1:], "/", 3)

	switch {
	case len(paths) != 3:
		return false
	case paths[0] != "store":
		return false
	case rootmulti.RequireProof("/" + paths[2]):
		return true
	}

	return false
}

func SDKDotStringifyEvent(e abci.Event) sdk.StringEvent {
	res := sdk.StringEvent{Type: e.Type}

	for _, attr := range e.Attributes {
		res.Attributes = append(
			res.Attributes,
			sdk.Attribute{Key: string(attr.Key), Value: string(attr.Value)},
		)
	}

	return res
}

func CommitmenttypesDotConvertProofs(tmProof *crypto.ProofOps) (commitmenttypes.MerkleProof, error) {
	if tmProof == nil {
		return commitmenttypes.MerkleProof{}, sdkerrors.Wrapf(commitmenttypes.ErrInvalidMerkleProof, "ostracon proof is nil")
	}
	// Unmarshal all proof ops to CommitmentProof
	proofs := make([]*ics23.CommitmentProof, len(tmProof.Ops))
	for i, op := range tmProof.Ops {
		var p ics23.CommitmentProof
		err := p.Unmarshal(op.Data)
		if err != nil || p.Proof == nil {
			return commitmenttypes.MerkleProof{}, sdkerrors.Wrapf(commitmenttypes.ErrInvalidMerkleProof, "could not unmarshal proof op into CommitmentProof at index %d: %v", i, err)
		}
		proofs[i] = &p
	}
	return commitmenttypes.MerkleProof{
		Proofs: proofs,
	}, nil
}
