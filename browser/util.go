package browser

import (
	"errors"
	"os"
)

var (
	internalSyscallErr = errors.New("syscall failed")
)

func getAttr(path string) (attr, error) {
	f, err := os.Lstat(path)
	if err != nil {
		return attr{}, err
	}
	return getAttrWithFileInfo(f)
}

func getAttrWithFileInfo(f os.FileInfo) (attr, error) {
	atime, ctime, mtime, err := statTimes(f)
	if err != nil {
		return attr{}, err
	}
	uid, gid, err := statId(f)
	return attr{
		Name:       f.Name(),
		Dir:        f.IsDir(),
		Mode:       f.Mode(),
		Size:       f.Size(),
		ModifyTime: mtime,
		AccessTime: atime,
		CreateTime: ctime,
		UserID:     uid,
		GroupID:    gid,
	}, nil
}
