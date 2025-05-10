package model

import (
	"context"
)

type NodeService interface {
	GetNode() (Node, error)
}

type Node interface {
	GetAccount(ctx context.Context, address string) ([]byte, error)
	GetJAccount(ctx context.Context, address, jettonContract string) ([]byte, error)
	GetSeqno(ctx context.Context, address string) (uint64, error)
	GetTxTrace(ctx context.Context, messageHash string) ([]byte, error)
	EmulateTxTrace(ctx context.Context, boc string) ([]byte, error)
	SendMsg(ctx context.Context, boc string) error
}
