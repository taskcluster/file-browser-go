package browser

// encoders and decoders for all command structs
// defined here.

import (
	"github.com/vmihailenco/msgpack"
)

// message struct to wrap fields in
type Message struct {
	code commandCode
	buf  []byte
}

func DecodeRequestMessage(m *Message) opRequest {
	var err error
	var req opRequest

	switch m.code {
	case SL_READDIR:
		req = &Sl_readdirRequest{}
	case SF_READDIR:
		req = &Sf_readdirRequest{}
	case SL_MKDIR:
		req = &Sl_mkdirRequest{}
	case SL_READ:
		req = &Sl_readRequest{}
	case SL_WRITE:
		req = &Sl_writeRequest{}
	case SL_CREATE:
		req = &Sl_createRequest{}
	case SL_TRUNC:
		req = &Sl_truncRequest{}
	case SL_STAT:
		req = &Sl_statRequest{}
	case SL_REMOVE:
		req = &Sl_removeRequest{}
	case SL_RENAME:
		req = &Sl_renameRequest{}
	case SF_OPEN:
		req = &Sf_openRequest{}
	case SF_CLOSE:
		req = &Sf_closeRequest{}
	}
	err = msgpack.Unmarshal(m.buf, &req)
	if err != nil {
		return nil
	}
	return req
}
