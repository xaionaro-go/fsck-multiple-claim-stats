package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/xaionaro-go/fsck-multiple-claim-stats/pkg/fscklog"
)

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "usage: %s /path/to/fsck/output/file\n", os.Args[0])
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	logsPath := flag.Arg(0)
	logsBytes, err := ioutil.ReadFile(logsPath)
	assertNoError(err)
	logs, err := fscklog.Parse(logsBytes)
	assertNoError(err)

	fmt.Printf("last progress report: %#+v\n", logs.Progress)

	fmt.Printf("inodes: %d/%d (%.1f%%)\n",
		len(logs.CloneMultiplyClaimedBlocks),
		len(logs.MultipleClaimedBlockInodes),
		float64(len(logs.CloneMultiplyClaimedBlocks))/float64(len(logs.MultipleClaimedBlockInodes)),
	)

	clonedBlocks, toBeClonedBlocks := uint(0), uint(0)
	for _, entry := range logs.MultipleClaimedBlockInodes {
		toBeClonedBlocks += uint(len(entry.BlockRanges))
	}
	for _, entry := range logs.CloneMultiplyClaimedBlocks {
		clonedBlocks += uint(entry.MultiplyClaimedBlocksAmount)
	}
	fmt.Printf("blocks: %d/%d  (%.1f%%)\n",
		clonedBlocks,
		toBeClonedBlocks,
		float64(clonedBlocks)/float64(toBeClonedBlocks),
	)
	blockSize := uint(4096)
	fmt.Printf("clone-size: %s/%s (assuming block-size is %d)\n",
		humanize.Bytes(uint64(clonedBlocks*blockSize)),
		humanize.Bytes(uint64(toBeClonedBlocks*blockSize)),
		blockSize,
	)
}
