package tonapi_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccount(t *testing.T) {
	t.Run("Success getting account", func(t *testing.T) {
		getNode(`{"ok": true}`, 200, func(resp *http.Response) {
			resp.Header.Set("X-Test", "true")
		})

		code := http.StatusOK
		response := `{
			"address": "0:abc",
			"balance": 4766049003,
			"last_activity": 123456,
			"status": "active",
			"interfaces": [],
			"get_methods": [],
			"is_wallet": true
		}`

		n := getNode(response, code)
		_, err := n.GetAccount(
			t.Context(),
			"0QBkEFTw8riiCs6WlfE2gaB9hajBCrpJuB-UTEANggY0yogK",
		)

		assert.NoError(t, err)
	})

	t.Run("Error marshalling the error resp", func(t *testing.T) {
		code := http.StatusBadRequest
		response := `{"bad request"}`

		n := getNode(response, code)
		_, err := n.GetAccount(
			t.Context(),
			"0QBkEFTw8riiCs6WlfE2gaB9hajBCrpJuB-UTEANggY0yogK",
		)

		assert.ErrorContains(t, err, "GetAccount: json.Unmarshal: status: 400")
	})

	t.Run("GetAccount error", func(t *testing.T) {
		code := http.StatusBadRequest
		response := `{"error": "bad request"}`

		n := getNode(response, code)
		_, err := n.GetAccount(
			t.Context(),
			"0QBkEFTw8riiCs6WlfE2gaB9hajBCrpJuB-UTEANggY0yogK",
		)

		assert.ErrorContains(t, err, "GetAccount: error: bad request")
	})
}

func TestGetSeqno(t *testing.T) {
	t.Run("Success getting seqno", func(t *testing.T) {
		code := http.StatusOK
		response := `{"seqno": 50}`

		n := getNode(response, code)
		_, err := n.GetSeqno(
			t.Context(),
			"0QBkEFTw8riiCs6WlfE2gaB9hajBCrpJuB-UTEANggY0yogK",
		)

		assert.NoError(t, err)
	})

	t.Run("Error marshalling the error resp", func(t *testing.T) {
		code := http.StatusBadRequest
		response := `{"bad request"}`

		n := getNode(response, code)
		_, err := n.GetSeqno(
			t.Context(),
			"0QBkEFTw8riiCs6WlfE2gaB9hajBCrpJuB-UTEANggY0yogK",
		)

		assert.ErrorContains(t, err, "GetSeqno: json.Unmarshal: status: 400")
	})

	t.Run("GetSeqno error", func(t *testing.T) {
		code := http.StatusBadRequest
		response := `{"error": "bad request"}`

		n := getNode(response, code)
		_, err := n.GetSeqno(
			t.Context(),
			"0QBkEFTw8riiCs6WlfE2gaB9hajBCrpJuB-UTEANggY0yogK",
		)

		assert.ErrorContains(t, err, "GetSeqno: error: bad request")
	})

	t.Run("Unmarshal error", func(t *testing.T) {
		code := http.StatusOK
		response := `{"seqno": "50"}`

		n := getNode(response, code)
		_, err := n.GetSeqno(
			t.Context(),
			"0QBkEFTw8riiCs6WlfE2gaB9hajBCrpJuB-UTEANggY0yogK",
		)

		assert.ErrorContains(t, err, "json.Unmarshal")
	})
}
