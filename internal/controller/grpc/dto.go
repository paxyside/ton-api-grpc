package grpc

import (
	model "ton-node/internal/domain/ton"

	"emperror.dev/errors"
)

func ValidateAddress(addresses ...string) error {
	for _, address := range addresses {
		if err := model.ValidateAddress(address); err != nil {
			return errors.Wrap(err, "model.ValidateAddress")
		}
	}
	return nil
}
