package main

import (
	"errors"
	"fmt"
	pbbstream "github.com/streamingfast/bstream/pb/sf/bstream/v1"
	"github.com/streamingfast/firehose-aelf/block"
	"github.com/streamingfast/firehose-aelf/pb/aelf"
	pbaelf "github.com/streamingfast/firehose-aelf/pb/sf/aelf/type/v1"
	firecore "github.com/streamingfast/firehose-core"
	fhCmd "github.com/streamingfast/firehose-core/cmd"
	"github.com/streamingfast/firehose-core/firehose/info"
	"github.com/streamingfast/firehose-core/node-manager/mindreader"
	"github.com/streamingfast/logging"
	pbfirehose "github.com/streamingfast/pbgo/sf/firehose/v2"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"strings"
)

type ReaderWithConverter struct {
	inner       mindreader.ConsolerReader
	fromTypeUrl string
	toTypeUrl   string
}

func (r ReaderWithConverter) ReadBlock() (blk *pbbstream.Block, err error) {

	blk, err = r.inner.ReadBlock()
	if err != nil {
		return blk, err
	}
	if clean(blk.Payload.TypeUrl) != r.fromTypeUrl {
		return nil, errors.New("unrecognized type to convert from")
	}
	var parsed aelf.Block
	err = proto.Unmarshal(blk.Payload.Value, &parsed)
	if err != nil {
		return nil, errors.New("unable to unmarshal aelf.Block")
	}
	converted := block.ConvertBlock(blk.Id, &parsed)
	newPayloadBytes, err := proto.Marshal(converted)
	if err != nil {
		return nil, errors.New("unable to marshal pbaelf.Block")
	}
	blk.Payload = &anypb.Any{
		TypeUrl: r.toTypeUrl,
		Value:   newPayloadBytes,
	}

	return blk, nil
}

func (r ReaderWithConverter) Done() <-chan interface{} {
	return r.inner.Done()
}

func clean(in string) string {
	return strings.Replace(in, "type.googleapis.com/", "", 1)
}

func newReaderWithConverter(lines chan string, blockEncoder firecore.BlockEncoder, logger *zap.Logger, tracer logging.Tracer) (mindreader.ConsolerReader, error) {

	inner, err := firecore.NewConsoleReader(lines, blockEncoder, logger, tracer)
	if err != nil {
		return inner, err
	}
	fromTypeUrl := new(aelf.Block).ProtoReflect().Descriptor().FullName()
	toTypeUrl := new(pbaelf.Block).ProtoReflect().Descriptor().FullName()
	return &ReaderWithConverter{
		inner:       inner,
		fromTypeUrl: clean(string(fromTypeUrl)),
		toTypeUrl:   string(toTypeUrl),
	}, nil
}

func main() {
	fhCmd.Main(&firecore.Chain[*pbaelf.Block]{
		ShortName:            "aelf",
		LongName:             "AElf",
		ExecutableName:       "AElf.Launcher",
		FullyQualifiedModule: "github.com/streamingfast/firehose-aelf",
		Version:              version,

		FirstStreamableBlock: 1,

		BlockFactory:         func() firecore.Block { return new(pbaelf.Block) },
		ConsoleReaderFactory: newReaderWithConverter,
		InfoResponseFiller: func(firstStreamableBlock *pbbstream.Block, resp *pbfirehose.InfoResponse, validate bool) error {
			aelfBlock := &pbaelf.Block{}
			if err := firstStreamableBlock.Payload.UnmarshalTo(aelfBlock); err != nil && validate {
				return fmt.Errorf("cannot decode first streamable block: %w", err)
			}
			if err := info.DefaultInfoResponseFiller(firstStreamableBlock, resp, validate); err != nil && validate {
				return err
			}
			return nil
		},

		Tools: &firecore.ToolsConfig[*pbaelf.Block]{},
	})
}

// Version value, injected via go build `ldflags` at build time, **must** not be removed or inlined
var version = "dev"
