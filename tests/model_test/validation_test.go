package model_test

import (
	"testing"
	model "ton-node/internal/domain/ton"

	"github.com/stretchr/testify/require"
)

func TestAddressValidator_Validate(t *testing.T) {
	t.Run("Valid address", func(t *testing.T) {
		address := "EQBkEFTw8riiCs6WlfE2gaB9hajBCrpJuB-UTEANggY0ym5F"
		err := model.ValidateAddress(address)
		require.NoError(t, err)
	})
	t.Run("Invalid address (not a string)", func(t *testing.T) {
		address := 123456
		err := model.ValidateAddress(address)
		require.Contains(t, err.Error(), "address is not a string")
	})
	t.Run("Invalid address (error decoding address)", func(t *testing.T) {
		address := "EQBkEFTw8riiCs6WlfE2gaB9hajBCrpJuB-UTEANggY0ym5F!"
		err := model.ValidateAddress(address)
		require.Contains(t, err.Error(), "base64.RawURLEncoding.DecodeString")
	})
	t.Run("Invalid address (invalid address length)", func(t *testing.T) {
		address := "EQBkEFTw8riiCs6WlfE2gaB9hajBCrpJuB"
		err := model.ValidateAddress(address)
		require.Contains(t, err.Error(), "invalid address length")
	})
	t.Run("Invalid address (invalid address checksum)", func(t *testing.T) {
		address := "EQBkEFTw8riiCs6WlfE2gaB9hajBCrpJuB-UTEANggY0ym5X"
		err := model.ValidateAddress(address)
		require.Contains(t, err.Error(), "invalid address checksum")
	})
}
