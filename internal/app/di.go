package app

import (
	"time"
	"ton-node/infra/node"
	grpcController "ton-node/internal/controller/grpc"
	"ton-node/internal/domain/logger"
	tonModel "ton-node/internal/domain/ton"
	"ton-node/internal/domain/usecase"
	uc "ton-node/internal/usecase"

	"github.com/spf13/viper"
)

type App struct {
	Controller *grpcController.TonNodeController
	UseCase    usecase.UseCase
	NodeSvc    tonModel.NodeService
}

func newApp(l logger.Loggerer) *App {
	nodeSvc := newNodeService(
		viper.GetString("app.node.url"),
		viper.GetString("app.node.api_key"),
		viper.GetDuration("app.node.timeout"),
		viper.GetInt("app.node.rate_limit"),
		viper.GetInt("app.node.rate_burst"),
	)

	useCase := uc.NewUseCase(nodeSvc)
	controller := grpcController.NewTonNodeController(useCase, l)

	return &App{
		NodeSvc:    nodeSvc,
		UseCase:    useCase,
		Controller: controller,
	}
}

func newNodeService(rpcURL, apiKey string, timeout time.Duration, rps, burst int) tonModel.NodeService {
	return node.NewService(rpcURL, apiKey, timeout, rps, burst)
}
