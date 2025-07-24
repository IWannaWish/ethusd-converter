package main

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"time"

	"github.com/IWannaWish/ethusd-converter/cmd/cli/display"
	"github.com/IWannaWish/ethusd-converter/internal/applog"
	"github.com/IWannaWish/ethusd-converter/internal/config"
	"github.com/IWannaWish/ethusd-converter/internal/core/mapper"
	ethusdpb "github.com/IWannaWish/ethusd-converter/proto/ethusd/gen"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load("config.env")
	cfg := config.Load()
	logger := applog.NewLogger(cfg)
	ctx := context.Background()

	if len(os.Args) < 2 {
		logger.Error(ctx, "Usage: ./ethusd-converter <ethereum_address>")
		os.Exit(1)
	}

	rawAddr := os.Args[1]
	if !common.IsHexAddress(rawAddr) {
		logger.Error(ctx, "Неверный Ethereum адрес", applog.String("address", rawAddr))
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		}),
	)

	if err != nil {
		logger.Error(ctx, "Ошибка подключения к gRPC-серверу", applog.WithStack(err)...)
		os.Exit(1)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Error(ctx, "Ошибка при закрытии gRPC-соединения", applog.WithStack(err)...)
		}
	}()

	client := ethusdpb.NewEthusdConverterClient(conn)

	resp, err := client.GetAssets(ctx, &ethusdpb.GetAssetsRequest{Address: rawAddr})
	if err != nil {
		logger.Error(ctx, "Ошибка при вызове GetAssets", applog.WithStack(err)...)
		os.Exit(1)
	}

	assets, total, err := mapper.FromProtoAssets(resp.Assets)
	if err != nil {
		logger.Error(ctx, "Ошибка преобразования активов", applog.WithStack(err)...)
		os.Exit(1)
	}

	display.NewTablePrinter().Print(assets, total)
}
