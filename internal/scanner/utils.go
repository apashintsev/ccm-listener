package scanner

import (
	"context"
	"fmt"
	"photon-listener/internal/storage"
	"time"

	"github.com/xssnick/tonutils-go/ton"
	"gorm.io/gorm"
)

func (s *scanner) getShardID(shard *ton.BlockIDExt) string {
	return fmt.Sprintf("%d|%d", shard.Workchain, shard.Shard)
}

func (s *scanner) getNonSeenShards(
	ctx context.Context,
	shard *ton.BlockIDExt,
) (ret []*ton.BlockIDExt, err error) {
	if seqno, ok := s.shardLastSeqno[s.getShardID(shard)]; ok && seqno == shard.SeqNo {
		return nil, nil
	}

	block, err := s.api.GetBlockData(ctx, shard)
	if err != nil {
		return nil, fmt.Errorf("get block data err: ", err)
	}

	parents, err := block.BlockInfo.GetParentBlocks()
	if err != nil {
		return nil, fmt.Errorf("get parent blocks (%d:%d): %w", shard.Workchain, shard.Shard, err)
	}

	for _, parent := range parents {
		ext, err := s.getNonSeenShards(ctx, parent)
		if err != nil {
			return nil, err
		}

		ret = append(ret, ext...)
	}

	ret = append(ret, shard)
	return ret, nil
}

// addBlock добавляет новый блок в базу данных на основе информации из мастерчейна.
// Блок сохраняется в таблице через транзакцию базы данных (dbtx), после чего обновляется
// последний сохранённый блок с увеличением его SeqNo.
//
// Аргументы:
//
//	master ton.BlockIDExt - структура с информацией о блоке из мастерчейна.
//	dbtx *gorm.DB - транзакция базы данных, используемая для сохранения нового блока.
//
// Возвращает:
//
//	error - ошибка, если не удалось сохранить блок в базу данных.
//
// Примечание:
//
//	Если сохранение блока прошло успешно, последний блок (`s.lastBlock`) обновляется,
//	и его порядковый номер (SeqNo) увеличивается на 1.
func (s *scanner) addBlock(
	master ton.BlockIDExt,
	dbtx *gorm.DB,
) error {
	newBlock := storage.Block{
		SeqNo:       master.SeqNo,
		WorkChain:   master.Workchain,
		Shard:       master.Shard,
		ProcessedAt: time.Now(),
	}

	if err := dbtx.Create(&newBlock).Error; err != nil {
		return err
	}

	s.lastBlock = newBlock
	s.lastBlock.SeqNo += 1
	return nil
}

// getLastBlockSeqno получает и возвращает последний SeqNo мастерчейна.
// Для получения информации о мастерчейне функция вызывает метод GetMasterchainInfo у API.
//
// Возвращает:
//
//	uint32 - последний seqno мастерчейна.
//	error - ошибка, если не удалось получить информацию с мастерчейна.
//
// Пример ошибки: ошибка сети при обращении к API.
func (s *scanner) getLastBlockSeqno() (uint32, error) {
	lastMaster, err := s.api.GetMasterchainInfo(context.Background())
	if err != nil {
		return 0, err
	}

	return lastMaster.SeqNo, nil
}

// getUniqueShards принимает срез шардов и возвращает срез уникальных шардов.
// Уникальность определяется на основе идентификаторов шардов, полученных с помощью getShardID.
// Если у нескольких шардов одинаковый идентификатор, в результат попадет только один из них.
//
// Аргументы:
//
//	shards []*ton.BlockIDExt - список шардов для обработки.
//
// Возвращает:
//
//	[]*ton.BlockIDExt - список уникальных шардов.
func (s *scanner) getUniqueShards(shards []*ton.BlockIDExt) []*ton.BlockIDExt {
	shardMap := map[string]struct{}{}
	var uniqueShards []*ton.BlockIDExt

	for _, shard := range shards {
		shardID := s.getShardID(shard)
		if _, exists := shardMap[shardID]; !exists {
			shardMap[shardID] = struct{}{}
			uniqueShards = append(uniqueShards, shard)
		}
	}

	return uniqueShards
}
