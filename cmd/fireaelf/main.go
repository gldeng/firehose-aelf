package main

import (
	"fmt"
	pbbstream "github.com/streamingfast/bstream/pb/sf/bstream/v1"
	pbaelf "github.com/streamingfast/firehose-aelf/pb/aelf"
	firecore "github.com/streamingfast/firehose-core"
	fhCmd "github.com/streamingfast/firehose-core/cmd"
	"github.com/streamingfast/firehose-core/firehose/info"
	pbfirehose "github.com/streamingfast/pbgo/sf/firehose/v2"
)

func main() {
	fhCmd.Main(&firecore.Chain[*pbaelf.Block]{
		ShortName:            "aelf",
		LongName:             "AElf",
		ExecutableName:       "AElf.Launcher",
		FullyQualifiedModule: "github.com/streamingfast/firehose-aelf",
		Version:              version,

		FirstStreamableBlock: 1,

		BlockFactory:         func() firecore.Block { return new(pbaelf.Block) },
		ConsoleReaderFactory: firecore.NewConsoleReader,
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
