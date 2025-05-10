package tonapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"emperror.dev/errors"
)

func (n *Node) GetAccount(ctx context.Context, address string) ([]byte, error) {
	url := fmt.Sprintf("%s/v2/accounts/%s", n.cli.URL, address)

	req, code, err := n.cli.GetRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "n.cli.GetRequest")
	}

	if code != http.StatusOK {
		return nil, errors.Wrap(
			n.cli.HandleBusinessError(req, code), "GetAccount",
		)
	}

	return req, nil
}

func (n *Node) GetSeqno(ctx context.Context, address string) (uint64, error) {
	url := fmt.Sprintf("%s/v2/wallet/%s/seqno", n.cli.URL, address)

	req, code, err := n.cli.GetRequest(ctx, url)
	if err != nil {
		return 0, errors.Wrap(err, "n.cli.GetRequest")
	}

	if code != http.StatusOK {
		return 0, errors.Wrap(
			n.cli.HandleBusinessError(req, code), "GetSeqno",
		)
	}

	res := struct {
		Seqno uint64 `json:"seqno"`
	}{}

	if err = json.Unmarshal(req, &res); err != nil {
		return 0, errors.Wrap(err, "json.Unmarshal")
	}

	return res.Seqno, nil
}
