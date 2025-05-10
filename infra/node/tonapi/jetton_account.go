package tonapi

import (
	"context"
	"fmt"
	"net/http"

	"emperror.dev/errors"
)

func (n *Node) GetJAccount(ctx context.Context, address, jettonContract string) ([]byte, error) {
	url := fmt.Sprintf("%s/v2/accounts/%s/jettons/%s", n.cli.URL, address, jettonContract)

	req, code, err := n.cli.GetRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "n.cli.GetRequest")
	}

	if code != http.StatusOK {
		return nil, errors.Wrap(
			n.cli.HandleBusinessError(req, code), "GetJAccount",
		)
	}

	return req, nil
}
