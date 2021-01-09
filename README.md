```
xaionaro@void:~/go/src/github.com/xaionaro-go/fsck-multiple-claim-stats$ go run ./cmd/multiple-claim-progress/ /tmp/fsck.log
last progress report: fscklog.ProgressEntry{Stage:0x1, Complete:0x2461e, Total:0x2461e, Device:"/dev/md127"}
inodes: 7894/133402 (0.1%)
blocks: 20847/3527526  (0.0%)
clone-size: 85 MB/14 GB (assuming block-size is 4096)

xaionaro@void:~/go/src/github.com/xaionaro-go/fsck-multiple-claim-stats$ go run ./cmd/top-multiple-claim-inodes/ /tmp/fsck.log
inode #190042244: 960705 blocks
```
