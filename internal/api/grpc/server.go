package grpc

import (
	"context"

	"github.com/IWannaWish/ethusd-converter/internal/core"
	ethusd_pb "github.com/IWannaWish/ethusd-converter/proto/ethusd/gen"
	"github.com/ethereum/go-ethereum/common"
)

type EthusdGRPCServer struct {
	ethusd_pb.UnimplementedEthusdConverterServer
	assetService core.AssetService
}

func NewEthusdGRPCServer(assetService core.AssetService) *EthusdGRPCServer {
	return &EthusdGRPCServer{
		assetService: assetService,
	}
}

func (s *EthusdGRPCServer) GetAssets(ctx context.Context, req *ethusd_pb.GetAssetsRequest) (*ethusd_pb.GetAssetsResponse, error) {
	address := common.HexToAddress(req.GetAddress())

	assets, err := s.assetService.GetAssets(ctx, address)
	if err != nil {
		return nil, err
	}

	var responseAssets []*ethusd_pb.Asset
	for _, a := range assets {
		responseAssets = append(responseAssets, &ethusd_pb.Asset{
			Symbol:   a.Symbol,
			Balance:  a.Balance,
			UsdValue: a.USDValue,
		})
	}

	return &ethusd_pb.GetAssetsResponse{
		Assets: responseAssets,
	}, nil
}
