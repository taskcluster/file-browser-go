// +build darwin

package browser

import (
	"os"
	"syscall"
)

func statTimes(fi os.FileInfo) (atime int64, mtime int64, ctime int64, err error) {
	if fi == nil || fi.Sys() == nil {
		return 0, 0, 0, internalSyscallErr
	}
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, 0, 0, internalSyscallErr
	}
	err = nil
	atime, _ = stat.Atimespec.Unix()
	ctime, _ = stat.Ctimespec.Unix()
	mtime, _ = stat.Mtimespec.Unix()
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
