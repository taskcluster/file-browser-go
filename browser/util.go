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

func streamReadResponse(res *ReadResponse) ([]byte, error) {
	if res == nil {
		panic("bad request")
	}

	if !res.IsStreamResponse() {
		return nil, res.Error
	}
	out := make(chan interface{}, 1)
	done := make(chan struct{}, 1)
	data := []byte{}
	var err error

	go func() {
		_ = res.StreamToChannel(out)
		close(done)
	}()

	rb := int(res.requestedBytes)
	for err == nil && len(data) < rb {
		rr := (<-out).(*ReadResponseFrame)
		if rr.Bytes != 0 {
			data = append(data, rr.Buffer...)
		}
		err = rr.Error
	}

	<-done

	return data, err
}
