package tonapi

import (
	"net/http"
	"ton-node/infra/node/tonapi/tonclient"
	tonModel "ton-node/internal/domain/ton"

	"emperror.dev/errors"
)

type Node struct {
	cli *tonclient.Client
}

var _ tonModel.Node = (*Node)(nil)

func NewNode(nodeURL, apiKey string, opts ...NodeOptFn) (*Node, error) {
	if nodeURL == "" {
		return nil, errors.New("node.NewNode: ton nodeURL not set")
	}

	if apiKey == "" {
		return nil, errors.New("node.NewNode: ton api key not set")
	}

	cfg := &NodeOpts{}

	for idx, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, errors.Wrapf(err, "node.NewNode: ton node option at index %d failed", idx)
		}
	}

	if cfg.Client == nil {
		cfg.Client = &http.Client{
			Transport: http.DefaultTransport,
		}
	}

	return &Node{
		cli: tonclient.NewClient(nodeURL, apiKey, cfg.Client),
	}, nil
}
