package node

import (
	"sync"
	"time"
	ton_api "ton-node/infra/node/tonapi"
	tonModel "ton-node/internal/domain/ton"

	"emperror.dev/errors"
)

type Service struct {
	rpcURL string
	apiKey string

	timeout time.Duration
	rps     int
	burst   int

	mx   *sync.Mutex
	node tonModel.Node
}

func NewService(rpcURL, apiKey string, timeout time.Duration, rps, burst int) *Service {
	return &Service{
		rpcURL:  rpcURL,
		apiKey:  apiKey,
		timeout: timeout,
		rps:     rps,
		burst:   burst,
		mx:      &sync.Mutex{},
	}
}

var _ tonModel.NodeService = (*Service)(nil)

func (n *Service) GetNode() (tonModel.Node, error) {
	n.mx.Lock()
	defer n.mx.Unlock()

	if n.node != nil {
		return n.node, nil
	}

	cli, err := n.buildNode(n.rpcURL, n.apiKey)
	if err != nil {
		return nil, errors.Wrap(err, "n.buildNode")
	}

	n.node = cli
	return cli, nil
}

func (n *Service) buildNode(url, apiKey string) (tonModel.Node, error) {
	opts := []ton_api.NodeOptFn{
		ton_api.WithTimeout(n.timeout),
		ton_api.WithRateLimit(n.rps, n.burst),
	}

	cli, err := ton_api.NewNode(url, apiKey, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "ton.NewNode")
	}

	return cli, nil
}
