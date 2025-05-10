package tonapi

import (
	"context"
	"encoding/json"
	"net/http"

	"emperror.dev/errors"
)

func (n *Node) SendMsg(ctx context.Context, boc string) error {
	body := map[string]string{"boc": boc}

	bytesBody, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	url := n.cli.URL + "/v2/blockchain/message"

	req, code, err := n.cli.PostRequest(ctx, url, bytesBody)
	if err != nil {
		return errors.Wrap(err, "n.cli.PostRequest")
	}

	if code != http.StatusOK {
		return errors.Wrap(
			n.cli.HandleBusinessError(req, code), "SendMsg",
		)
	}

	return nil
}
