package parser

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/xssnick/tonutils-go/tvm/cell"

	"github.com/tonindexer/anton/internal/app"
	"github.com/tonindexer/anton/internal/core"
)

func (s *Service) parseDirectedMessage(ctx context.Context, acc *core.AccountState, msg *core.Message) error {
	if acc == nil {
		return errors.Wrap(app.ErrImpossibleParsing, "no account data")
	}
	if len(acc.Types) == 0 {
		return errors.Wrap(app.ErrImpossibleParsing, "no interfaces")
	}

	outgoing := acc.Address == msg.SrcAddress
	if outgoing && len(acc.Types) == 1 {
		msg.SrcContract = acc.Types[0]
	}
	if !outgoing && len(acc.Types) == 1 {
		msg.DstContract = acc.Types[0]
	}

	op, err := s.ContractRepo.GetOperationByID(ctx, msg.Type, acc.Types, outgoing, msg.OperationID)
	if errors.Is(err, core.ErrNotFound) {
		return errors.Wrap(app.ErrImpossibleParsing, "unknown operation")
	}
	if err != nil {
		return errors.Wrap(err, "get contract operations")
	}

	msg.OperationName = op.OperationName
	if outgoing {
		msg.SrcContract = op.ContractName
	} else {
		msg.DstContract = op.ContractName
	}

	payloadCell, err := cell.FromBOC(msg.Body)
	if err != nil {
		return errors.Wrap(err, "msg body from boc")
	}

	msgParsed, err := op.Schema.FromCell(payloadCell)
	if err != nil {
		return errors.Wrap(err, "msg body from boc")
	}

	msg.DataJSON, err = json.Marshal(msgParsed)
	if err != nil {
		return errors.Wrap(err, "json marshal parsed payload")
	}

	return nil
}

func (s *Service) ParseMessagePayload(ctx context.Context, msg *core.Message) error {
	var err = app.ErrImpossibleParsing // save message parsing error to a database to look at it later

	// you can parse separately incoming messages to known contracts and outgoing message from them

	if len(msg.Body) == 0 {
		return errors.Wrap(app.ErrImpossibleParsing, "no message body")
	}

	errIn := s.parseDirectedMessage(ctx, msg.DstState, msg)
	if errIn != nil && !errors.Is(errIn, app.ErrImpossibleParsing) {
		err = errors.Wrap(errIn, "incoming")
	}
	if errIn == nil {
		return nil
	}

	errOut := s.parseDirectedMessage(ctx, msg.SrcState, msg)
	if errOut != nil && !errors.Is(errOut, app.ErrImpossibleParsing) {
		err = errors.Wrap(errOut, "outgoing")
	}
	if errOut == nil {
		return nil
	}

	return err
}
