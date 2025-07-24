package grpc

import (
	"context"

	"github.com/IWannaWish/ethusd-converter/internal/core"
	ethusdpb "github.com/IWannaWish/ethusd-converter/proto/ethusd/gen"
	"github.com/ethereum/go-ethereum/common"
)

type EthusdGRPCServer struct {
	ethusdpb.UnimplementedEthusdConverterServer
	assetService core.AssetService
}

func NewEthusdGRPCServer(assetService core.AssetService) *EthusdGRPCServer {
	return &EthusdGRPCServer{
		assetService: assetService,
	}
}

func (s *EthusdGRPCServer) GetAssets(ctx context.Context, req *ethusdpb.GetAssetsRequest) (*ethusdpb.GetAssetsResponse, error) {
	address := common.HexToAddress(req.GetAddress())

	assets, err := s.assetService.GetAssets(ctx, address)
	if err != nil {
		return nil, err
	}

	var responseAssets []*ethusdpb.Asset
	for _, a := range assets {
		responseAssets = append(responseAssets, &ethusdpb.Asset{
			Symbol:   a.Symbol,
			Balance:  a.Balance,
			UsdValue: a.USDValue,
		})
	}

	return &ethusdpb.GetAssetsResponse{
		Assets: responseAssets,
	}, nil
}
