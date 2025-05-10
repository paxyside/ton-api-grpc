package tonapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	tonModel "ton-node/internal/domain/ton"

	"emperror.dev/errors"
)

func (n *Node) GetTxTrace(ctx context.Context, messageHash string) ([]byte, error) {
	url := fmt.Sprintf("%s/v2/traces/%s", n.cli.URL, messageHash)

	req, code, err := n.cli.GetRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "n.cli.GetRequest")
	}

	if code != http.StatusOK {
		errMap := map[string]error{
			tonModel.EntityNotFound: tonModel.ErrEntityNotFound,
		}

		return nil, errors.Wrap(
			n.cli.HandleBusinessErrorWithMap(req, code, errMap), "GetTxTrace",
		)
	}

	return req, nil
}

func (n *Node) EmulateTxTrace(ctx context.Context, boc string) ([]byte, error) {
	body := map[string]string{"boc": boc}

	bytesBody, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}

	url := n.cli.URL + "/v2/traces/emulate"

	req, code, err := n.cli.PostRequest(ctx, url, bytesBody)
	if err != nil {
		return nil, errors.Wrap(err, "n.cli.PostRequest")
	}

	if code != http.StatusOK {
		return nil, errors.Wrap(
			n.cli.HandleBusinessError(req, code), "EmulateTxTrace",
		)
	}

	return req, nil
}
