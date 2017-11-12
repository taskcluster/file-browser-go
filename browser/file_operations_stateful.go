package browser

import (
	"context"
	"os"
)

func init() {
	localRegistry.registerCommand(SF_OPEN, sf_open)
	localRegistry.registerCommand(SF_CLOSE, sf_close)
}

// Stateful Open request
type Sf_openRequest struct {
	RequestID requestID   `msgpack:requestID`
	Path      string      `msgpack:path`
	Flags     int         `msgpack:flags`
	Perm      os.FileMode `msgpack:perm`
}

func (op *Sf_openRequest) GetRequestID() requestID {
	return op.RequestID
}

func (op *Sf_openRequest) GenerateErrorResponse(err error) opResponse {
	return &Sf_openResponse{
		RequestID: op.RequestID,
		HandleID:  0,
		Error:     err,
	}
}

type Sf_openResponse struct {
	opResponseBase
	RequestID requestID `msgpack:requestID`
	HandleID  uint64    `msgpack:handleID`
	Error     error     `msgpack:error`
}

func sf_open(ctx context.Context, req opRequest, callback func()) opResponse {
	if callback != nil {
		defer callback()
	}
	op, ok := req.(*Sf_openRequest)
	if !ok {
		return &Sf_openResponse{
			HandleID:  0,
			RequestID: req.GetRequestID(),
			Error:     errBadRequest,
		}
	}

	select {
	case <-ctx.Done():
		return op.GenerateErrorResponse(errInterrupted)
	default:
	}

	handleID, err := localFileRegistry.openFile(op.Path, op.Flags, op.Perm)
	if err != nil {
		return op.GenerateErrorResponse(err)
	}
	return &Sf_openResponse{
		RequestID: op.RequestID,
		HandleID:  handleID,
	}
}

// Close
type Sf_closeRequest struct {
	RequestID requestID `msgpack:requestID`
	HandleID  uint64    `msgpack:handleID`
}

func (cr *Sf_closeRequest) GetRequestID() requestID {
	return cr.RequestID
}

func (cr *Sf_closeRequest) GenerateErrorResponse(err error) opResponse {
	return &Sf_closeResponse{
		RequestID: cr.RequestID,
		Error:     err,
	}
}

type Sf_closeResponse struct {
	opResponseBase
	RequestID requestID `msgpack:requestID`
	Error     error     `msgpack:error`
}

func sf_close(ctx context.Context, req opRequest, callback func()) opResponse {
	if callback != nil {
		defer callback()
	}
	cr, ok := req.(*Sf_closeRequest)
	if !ok {
		return &Sf_closeResponse{
			RequestID: cr.GetRequestID(),
			Error:     errBadRequest,
		}
	}

	select {
	case <-ctx.Done():
		return cr.GenerateErrorResponse(errInterrupted)
	default:
	}

	return cr.GenerateErrorResponse(localFileRegistry.closeFile(cr.HandleID))
}
