package pbaelf

import (
	firecore "github.com/streamingfast/firehose-core"
	"time"
)

var _ firecore.Block = (*Block)(nil)

func (b *Block) GetFirehoseBlockID() string {
	return b.BlockHash
}

func (b *Block) GetFirehoseBlockNumber() uint64 {
	return uint64(b.Header.Height)
}

func (b *Block) GetFirehoseBlockParentID() string {
	return b.Header.PreviousBlockHash
}

func (b *Block) GetFirehoseBlockParentNumber() uint64 {
	return uint64(b.Header.Height) - uint64(1)
}

func (b *Block) GetFirehoseBlockTime() time.Time {
	return b.Header.Time.AsTime().UTC()
}

func (b *Block) GetFirehoseBlockLIBNum() uint64 {
	return 1 // TODO: Fix me
}
