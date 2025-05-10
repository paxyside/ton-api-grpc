package tonapi_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetJBalance(t *testing.T) {
	t.Run("Success getting jetton balance", func(t *testing.T) {
		code := http.StatusOK
		response := `{
		  "balance": "123456789000",
		  "wallet_address": {
			"address": "0:abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			"is_scam": false,
			"is_wallet": true
		  },
		  "jetton": {
			"address": "0:deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
			"name": "MegaTON",
			"symbol": "MGT",
			"decimals": 9,
			"image": "https://ton.org/jettons/megatron/logo.png",
			"verification": "verified",
			"score": 95
		  }
		}`

		n := getNode(response, code)
		_, err := n.GetJAccount(
			t.Context(),
			"0QBkEFTw8riiCs6WlfE2gaB9hajBCrpJuB-UTEANggY0yogK",
			"0QBkEFTw8riiCs6WlfE2gaB9hajBCrpJuB-UTEANggY0yogK",
		)

		assert.NoError(t, err)
	})

	t.Run("Error marshalling the error resp", func(t *testing.T) {
		code := http.StatusBadRequest
		response := `{"bad request"}`

		n := getNode(response, code)
		_, err := n.GetJAccount(
			t.Context(),
			"0QBkEFTw8riiCs6WlfE2gaB9hajBCrpJuB-UTEANggY0yogK",
			"0QBkEFTw8riiCs6WlfE2gaB9hajBCrpJuB-UTEANggY0yogK",
		)

		assert.ErrorContains(t, err, "GetJAccount: json.Unmarshal: status: 400")
	})

	t.Run("GetJAccount error", func(t *testing.T) {
		code := http.StatusBadRequest
		response := `{"error": "bad request"}`

		n := getNode(response, code)
		_, err := n.GetJAccount(
			t.Context(),
			"0QBkEFTw8riiCs6WlfE2gaB9hajBCrpJuB-UTEANggY0yogK",
			"0QBkEFTw8riiCs6WlfE2gaB9hajBCrpJuB-UTEANggY0yogK",
		)

		assert.ErrorContains(t, err, "GetJAccount: error: bad request")
	})
}
