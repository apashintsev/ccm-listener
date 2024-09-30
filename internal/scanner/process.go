package scanner

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"photon-listener/internal/storage"
	"photon-listener/internal/structures"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"gorm.io/gorm"
)

func handleOpcode(opcode uint64, slice *cell.Slice) {
	switch opcode {
	case TransmitterAdded:
		var tradded structures.TransmitterAdded
		if err1 := tlb.LoadFromCell(&tradded, slice); err1 != nil {
			logrus.Error(err1)
		}
		logrus.Info(&tradded.ProtocolId, "  ", &tradded.Transmitter)
	case TransmitterRemoved:
		fmt.Println("Handling TransmitterRemoved opcode")
		// Логика для TransmitterRemoved

	case Propose:
		fmt.Println("Handling Propose opcode")
		// Логика для Propose

	case ProposalExecuted:
		fmt.Println("Handling ProposalExecuted opcode")
		if opHash, err := slice.LoadUInt(256); err != nil {
			fmt.Println("err")
		} else {
			fmt.Println("opHash=", opHash)
		}
		if success, err := slice.LoadBoolBit(); err != nil {
			fmt.Println("err success")
		} else {
			fmt.Println("success=", success)
		}

	case AllowedProtocolAddressAdded:
		fmt.Println("Handling AllowedProtocolAddressAdded opcode")
		// Логика для AllowedProtocolAddressAdded

	case AllowedProtocolAddressRemoved:
		fmt.Println("Handling AllowedProtocolAddressRemoved opcode")
		// Логика для AllowedProtocolAddressRemoved

	case ExecutorAddressAdded:
		fmt.Println("Handling ExecutorAddressAdded opcode")
		// Логика для ExecutorAddressAdded

	case ExecutorAddressRemoved:
		fmt.Println("Handling ExecutorAddressRemoved opcode")
		// Логика для ExecutorAddressRemoved

	case ConsensusTargetRateSetted:
		fmt.Println("Handling ConsensusTargetRateSetted opcode")
		// Логика для ConsensusTargetRateSetted

	case AllowedProtocolAdded:
		fmt.Println("Handling AllowedProtocolAdded opcode")
		// Логика для AllowedProtocolAdded

	case AllowedProposerAddressAdded:
		fmt.Println("Handling AllowedProposerAddressAdded opcode")
		// Логика для AllowedProposerAddressAdded

	case AllowedProposerAddressRemoved:
		fmt.Println("Handling AllowedProposerAddressRemoved opcode")
		// Логика для AllowedProposerAddressRemoved

	default:
		fmt.Println("Unknown opcode: ", opcode)
	}
}

func (s *scanner) processTransaction(
	trans *tlb.Transaction,
	dbtx *gorm.DB,
	sender *address.Address,
) error {
	if trans.IO.Out == nil {
		return nil
	}

	outMsgs, err := trans.IO.Out.ToSlice()
	if err != nil {
		return nil
	}

	for _, out := range outMsgs {
		if out.MsgType != tlb.MsgTypeExternalOut {
			continue
		}

		if out.Msg.SenderAddr().String() != sender.String() {
			continue
		}

		externalOut := out.AsExternalOut()
		if externalOut.Body == nil {
			continue
		}
		slice := externalOut.Body.BeginParse()
		opcode, err := slice.LoadUInt(32)
		if err != nil {
			continue
		}
		handleOpcode(opcode, externalOut.Body.BeginParse())

		event := storage.Event{
			OpCode:      opcode,
			Data:        "",
			CreatedAt:   time.Unix(int64(externalOut.CreatedAt), 0),
			ProcessedAt: time.Now(),
		}

		if err := dbtx.Create(&event).Error; err != nil {
			return err
		}
		tx := storage.Transaction{
			LT:          trans.LT,
			ProcessedAt: time.Now(),
		}

		if err := dbtx.Create(&tx).Error; err != nil {
			return err
		}
	}

	return nil
}
