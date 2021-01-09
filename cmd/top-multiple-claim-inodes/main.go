package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

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

	if len(logs.MultipleClaimedBlockInodes) == 0 {
		return
	}

	sort.Slice(logs.MultipleClaimedBlockInodes, func(i, j int) bool {
		return logs.MultipleClaimedBlockInodes[i].BlockRanges.Count() > logs.MultipleClaimedBlockInodes[j].BlockRanges.Count()
	})

	top := logs.MultipleClaimedBlockInodes[0]
	fmt.Printf("inode #%d: %d blocks\n", top.Inode, top.BlockRanges.Count())
}
