package fscklog

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/xaionaro-go/unsafetools"
)

func Parse(_b []byte) (*Log, error) {
	fsckLog := &Log{}
	s := unsafetools.CastBytesToString(_b)
	clonedMultiplyHeaderParseRegexStr := `^File (.*) \(inode #(\d+), mod time (.*)\)\s*\n\s+has (\d+) multiply-claimed block\(s\), shared with (\d+) file\(s\):`
	clonedMultiplyHeaderParseRegex := regexp.MustCompile(clonedMultiplyHeaderParseRegexStr)
	clonedMultiplyEntryParseRegexStr := `^\s+(<filesystem metadata>|([^\n]+) \(inode #(\d+), mod time (.*)\))\s*\n`
	clonedMultiplyEntryParseRegex := regexp.MustCompile(clonedMultiplyEntryParseRegexStr)
	clonedMultiplyFooterParseRegexStr := `^Clone multiply-claimed blocks\? (yes|no)`
	clonedMultiplyFooterParseRegex := regexp.MustCompile(clonedMultiplyFooterParseRegexStr)
	illegalBlockNumberRegex := regexp.MustCompile(`Illegal block number passed to ext2fs_test_block_bitmap #(\d+) for multiply claimed block map$`)
	for len(s) > 0 {
		endIdx := strings.Index(s, "\n")
		if endIdx == 0 {
			s = s[1:]
			continue
		}
		if endIdx < 0 {
			break
		}
		curS := s
		line := s[:endIdx]
		s = s[endIdx+1:]
		words := strings.SplitN(line, " ", 6)
		switch {
		case len(words) == 4 && isUint8(words[0]) && isUint64(words[1]) && isUint64(words[2]):
			fsckLog.Progress.Stage = mustParseUint8(words[0])
			fsckLog.Progress.Complete = mustParseUint64(words[1])
			fsckLog.Progress.Total = mustParseUint64(words[2])
			fsckLog.Progress.Device = words[3]
		case strings.HasPrefix(line, "Multiply-claimed block(s) in inode "):
			entry := MultipleClaimedBlockInodeEntry{}
			inodeStr := words[4][:len(words[4])-1]
			inode, err := strconv.ParseUint(inodeStr, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("unable to parse inode number '%s': %w", inodeStr, err)
			}
			entry.Inode = inode
			blockRangesStr := words[5]

			matches := illegalBlockNumberRegex.FindAllStringSubmatch(line, 1)
			if matches != nil {
				block, err := strconv.ParseUint(matches[0][1], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("unable to parse block number '%s': %w", matches[0][1], err)
				}
				entry.IllegalBlockNumberPassed = Uint64Null{
					Uint64: block,
					Valid:  true,
				}
				blockRangesStr = blockRangesStr[:len(blockRangesStr)-len(matches[0][0])]
			}

			blockRangeStrs := strings.Split(blockRangesStr, " ")
			entry.BlockRanges = make([]BlockRange, 0, len(blockRangeStrs))
			for _, blockRangeStr := range blockRangeStrs {
				parts := strings.Split(blockRangeStr, `--`)
				start, err := strconv.ParseUint(parts[0], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("unable to parse block number '%s': %w", parts[0], err)
				}
				var end uint64
				if len(parts) > 1 {
					end, err = strconv.ParseUint(parts[1], 10, 64)
					if err != nil {
						return nil, fmt.Errorf("unable to parse block number '%s': %w", parts[1], err)
					}
				} else {
					end = start
				}
				entry.BlockRanges = append(entry.BlockRanges, BlockRange{
					Start: start,
					End:   end,
				})
			}
			fsckLog.MultipleClaimedBlockInodes = append(fsckLog.MultipleClaimedBlockInodes, entry)
		case strings.HasPrefix(line, "File ") && strings.HasPrefix(s, "  has "):
			entry := CloneMultiplyClaimedBlocksEntry{}
			{
				matches := clonedMultiplyHeaderParseRegex.FindAllStringSubmatch(curS, -1)
				match := matches[0]
				entry.File.FilePath = mustParseFilePath(match[1])
				entry.File.FileInode = mustParseUint64(match[2])
				entry.File.ModTime = mustParseModTime(match[3])
				entry.MultiplyClaimedBlocksAmount = mustParseUint64(match[4])
				entry.SharedWithFiles = make([]FileEntry, mustParseUint64(match[5]))
				s = curS[len(match[0]):]
			}
			for idx := range entry.SharedWithFiles {
				f := &entry.SharedWithFiles[idx]
				matches := clonedMultiplyEntryParseRegex.FindAllStringSubmatch(s, -1)
				match := matches[0]
				if match[4] != "" {
					f.FilePath = mustParseFilePath(match[2])
					f.FileInode = mustParseUint64(match[3])
					f.ModTime = mustParseModTime(match[4])
				} else {
					f.FilePath = mustParseFilePath(match[1])
				}
				_ = f
				s = s[len(match[0]):]
			}
			{
				matches := clonedMultiplyFooterParseRegex.FindAllStringSubmatch(s, -1)
				if matches != nil {
					match := matches[0]
					s = s[len(match[0]):]
					switch match[1] {
					case "yes":
						entry.Answer = sql.NullBool{
							Bool:  true,
							Valid: true,
						}
					case "no":
						entry.Answer = sql.NullBool{
							Bool:  false,
							Valid: true,
						}
					default:
						panic(fmt.Sprintf("%v", match))
					}
				}
			}
			fsckLog.CloneMultiplyClaimedBlocks = append(fsckLog.CloneMultiplyClaimedBlocks, entry)
		}
	}
	return fsckLog, nil
}
