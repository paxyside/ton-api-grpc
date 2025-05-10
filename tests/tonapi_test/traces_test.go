package tonapi_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmulateTxTrace(t *testing.T) {
	t.Run("Success emulate trace", func(t *testing.T) {
		code := http.StatusOK
		response := `{
			"transaction": {
				"hash": "deadbeef",
				"success": true,
				"total_fees": 1000,
				"account": {
					"address": "0:abc"
				},
				"in_msg": {
					"source": {
						"address": "0:def"
					}
				}
			},
			"children": []
		}`

		n := getNode(response, code)
		_, err := n.EmulateTxTrace(t.Context(), "")

		assert.NoError(t, err)
	})

	t.Run("Error marshalling the error resp", func(t *testing.T) {
		code := http.StatusBadRequest
		response := `{"bad request"}`

		n := getNode(response, code)
		_, err := n.EmulateTxTrace(t.Context(), "")

		assert.ErrorContains(t, err, "EmulateTxTrace: json.Unmarshal: status: 400")
	})

	t.Run("EmulateTxTrace error", func(t *testing.T) {
		code := http.StatusBadRequest
		response := `{"error": "bad request"}`

		n := getNode(response, code)
		_, err := n.EmulateTxTrace(t.Context(), "")

		assert.ErrorContains(t, err, "EmulateTxTrace: error: bad request")
	})
}

func TestGetTxTrace(t *testing.T) {
	t.Run("Success getting trace", func(t *testing.T) {
		code := http.StatusOK
		response := `{
			"transaction": {
				"hash": "deadbeef",
				"success": true,
				"total_fees": 1000,
				"account": {
					"address": "0:abc"
				},
				"in_msg": {
					"source": {
						"address": "0:def"
					}
				}
			},
			"children": []
		}`

		n := getNode(response, code)
		_, err := n.GetTxTrace(t.Context(), "0xdeadbeef")

		assert.NoError(t, err)
	})

	t.Run("Error marshalling the error resp", func(t *testing.T) {
		code := http.StatusBadRequest
		response := `{"bad request"}`

		n := getNode(response, code)
		_, err := n.GetTxTrace(t.Context(), "0xdeadbeef")

		assert.ErrorContains(t, err, "GetTxTrace: json.Unmarshal: status: 400")
	})

	t.Run("EntityNotFound error", func(t *testing.T) {
		code := http.StatusBadRequest
		response := `{"error": "entity not found"}`

		n := getNode(response, code)
		_, err := n.GetTxTrace(t.Context(), "0xdeadbeef")

		assert.ErrorContains(t, err, "entity not found")
	})

	t.Run("GetTxTrace error", func(t *testing.T) {
		code := http.StatusBadRequest
		response := `{"error": "bad request"}`

		n := getNode(response, code)
		_, err := n.GetTxTrace(t.Context(), "0xdeadbeef")

		assert.ErrorContains(t, err, "GetTxTrace: error: bad request")
	})
}
