package tonapi_test

import (
	"net/http"
	"ton-node/infra/node/tonapi"
	"ton-node/infra/node/tonapi/testutils"
	tonModel "ton-node/internal/domain/ton"
)

func getNode(resp string, code int, opts ...func(response *http.Response)) tonModel.Node {
	cli, _ := tonapi.NewNode("https://example.com", "apikey", tonapi.WithHTTPClient(&http.Client{
		Transport: testutils.StubRT{
			Resp:       resp,
			StatusCode: code,
			Opts:       opts,
		}}))

	return cli
}
