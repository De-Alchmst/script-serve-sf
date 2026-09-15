# Script-Serve-FS

Serve regular files regularly and executable files as their output

## build

```
go mod tity
go build
```

## usage
```
script-serve-fs <source-dir> <mountpoint>
```

This fs, by design, allows other users to access it. Thus, you need to enable
`user_allow_other` in `/etc/fuse.conf`.

Mounts `source-dir` to `mountpoint`, so that non-executable files are readable
and executable files are readable and writable.
On read, the file is executed with reader's `PID` as it's only argument and
`STDOUT` is returned to the reader.
On write, the file is executed with reader's `PID` as it's first argument and
the written data as it's second. The `STDOUT` will be given on the next read.

Then serve over `sftp` or `9p` or whatever else you fancy.

To connect with [`sshfs`](https://github.com/libfuse/sshfs),
use the `-o direct_io` flag, else you can't read the executable files.
