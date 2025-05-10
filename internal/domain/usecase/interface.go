package usecase

import "context"

type UseCase interface {
	GetAccountInfo(ctx context.Context, address string) ([]byte, error)
	GetJAccountInfo(ctx context.Context, address, jettonContract string) ([]byte, error)
	GetSeqNumber(ctx context.Context, address string) (uint64, error)
	GetTransactionTrace(ctx context.Context, messageHash string) ([]byte, error)
	EmulateTransactionTrace(ctx context.Context, boc string) ([]byte, error)
	SendMessage(ctx context.Context, boc string) (string, error)
}
