package scanner

import (
	"context"
	"log"
	"photon-listener/internal/app"
	"photon-listener/internal/storage"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"gopkg.in/tomb.v1"
)

type scanner struct {
	api            *ton.APIClient
	lastBlock      storage.Block
	shardLastSeqno map[string]uint32
}

func NewScanner() (*scanner, error) {
	client := liteclient.NewConnectionPool()

	if err := client.AddConnectionsFromConfig(
		context.Background(),
		app.CFG.NetworkConfig,
	); err != nil {
		return nil, err
	}

	api := ton.NewAPIClient(client)

	return &scanner{
		api:            api,
		lastBlock:      storage.Block{},
		shardLastSeqno: make(map[string]uint32),
	}, nil
}

func (s *scanner) Listen() {
	logrus.Info("[SCN] start scanning blocks")

	if err := app.DB.Last(&s.lastBlock).Error; err != nil {
		s.lastBlock.SeqNo = 0
	} else {
		s.lastBlock.SeqNo += 1
	}

	if s.lastBlock.SeqNo == 0 {
		lastMaster, err := s.api.GetMasterchainInfo(context.Background())
		for err != nil {
			time.Sleep(time.Second)
			logrus.Error("[SCN] error when get last master: ", err)
			lastMaster, err = s.api.GetMasterchainInfo(context.Background())
		}

		s.lastBlock.SeqNo = lastMaster.SeqNo
		s.lastBlock.Shard = lastMaster.Shard
		s.lastBlock.WorkChain = lastMaster.Workchain
	}

	s.processBlocks()
}

func (s *scanner) processBlocks() error {
	timeStart := time.Now()
	master, err := s.api.CurrentMasterchainInfo(context.Background()) // we fetch block just to trigger chain proof check
	if err != nil {
		logrus.Info("get masterchain info err: ", err.Error())
		return err
	}
	account, err := s.api.GetAccount(context.Background(), master, app.CFG.TargetContractAddress)
	if err != nil {
		logrus.Error("get account err:", err.Error())
		return err
	}
	txList, err := s.api.ListTransactions(context.Background(), app.CFG.TargetContractAddress, 50, account.LastTxLT, account.LastTxHash)
	if err != nil {
		log.Printf("send err: %s", err.Error())
		return err
	}
	// Получение списка последних транзакций из базы данных
	var lastTxs []storage.Transaction
	if err := app.DB.Order("LT").Limit(50).Find(&lastTxs).Error; err != nil {
		logrus.Error("error fetching last transactions: ", err)
		return err
	}

	// Фильтруем txList, убирая элементы с совпадающим LT
	filteredTxList := []*tlb.Transaction{}

	for _, tx := range txList {
		found := false
		for _, lastTx := range lastTxs {
			if tx.LT == lastTx.LT {
				found = true
				break
			}
		}
		if !found {
			filteredTxList = append(filteredTxList, tx)
		}
	}

	var wg sync.WaitGroup
	var tombTrans tomb.Tomb

	// Создаем канал, чтобы отслеживать завершение всех горутин
	allDone := make(chan struct{})

	for _, transaction := range filteredTxList {
		wg.Add(1)

		go func(transaction *tlb.Transaction) {
			defer wg.Done()

			// Создаем новую транзакцию базы данных для каждой горутины
			dbtx := app.DB.Begin()
			defer func() {
				if r := recover(); r != nil {
					dbtx.Rollback()
				}
			}()

			if err := s.processTransaction(transaction, dbtx, app.CFG.TargetContractAddress); err != nil {
				tombTrans.Kill(err)
				dbtx.Rollback()
				return
			}

			if err := dbtx.Commit().Error; err != nil {
				tombTrans.Kill(err)
			}
		}(transaction)
	}

	// Горутин, которая закрывает канал, когда все транзакции завершены
	go func() {
		wg.Wait()
		close(allDone)
	}()

	select {
	case <-allDone:
		// Все транзакции успешно завершены
	case <-tombTrans.Dying():
		// Если произошла ошибка, вернуть ошибку
		logrus.Error("[SCN] err when processing transactions: ", tombTrans.Err())
		return tombTrans.Err()
	}

	// Запись блока после успешной обработки транзакций
	if err := s.addBlock(*master, app.DB); err != nil {
		return err
	}

	lastSeqno, err := s.getLastBlockSeqno()
	if err != nil {
		logrus.Infof("[SCN] success process block [%d] time to process block [%0.2fs] trans count [%d]",
			master.SeqNo,
			time.Since(timeStart).Seconds(),
			len(txList),
		)
	} else {
		logrus.Infof("[SCN] success process block [%d|%d] time to process block [%0.2fs] trans count [%d]",
			master.SeqNo,
			lastSeqno,
			time.Since(timeStart).Seconds(),
			len(txList),
		)
	}
	time.Sleep(time.Second * 2)
	return nil
}
