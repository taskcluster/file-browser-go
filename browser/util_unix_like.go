// +build linux,unix

package browser

import (
	"os"
	"syscall"
)

// use only for linux
func statTimes(fi os.FileInfo) (atime int64, mtime int64, ctime int64, err error) {
	if fi == nil || fi.Sys() == nil {
		return 0, 0, 0, internalSyscallErr
	}
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, 0, 0, internalSyscallErr
	}
	err = nil
	atime, _ = stat.Atim.Unix()
	ctime, _ = stat.Ctim.Unix()
	mtime, _ = stat.Mtim.Unix()
	return
}

// use only for linux
func statId(fi os.FileInfo) (uid uint32, gid uint32, err error) {
	if fi == nil || fi.Sys == nil {
		return 0, 0, internalSyscallErr
	}
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, 0, internalSyscallErr
	}
	err = nil
	uid = stat.Uid
	gid = stat.Gid
	return
}
