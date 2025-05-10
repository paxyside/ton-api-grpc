package model

import (
	"encoding/base64"
	"encoding/binary"

	"emperror.dev/errors"
	"github.com/sigurn/crc16"
)

const DefaultDataLength = 36

func ValidateAddress(value interface{}) error {
	addressStr, ok := value.(string)
	if !ok {
		return errors.New("address is not a string")
	}

	data, err := base64.RawURLEncoding.DecodeString(addressStr)
	if err != nil {
		return errors.Wrap(err, "base64.RawURLEncoding.DecodeString")
	}

	if len(data) != DefaultDataLength {
		return errors.New("invalid address length")
	}

	checksum := data[len(data)-2:]
	if crc16.Checksum(data[:len(data)-2], crc16.MakeTable(crc16.CRC16_XMODEM)) != binary.BigEndian.Uint16(checksum) {
		return errors.New("invalid address checksum")
	}

	return nil
}
