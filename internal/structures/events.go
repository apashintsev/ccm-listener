package structures

import (
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"math/big"
)

type TransmitterAdded struct {
	_           tlb.Magic `tlb:"#2a4095d9"`
	ProtocolId  big.Int   `tlb:"## 256"`
	Transmitter big.Int   `tlb:"## 256"`
}
type TransmitterRemoved struct {
	_           tlb.Magic `tlb:"#4258513a"`
	ProtocolId  big.Int   `tlb:"## 256"`
	Transmitter big.Int   `tlb:"## 256"`
}

type ProposalExecuted struct {
	_       tlb.Magic  `tlb:"#7362d09c"`
	OpHash  big.Int    `tlb:"## 256"`
	Success bool       `tlb:"bool"`
	Ret     *cell.Cell `tlb:"either . ^"`
}

type Propose struct {
	_                tlb.Magic  `tlb:"#9ed365f1"`
	ProtocolId       big.Int    `tlb:"## 256"`
	Meta             big.Int    `tlb:"## 256"`
	Nonce            big.Int    `tlb:"## 256"`
	DestChainId      big.Int    `tlb:"## 256"`
	ProtocolAddress  *cell.Cell `tlb:"^"`
	FunctionSelector *cell.Cell `tlb:"^"`
	Params           *cell.Cell `tlb:"^"`
	Reserved         *cell.Cell `tlb:"^"`
}

type AllowedProtocolAdded struct {
	_                   tlb.Magic        `tlb:"#886517b7"`
	ProtocolId          big.Int          `tlb:"## 256"`
	ConsensusTargetRate big.Int          `tlb:"## 256"`
	Transmitters        *cell.Dictionary `tlb:"dict 256"`
}
type AllowedProtocolAddressAdded struct {
	_          tlb.Magic        `tlb:"#6481f58d"`
	ProtocolId big.Int          `tlb:"## 256"`
	Address    *address.Address `tlb:"addr"`
}
type AllowedProtocolAddressRemoved struct {
	_          tlb.Magic        `tlb:"#b0130cb5"`
	ProtocolId big.Int          `tlb:"## 256"`
	Address    *address.Address `tlb:"addr"`
}
type AllowedProposerAddressAdded struct {
	_          tlb.Magic        `tlb:"#cf3d3dfc"`
	ProtocolId big.Int          `tlb:"## 256"`
	Address    *address.Address `tlb:"addr"`
}
type AllowedProposerAddressRemoved struct {
	_          tlb.Magic        `tlb:"#1388a53c"`
	ProtocolId big.Int          `tlb:"## 256"`
	Address    *address.Address `tlb:"addr"`
}
type ExecutorAddressAdded struct {
	_          tlb.Magic        `tlb:"#ff4094eb"`
	ProtocolId big.Int          `tlb:"## 256"`
	Address    *address.Address `tlb:"addr"`
}
type ExecutorAddressRemoved struct {
	_          tlb.Magic        `tlb:"#b9010a9b"`
	ProtocolId big.Int          `tlb:"## 256"`
	Address    *address.Address `tlb:"addr"`
}
type ConsensusTargetRateSetted struct {
	_                   tlb.Magic `tlb:"#8ed04f65"`
	ProtocolId          big.Int   `tlb:"## 256"`
	ConsensusTargetRate big.Int   `tlb:"## 256"`
}

/*

message(0x0000000f) ErrorOccured{
    code: Int;//error code
    arg: String;//addition information about error
}
*/
