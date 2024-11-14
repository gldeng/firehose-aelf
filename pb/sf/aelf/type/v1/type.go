package pbaelf

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	firecore "github.com/streamingfast/firehose-core"
	"google.golang.org/protobuf/proto"
)

var _ firecore.Block = (*Block)(nil)

func (h *Hash) ToHex() string {
	return hex.EncodeToString(h.Value)
}

func (h *BlockHeader) GetHash() (*Hash, error) {
	var (
		bytes []byte
		err   error
	)
	bytes, err = proto.Marshal(h)
	if err != nil {
		return nil, err
	}
	if h.Signature != nil {
		header := &BlockHeader{}
		if err := proto.Unmarshal(bytes, header); err != nil {
			return nil, err
		}
		header.Signature = nil
		bytes, err = proto.Marshal(header)
		if err != nil {
			return nil, err
		}
	}
	hash := sha256.Sum256(bytes)
	return &Hash{Value: hash[:]}, nil
}

func (b *Block) GetFirehoseBlockID() string {
	hash, err := b.Header.GetHash()
	if err != nil {
		return ""
	}
	return hash.ToHex()
}

func (b *Block) GetFirehoseBlockNumber() uint64 {
	return uint64(b.Header.Height)
}

func (b *Block) GetFirehoseBlockParentID() string {
	if b.Header.PreviousBlockHash == nil {
		return ""
	}

	return b.Header.PreviousBlockHash.ToHex()
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
