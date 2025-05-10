package grpc

import (
	"context"
	"ton-node/internal/controller/grpc/tonnodepb"
	"ton-node/internal/domain/logger"
	"ton-node/internal/domain/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TonNodeController struct {
	tonnodepb.UnimplementedTonNodeServiceServer
	uc usecase.UseCase
	l  logger.Loggerer
}

func NewTonNodeController(uc usecase.UseCase, l logger.Loggerer) *TonNodeController {
	return &TonNodeController{
		uc: uc,
		l:  l,
	}
}

func (c *TonNodeController) GetAccount(
	ctx context.Context, req *tonnodepb.GetAccountRequest,
) (*tonnodepb.AccountInfoResponse, error) {
	err := ValidateAddress(req.GetAddress())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ValidateAddress: %v", err)
	}

	accInfo, err := c.uc.GetAccountInfo(ctx, req.GetAddress())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "c.uc.GetAccountInfo: %v", err)
	}

	return &tonnodepb.AccountInfoResponse{
		RawJsonAccountInfo: accInfo,
	}, nil
}

func (c *TonNodeController) GetJAccount(
	ctx context.Context, req *tonnodepb.GetJAccountRequest,
) (*tonnodepb.JettonAccountInfoResponse, error) {
	err := ValidateAddress(req.GetAddress(), req.GetJettonContract())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ValidateAddress: %v", err)
	}

	accInfo, err := c.uc.GetJAccountInfo(ctx, req.GetAddress(), req.GetJettonContract())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "c.uc.GetJAccountInfo: %v", err)
	}

	return &tonnodepb.JettonAccountInfoResponse{
		RawJsonJettonAccountInfo: accInfo,
	}, nil
}

func (c *TonNodeController) GetSeqno(
	ctx context.Context, req *tonnodepb.GetSeqnoRequest,
) (*tonnodepb.GetSeqnoResponse, error) {
	err := ValidateAddress(req.GetAddress())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ValidateAddress: %v", err)
	}

	seqno, err := c.uc.GetSeqNumber(ctx, req.GetAddress())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "c.uc.GetSeqno: %v", err)
	}

	return &tonnodepb.GetSeqnoResponse{
		Seqno: seqno,
	}, nil
}

func (c *TonNodeController) SendMsg(
	ctx context.Context, req *tonnodepb.SendMsgRequest,
) (*tonnodepb.SendMsgResponse, error) {
	success, err := c.uc.SendMessage(ctx, req.GetBoc())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "c.uc.SendMessage: %v", err)
	}

	return &tonnodepb.SendMsgResponse{
		Status: success,
	}, nil
}

func (c *TonNodeController) EmulateTxTrace(
	ctx context.Context, req *tonnodepb.EmulateTxTraceRequest,
) (*tonnodepb.TraceResponse, error) {
	trace, err := c.uc.EmulateTransactionTrace(ctx, req.GetBoc())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "c.uc.EmulateTransactionTrace: %v", err)
	}

	return &tonnodepb.TraceResponse{
		RawJsonTrace: trace,
	}, nil
}

func (c *TonNodeController) GetTxTrace(
	ctx context.Context, req *tonnodepb.GetTxTraceRequest,
) (*tonnodepb.TraceResponse, error) {
	trace, err := c.uc.GetTransactionTrace(ctx, req.GetMessageHash())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "c.uc.GetTransactionTrace: %v", err)
	}

	return &tonnodepb.TraceResponse{
		RawJsonTrace: trace,
	}, nil
}
