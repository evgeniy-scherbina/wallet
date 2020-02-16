package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/oklog/run"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pbw "github.com/evgeniy-scherbina/wallet/pb/wallet"
	"github.com/evgeniy-scherbina/wallet/server"
	"github.com/evgeniy-scherbina/wallet/walletdb"
)

var g run.Group

func main() {
	ctxb := context.Background()

	config, err := getConf()
	if err != nil {
		log.Fatal(err)
	}

	db, err := walletdb.New(config.DB.User, config.DB.Pass, config.DB.Addr, config.DB.Name)
	if err != nil {
		log.Fatal(err)
	}

	wallet := server.NewWalletServer(db)

	// PSS GRPC endpoints for connection to hub.
	{
		var grpcServer *grpc.Server
		ctxCancel, cancel := context.WithCancel(ctxb)
		g.Add(func() error {
			defer log.Info("Stop GRPC endpoints")
			for {
				select {
				case <-ctxCancel.Done():
					return nil
				default:
				}

				lis, err := net.Listen("tcp", config.Rpc.Listen)
				if err != nil {
					return fmt.Errorf("failed to listen: %v", err)
				}

				opts := []grpc.ServerOption{
					grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
						grpc_validator.UnaryServerInterceptor(),
					)),
					grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
						grpc_validator.StreamServerInterceptor(),
					)),
				}
				log.Infof("Start GRPC endpoints: %s", config.Rpc.Listen)

				grpcServer = grpc.NewServer(opts...)
				pbw.RegisterWalletServiceServer(grpcServer, wallet)
				_ = grpcServer.Serve(lis)
			}
		}, func(err error) {
			grpcServer.GracefulStop()
			cancel()
		})
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint).
	var httpServer *http.Server
	{
		// Register gRPC server endpoint.
		gwmux := runtime.NewServeMux()
		if err = pbw.RegisterWalletServiceHandlerServer(ctxb, gwmux, wallet); err != nil {
			log.Fatal(err)
		}

		ctxCancel, cancel := context.WithCancel(ctxb)
		g.Add(func() error {
			defer log.Info("Stop Http Server")
			for {
				select {
				case <-ctxCancel.Done():
					return nil
				default:
				}

				log.Infof("Start Http Server on: %s", config.Http.Listen)
				httpServer = &http.Server{
					Addr: config.Http.Listen,
					Handler: gwmux,
				}
				log.Error(httpServer.ListenAndServe())
			}
		}, func(err error) {
			_ = httpServer.Shutdown(context.Background())
			cancel()
		})
	}

	log.Infof("The wallet-service was terminated with: %v", g.Run())
}
