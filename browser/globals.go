package browser

import (
	"context"
)

var localRegistry *registry = &registry{
	cm: make(map[commandCode]func(context.Context, opRequest, func()) opResponse),
	op: make(map[requestID]context.CancelFunc),
}

var localFileRegistry *fileRegistry = &fileRegistry{
	openFiles:    make(map[uint64]*fileHandle),
	nextHandleID: 0x00000001,
}
