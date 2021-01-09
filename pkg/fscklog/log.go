package fscklog

import (
	"database/sql"
	"strings"
	"time"
)

type ProgressEntry struct {
	Stage    uint8
	Complete uint64
	Total    uint64
	Device   string
}

type BlockRange struct {
	Start uint64
	End   uint64
}

type BlockRanges []BlockRange

func (rs BlockRanges) Count() uint {
	var result uint
	for _, r := range rs {
		result += uint((r.End + 1) - r.Start)
	}
	return result
}

type MultipleClaimedBlockInodeEntry struct {
	Inode                    uint64
	BlockRanges              BlockRanges
	IllegalBlockNumberPassed Uint64Null
}

type FilePathPart struct {
	Type  PathType
	Inode uint64
	Name  string
}

type FilePath []FilePathPart

type FileEntry struct {
	FileInode uint64
	FilePath  FilePath
	ModTime   time.Time
}

type CloneMultiplyClaimedBlocksEntry struct {
	File                        FileEntry
	MultiplyClaimedBlocksAmount uint64
	SharedWithFiles             []FileEntry
	Answer                      sql.NullBool
}

type Log struct {
	Progress                   ProgressEntry
	MultipleClaimedBlockInodes []MultipleClaimedBlockInodeEntry
	CloneMultiplyClaimedBlocks []CloneMultiplyClaimedBlocksEntry
}

type PathType uint

const (
	PathTypeUndefined = PathType(iota)
	PathTypeRoot
	PathTypeInode
	PathTypeName
	PathTypeFilesystemMetadata
	PathTypeUnknown
)

func mustParseFilePath(path string) FilePath {
	parts := strings.Split(path, "/")
	r := make(FilePath, 0, len(parts))
	for _, p := range parts {
		r = append(r, mustParseFilePart(p))
	}
	return r
}

func mustParseFilePart(p string) FilePathPart {
	switch {
	case p == `???`:
		return FilePathPart{
			Type: PathTypeUnknown,
		}
	case p == "":
		return FilePathPart{
			Type: PathTypeRoot,
		}
	case strings.HasPrefix(p, "<") && strings.HasSuffix(p, ">"):
		v := p[1 : len(p)-1]
		switch v {
		case `filesystem metadata`:
			return FilePathPart{
				Type: PathTypeFilesystemMetadata,
			}
		default:
			return FilePathPart{
				Type:  PathTypeInode,
				Inode: mustParseUint64(v),
				Name:  "",
			}
		}
	default:
		return FilePathPart{
			Type: PathTypeName,
			Name: p,
		}
	}
}

func mustParseModTime(modTime string) time.Time {
	v, err := time.Parse(time.ANSIC, modTime)
	if err != nil {
		panic(err)
	}
	return v
}
