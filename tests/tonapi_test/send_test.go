package tonapi_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMsg(t *testing.T) {
	t.Run("Success sent message", func(t *testing.T) {
		code := http.StatusOK
		response := `{}`

		n := getNode(response, code)
		err := n.SendMsg(t.Context(), "")

		assert.NoError(t, err)
	})

	t.Run("Error marshalling the error resp", func(t *testing.T) {
		code := http.StatusBadRequest
		response := `{"bad request"}`

		n := getNode(response, code)
		err := n.SendMsg(t.Context(), "")

		assert.ErrorContains(t, err, "SendMsg: json.Unmarshal: status: 400")
	})

	t.Run("SendMsg error", func(t *testing.T) {
		code := http.StatusBadRequest
		response := `{"error": "bad request"}`

		n := getNode(response, code)
		err := n.SendMsg(t.Context(), "")

		assert.ErrorContains(t, err, "SendMsg: error: bad request")
	})
}
