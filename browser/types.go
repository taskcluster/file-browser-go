package browser

import (
	"errors"
	"os"
)

type requestID uint64

type commandCode uint8

const (
	SL_READ commandCode = iota
	SL_WRITE
	SL_TRUNC
	SL_CREATE
	SL_REMOVE
	SL_RENAME
	SL_MKDIR
	SL_READDIR
	SL_STAT

	SF_READ
	SF_WRITE
	SF_OPEN
	SF_CLOSE
	SF_READDIR
	SF_STAT

	OP_INTR
)

// If set, browser will allow stateful operations
var SF_SET bool = false

var (
	errStreamingNotSupported = errors.New("streaming is not supported by this operation")
	errBadRequest            = errors.New("bad request")
	errInternalError         = errors.New("internal error")
	errStateless             = errors.New("browser is running is stateless mode")
	errInvalidHandle         = errors.New("handle not valid")
	errInterrupted           = errors.New("operation interrupted")
)

/* attr describes the attributes of a file/dir.
 * Important attributes:
 * Name: Name of the file/directory (string)
 * Dir: true if directory (bool)
 * Mode: File Mode (os.FileMode)
 * Size: File size (uint64)
 * UserID: File UID (uint32)
 * GroupID: File GID (uint32)
 * AccessTime: Last access time since epoch (int64)
 * ModifyTime: Last modified time since epoch (int64)
 * CreateTime: Last create time since epoch (int64)
 */
type attr struct {
	Name       string      `msgpack:name`
	Dir        bool        `msgpack:dir`
	Mode       os.FileMode `msgpack:mode`
	Size       int64       `msgpack:size`
	UserID     uint32      `msgpack:uid`
	GroupID    uint32      `msgpack:gid`
	AccessTime int64       `msgpack:atim`
	ModifyTime int64       `msgpack:mtim`
	CreateTime int64       `msgpack:ctim`
}

/* dirent decribes a directory entry.
 * Since users will not be able to use the inode
 * number directly, this will not be returned.
 * Short of using a system call, the only information
 * returned through os.FileInfo is:
 * Name: name of the file (string)
 * Dir: true if directory (bool)
 * Size: length of file in bytes (system dependent) (int64)
 * Mode: file mode (os.FileMode which is a uint32)
 */
type dirent struct {
	Name string      `msgpack:name`
	Size int64       `msgpack:size`
	Dir  bool        `msgpack:dir`
	Mode os.FileMode `msgpack:mode`
}

type opRequest interface {
	GetRequestID() requestID
	GenerateErrorResponse(error) opResponse
}

type opResponse interface {
	IsStreamResponse() bool
	StreamToChannel(chan<- interface{}) error
}

type opResponseBase struct {
}

func (o *opResponseBase) IsStreamResponse() bool {
	return false
}

func (o *opResponseBase) StreamToChannel(out chan<- interface{}) error {
	return errStreamingNotSupported
}
