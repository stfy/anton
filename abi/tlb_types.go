package abi

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/pkg/errors"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type TLBType string

const (
	TLBAddr        TLBType = "addr"
	TLBBool        TLBType = "bool"
	TLBBigInt      TLBType = "bigInt"
	TLBString      TLBType = "string"
	TLBBytes       TLBType = "bytes"
	TLBCell        TLBType = "cell"
	TLBSlice       TLBType = "slice"
	TLBContentCell TLBType = "content"
	TLBStructCell  TLBType = "struct"
	TLBTag         TLBType = "tag"
)

type TelemintText struct {
	Len  uint8  // ## 8
	Text string // bits (len * 8)
}

func (x *TelemintText) LoadFromCell(loader *cell.Slice) error {
	l, err := loader.LoadUInt(8)
	if err != nil {
		return errors.Wrap(err, "load len uint8")
	}

	t, err := loader.LoadSlice(8 * uint(l))
	if err != nil {
		return errors.Wrap(err, "load text slice")
	}

	x.Len = uint8(l)
	x.Text = string(t)

	return nil
}

type StringSnake string

func (x *StringSnake) LoadFromCell(loader *cell.Slice) error {
	s, err := loader.LoadStringSnake()
	if err != nil {
		return err
	}
	*x = StringSnake(s)
	return nil
}

type Opcode string

func (c *Opcode) LoadFromCell(loader *cell.Slice) error {
	l, err := loader.LoadUInt(32)

	if err != nil {
		return errors.Wrap(err, "load len uint32")
	}

	*c = Opcode(fmt.Sprintf("0x%x", l))

	return nil
}

type LimitOrderData struct {
	Expiration       uint32     `tlb:"## 32" json:"expiration"`
	Direction        uint       `tlb:"## 1" json:"direction"`
	Amount           *tlb.Coins `tlb:"." json:"amount"`
	Leverage         uint64     `tlb:"## 64" json:"leverage"`
	LimitPrice       *tlb.Coins `tlb:"." json:"limit_price"`
	StopPrice        *tlb.Coins `tlb:"." json:"stop_price"`
	StopTriggerPrice *tlb.Coins `tlb:"." json:"stop_trigger_price"`
	TakeTriggerPrice *tlb.Coins `tlb:"." json:"take_trigger_price"`
}

type StopOrderData struct {
	Expiration   uint32     `tlb:"## 32" json:"expiration"`
	Direction    uint       `tlb:"## 1" json:"direction"`
	Amount       *tlb.Coins `tlb:"." json:"amount"`
	TriggerPrice *tlb.Coins `tlb:"." json:"limit_price"`
}

type Order struct {
	Value any `tlb:"[StopOrder,TakeOrder,LimitOrder,MarketOrder]"`
}

type StopOrder struct {
	_       tlb.Magic     `tlb:"$0000"`
	Payload StopOrderData `tlb:"."`
}

type TakeOrder struct {
	_       tlb.Magic     `tlb:"$0001"`
	Payload StopOrderData `tlb:"."`
}

type LimitOrder struct {
	_       tlb.Magic      `tlb:"$0010"`
	Payload LimitOrderData `tlb:"."`
}

type MarketOrder struct {
	_       tlb.Magic      `tlb:"$0011"`
	Payload LimitOrderData `tlb:"."`
}

type Orders struct {
	List map[int]Order `json:"list"`
}

func (o *Orders) LoadFromCell(loader *cell.Slice) error {

	d, err := loader.ToDict(4)

	if err != nil {
		return err
	}

	ret := map[int]Order{}

	fmt.Println("LOAD ORDERS", ret)

	for _, item := range d.All() {
		v := Order{}

		key, err := item.Key.BeginParse().LoadUInt(3)

		if err != nil {
			return err
		}

		ref, err := item.Value.BeginParse().LoadRef()

		if err != nil {
			return err
		}

		if err = tlb.LoadFromCell(&v, ref); err != nil {
			return err
		}

		ret[int(key)] = v
	}

	fmt.Println("LOAD ORDERS", ret)

	o.List = ret

	return nil
}

var (
	typeNameRMap = map[reflect.Type]TLBType{
		reflect.TypeOf([]uint8{}): TLBBytes,
	}
	typeNameMap = map[TLBType]reflect.Type{
		TLBBool:        reflect.TypeOf(false),
		"int8":         reflect.TypeOf(int8(0)),
		"int16":        reflect.TypeOf(int16(0)),
		"int32":        reflect.TypeOf(int32(0)),
		"int64":        reflect.TypeOf(int64(0)),
		"uint8":        reflect.TypeOf(uint8(0)),
		"uint16":       reflect.TypeOf(uint16(0)),
		"uint32":       reflect.TypeOf(uint32(0)),
		"uint64":       reflect.TypeOf(uint64(0)),
		TLBBytes:       reflect.TypeOf([]byte{}),
		TLBBigInt:      reflect.TypeOf(big.NewInt(0)),
		TLBCell:        reflect.TypeOf((*cell.Cell)(nil)),
		"dict":         reflect.TypeOf((*cell.Dictionary)(nil)),
		TLBTag:         reflect.TypeOf(tlb.Magic{}),
		"opcode":       reflect.TypeOf((*Opcode)(nil)),
		"coins":        reflect.TypeOf(tlb.Coins{}),
		TLBAddr:        reflect.TypeOf((*address.Address)(nil)),
		TLBString:      reflect.TypeOf((*StringSnake)(nil)),
		"telemintText": reflect.TypeOf((*TelemintText)(nil)),
		"order":        reflect.TypeOf((*Order)(nil)),
		"orders":       reflect.TypeOf((*Orders)(nil)),
	}

	registeredDefinitions = map[TLBType]TLBFieldsDesc{}
)

func init() {
	for n, t := range typeNameMap {
		typeNameRMap[t] = n
	}
}
