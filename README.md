```
$ go get github.com/xaionaro-go/fsck-multiple-claim-stats/...

$ "$(go env GOPATH)"/bin/multiple-claim-progress /tmp/fsck.log
last progress report: fscklog.ProgressEntry{Stage:0x1, Complete:0x2461e, Total:0x2461e, Device:"/dev/md127"}
inodes: 7894/133402 (5.9%)
blocks: 20847/3527526  (0.6%)
clone-size: 85 MB/14 GB (assuming block-size is 4096)

$ "$(go env GOPATH)"/bin/top-multiple-claim-inodes /tmp/fsck.log
inode #190042244: 960705 blocks
```
