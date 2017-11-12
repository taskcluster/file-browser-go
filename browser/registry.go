package browser

import (
	"context"
	"sync"
)

type registry struct {
	cm map[commandCode]func(context.Context, opRequest, func()) opResponse
	sync.RWMutex
	op map[requestID]context.CancelFunc
}

func (r *registry) registerCommand(code commandCode, f func(context.Context, opRequest, func()) opResponse) {
	r.Lock()
	defer r.Unlock()
	r.cm[code] = f
}

func (r *registry) getCommand(code commandCode) func(context.Context, opRequest, func()) opResponse {
	r.RLock()
	defer r.RUnlock()
	return r.cm[code]
}

func (r *registry) register(reqID requestID, cancel context.CancelFunc) {
	r.Lock()
	defer r.Unlock()
	r.op[reqID] = cancel
}

func (r *registry) unregister(reqID requestID) {
	r.Lock()
	defer r.Unlock()
	delete(r.op, reqID)
}

func (r *registry) callCancelFunc(reqID requestID) {
	r.RLock()
	defer r.RUnlock()
	cancel := r.op[reqID]
	if cancel != nil {
		cancel()
	}
}
