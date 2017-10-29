package browser

import (
	"context"
	"os"
)

const defaultFrameSize = 1024

func init() {
	localRegistry.registerCommand(SL_READ, sl_read)
	localRegistry.registerCommand(SL_WRITE, sl_write)
	localRegistry.registerCommand(SL_CREATE, sl_create)
	localRegistry.registerCommand(SL_TRUNC, sl_trunc)
	localRegistry.registerCommand(SF_OPEN, sf_open)
	localRegistry.registerCommand(SF_CLOSE, sf_close)
}

// Read
type Sl_readRequest struct {
	RequestID requestID `msgpack:requestID`
	// File path
	Path string `msgpack:path`
	// Offset is always relative to file file origin
	// To read from end, set offset = size-1-offset
	Offset      int64 `msgpack:offset`
	BytesToRead int64 `msgpack:bytesToRead`
}

func (r *Sl_readRequest) GetRequestID() requestID {
	return r.RequestID
}

func (r *Sl_readRequest) GenerateErrorResponse(err error) opResponse {
	return &ReadResponse{
		RequestID: r.RequestID,
		Error:     err,
		FrameSize: defaultFrameSize,
	}
}

// Common to both stateless and stateful versions
type ReadResponse struct {
	RequestID requestID `msgpack:requestID`
	Error     error     `msgpack:error`
	FrameSize int64     `msgpack:frameSize`
	// Number of bytes in the read request
	requestedBytes int64
	// Handle to the file so that data can be streamed
	// It is the responsibility of the read function to
	// seek to the required position
	file *os.File
	// Wait for a stream response only if this value is
	// true
	Stream bool `msgpack:streaming`
	// context used to check for cancellation
	ctx context.Context
}

type ReadResponseFrame struct {
	Buffer []byte `msgpack:buffer`
	// Number of bytes in the buffer
	Bytes int `msgpack:bytes`
	// Used to set EOF
	Error error `msgpack:error`
}

func (r *ReadResponse) IsStreamResponse() bool {
	return true
}

func (r *ReadResponse) StreamToChannel(out chan<- interface{}) error {
	bytesSent := int64(0)
	// Close the file struct after streaming
	defer func() {
		_ = r.file.Close()
	}()
	for bytesSent != r.requestedBytes {
		bufsize := r.FrameSize
		if r.requestedBytes-bytesSent < r.FrameSize {
			bufsize = r.requestedBytes - bytesSent
		}

		select {
		case <-r.ctx.Done():
			out <- &ReadResponseFrame{
				Error: errInterrupted,
			}
			return errInterrupted
		default:
		}

		data := make([]byte, bufsize)
		n, err := r.file.Read(data)
		data = data[:n]
		frame := &ReadResponseFrame{
			Buffer: data,
			Bytes:  n,
			Error:  err,
		}

		out <- frame
		bytesSent += int64(n)
		if err != nil {
			return err
		}
	}
	return nil
}

func sl_read(ctx context.Context, req opRequest, callback func()) opResponse {
	if callback != nil {
		defer callback()
	}
	rr, ok := req.(*Sl_readRequest)
	if !ok {
		return &ReadResponse{
			RequestID: req.GetRequestID(),
			Error:     errBadRequest,
		}
	}
	file, err := os.Open(rr.Path)
	if err != nil {
		return rr.GenerateErrorResponse(err)
	}
	_, err = file.Seek(rr.Offset, 0) // Seek from file origin
	if err != nil {
		return rr.GenerateErrorResponse(err)
	}

	// Streaming is the most expensive operation so check for
	// cancellation before streaming
	select {
	case <-ctx.Done():
		_ = file.Close()
		return rr.GenerateErrorResponse(errInterrupted)
	default:
	}
	// Calling StreamToChannel should safely stream the data out
	return &ReadResponse{
		RequestID:      rr.RequestID,
		FrameSize:      defaultFrameSize,
		requestedBytes: rr.BytesToRead,
		Stream:         true,
		// The offset for the file has been set
		file: file,
	}
}

// Write requests
// 1. Stateless write requests
type Sl_writeRequest struct {
	RequestID requestID `msgpack:requestID`
	Path      string    `msgpack:path`
	Offset    int64     `msgpack:offset`
	Buf       []byte    `msgpack:buf`
}

func (wr *Sl_writeRequest) GetRequestID() requestID {
	return wr.RequestID
}

func (wr *Sl_writeRequest) GenerateErrorResponse(err error) opResponse {
	return &WriteResponse{
		RequestID: wr.RequestID,
		Error:     err,
	}
}

type WriteResponse struct {
	opResponseBase
	RequestID    requestID `msgpack:requestID`
	BytesWritten int       `msgpack:bytesWritten`
	Error        error     `msgpack:error`
}

// sl_write requires that the file exists before writing.
// If the file does not exist prior to writing, use sl_create
// to create the file and then call sl_write.
func sl_write(ctx context.Context, req opRequest, callback func()) opResponse {
	if callback != nil {
		defer callback()
	}
	wr, ok := req.(*Sl_writeRequest)
	if !ok {
		return &WriteResponse{
			RequestID: req.GetRequestID(),
			Error:     errBadRequest,
		}
	}
	file, err := os.OpenFile(wr.Path, os.O_WRONLY, os.FileMode(0666))
	if err != nil {
		return wr.GenerateErrorResponse(err)
	}
	defer func() {
		_ = file.Close()
	}()

	select {
	case <-ctx.Done():
		return wr.GenerateErrorResponse(errInterrupted)
	default:
	}

	n, err := file.WriteAt(wr.Buf, wr.Offset)
	return &WriteResponse{
		RequestID:    wr.RequestID,
		BytesWritten: n,
		Error:        err,
	}
}

// Stateless create
type Sl_createRequest struct {
	RequestID requestID   `msgpack:requestID`
	Path      string      `msgpack:path`
	Mode      os.FileMode `msgpack:mode`
}

func (cr *Sl_createRequest) GetRequestID() requestID {
	return cr.RequestID
}

func (cr *Sl_createRequest) GenerateErrorResponse(err error) opResponse {
	return &CreateResponse{
		RequestID: cr.RequestID,
		Error:     err,
	}
}

type CreateResponse struct {
	opResponseBase
	RequestID requestID `msgpack:requestID`
	Error     error     `msgpack:error`
}

func sl_create(ctx context.Context, req opRequest, callback func()) opResponse {
	if callback != nil {
		defer callback()
	}
	cr, ok := req.(*Sl_createRequest)
	if !ok {
		return &CreateResponse{
			RequestID: cr.RequestID,
			Error:     errBadRequest,
		}
	}

	select {
	case <-ctx.Done():
		return cr.GenerateErrorResponse(errInterrupted)
	default:
	}

	file, err := os.OpenFile(cr.Path, os.O_CREATE, cr.Mode)
	if err != nil {
		return cr.GenerateErrorResponse(err)
	}
	defer func() {
		_ = file.Close()
	}()
	return cr.GenerateErrorResponse(err)
}

// Stateless truncate
type Sl_truncRequest struct {
	RequestID requestID `msgpack:requestID`
	Path      string    `msgpack:path`
	Size      int64     `msgpack:size`
}

func (tr *Sl_truncRequest) GetRequestID() requestID {
	return tr.RequestID
}

func (tr *Sl_truncRequest) GenerateErrorResponse(err error) opResponse {
	return &TruncResponse{
		RequestID: tr.RequestID,
		Error:     err,
	}
}

type TruncResponse struct {
	opResponseBase
	RequestID requestID `msgpack:requestID`
	Error     error     `msgpack:error`
}

func sl_trunc(ctx context.Context, req opRequest, callback func()) opResponse {
	if callback != nil {
		defer callback()
	}
	tr, ok := req.(*Sl_truncRequest)
	if !ok {
		return &TruncResponse{
			RequestID: tr.RequestID,
			Error:     errBadRequest,
		}
	}

	select {
	case <-ctx.Done():
		return tr.GenerateErrorResponse(errInterrupted)
	default:
	}

	err := os.Truncate(tr.Path, tr.Size)
	return tr.GenerateErrorResponse(err)
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
