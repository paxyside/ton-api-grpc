package usecase

import (
	"context"
	tonModel "ton-node/internal/domain/ton"

	"emperror.dev/errors"
)

type UseCase struct {
	nodeSvc tonModel.NodeService
}

func NewUseCase(nodeSvc tonModel.NodeService) *UseCase { return &UseCase{nodeSvc: nodeSvc} }

func (u *UseCase) GetAccountInfo(ctx context.Context, address string) ([]byte, error) {
	cli, err := u.nodeSvc.GetNode()
	if err != nil {
		return nil, errors.Wrap(err, "u.nodeSvc.GetNode")
	}

	return cli.GetAccount(ctx, address)
}

func (u *UseCase) GetJAccountInfo(ctx context.Context, address, jettonContract string) ([]byte, error) {
	cli, err := u.nodeSvc.GetNode()
	if err != nil {
		return nil, errors.Wrap(err, "u.nodeSvc.GetNode")
	}

	return cli.GetJAccount(ctx, address, jettonContract)
}

func (u *UseCase) GetSeqNumber(ctx context.Context, address string) (uint64, error) {
	cli, err := u.nodeSvc.GetNode()
	if err != nil {
		return 0, errors.Wrap(err, "u.nodeSvc.GetNode")
	}

	return cli.GetSeqno(ctx, address)
}

func (u *UseCase) GetTransactionTrace(ctx context.Context, messageHash string) ([]byte, error) {
	cli, err := u.nodeSvc.GetNode()
	if err != nil {
		return nil, errors.Wrap(err, "u.nodeSvc.GetNode")
	}

	return cli.GetTxTrace(ctx, messageHash)
}

func (u *UseCase) EmulateTransactionTrace(ctx context.Context, boc string) ([]byte, error) {
	cli, err := u.nodeSvc.GetNode()
	if err != nil {
		return nil, errors.Wrap(err, "u.nodeSvc.GetNode")
	}

	return cli.EmulateTxTrace(ctx, boc)
}

func (u *UseCase) SendMessage(ctx context.Context, boc string) (string, error) {
	cli, err := u.nodeSvc.GetNode()
	if err != nil {
		return "", errors.Wrap(err, "u.nodeSvc.GetNode")
	}

	if err = cli.SendMsg(ctx, boc); err != nil {
		return "", errors.Wrap(err, "cli.SendMsg")
	}

	return "success", nil
}
