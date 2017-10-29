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
	case SL_READ:
		req = &Sl_readRequest{}
	case SL_READDIR:
		req = &Sl_readdirRequest{}
	case SF_READDIR:
		req = &Sf_readdirRequest{}
	case SF_OPEN:
		req = &Sf_openRequest{}
	case SF_CLOSE:
		req = &Sf_closeRequest{}
	case SL_MKDIR:
		req = &Sl_mkdirRequest{}
	}
	err = msgpack.Unmarshal(m.buf, &req)
	if err != nil {
		return nil
	}
	return req
}
