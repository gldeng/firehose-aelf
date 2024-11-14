package main

import (
	pbaelf "github.com/streamingfast/firehose-aelf/pb/sf/aelf/type/v1"
	firecore "github.com/streamingfast/firehose-core"
	fhCmd "github.com/streamingfast/firehose-core/cmd"
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

		Tools: &firecore.ToolsConfig[*pbaelf.Block]{},
	})
}

// Version value, injected via go build `ldflags` at build time, **must** not be removed or inlined
var version = "dev"
