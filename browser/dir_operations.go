package browser

import (
	"context"
	"io/ioutil"
	"os"
)

func init() {
	localRegistry.registerCommand(SL_READDIR, sl_readdir)
	localRegistry.registerCommand(SL_MKDIR, sl_mkdir)
}

// listRequest implements opRequest
type Sl_readdirRequest struct {
	// Path of directory to list
	Path      string      `msgpack:path`
	RequestID requestID   `msgpack:requestID`
	Code      commandCode `msgpack:code`
}

func (l *Sl_readdirRequest) GetRequestID() requestID {
	return l.RequestID
}

func (l *Sl_readdirRequest) GenerateErrorResponse(err error) opResponse {
	return &ReaddirResponse{
		RequestID: l.RequestID,
		Error:     err,
	}
}

// readdirResponse implements opResponse
type ReaddirResponse struct {
	opResponseBase
	Code      commandCode `msgpack:code`
	RequestID requestID   `msgpack:requestID`
	Entries   []dirent    `msgpack:entries`
	Error     error       `msgpack:error`
}

func sl_readdir(ctx context.Context, req opRequest, callback func()) opResponse {
	if callback != nil {
		defer callback()
	}

	lr, ok := req.(*Sl_readdirRequest)
	if !ok {
		return &ReaddirResponse{
			RequestID: req.GetRequestID(),
			Error:     errBadRequest,
		}
	}
	// Check if operation was cancelled before executing call
	select {
	case <-ctx.Done():
		return lr.GenerateErrorResponse(errInterrupted)
	default:
	}

	fileInfo, err := ioutil.ReadDir(lr.Path)
	if err != nil {
		return lr.GenerateErrorResponse(err)
	}
	entries := []dirent{}
	for _, fi := range fileInfo {
		d := dirent{
			Name: fi.Name(),
			Dir:  fi.IsDir(),
			Size: fi.Size(),
			Mode: fi.Mode(),
		}
		entries = append(entries, d)
	}
	return &ReaddirResponse{
		RequestID: lr.RequestID,
		Entries:   entries,
	}
}

// Stateful readdir
type Sf_readdirRequest struct {
	RequestID requestID   `msgpack:requestID`
	Code      commandCode `msgpack:code`
	HandleID  uint64      `msgpack:handleID`
}

func (r *Sf_readdirRequest) GetRequestID() requestID {
	return r.RequestID
}

func (r *Sf_readdirRequest) GenerateErrorResponse(err error) opResponse {
	return &ReaddirResponse{
		RequestID: r.RequestID,
		Error:     err,
	}
}

func sf_readdir(ctx context.Context, req opRequest, callback func()) opResponse {
	if callback != nil {
		defer callback()
	}
	rr, ok := req.(*Sf_readdirRequest)
	if !ok {
		return &ReaddirResponse{
			RequestID: req.GetRequestID(),
			Error:     errBadRequest,
		}
	}
	fh, ok := localFileRegistry.get(rr.HandleID)
	if !ok {
		return rr.GenerateErrorResponse(errInvalidHandle)
	}

	select {
	case <-ctx.Done():
		return rr.GenerateErrorResponse(errInterrupted)
	default:
	}

	fileInfo, err := fh.file.Readdir(0)
	if err != nil {
		return rr.GenerateErrorResponse(err)
	}
	entries := []dirent{}
	for _, fi := range fileInfo {
		d := dirent{
			Name: fi.Name(),
			Dir:  fi.IsDir(),
			Size: fi.Size(),
			Mode: fi.Mode(),
		}
		entries = append(entries, d)
	}
	return &ReaddirResponse{
		RequestID: rr.RequestID,
		Entries:   entries,
	}
}

// mkdir syscall. Stateless only.
type Sl_mkdirRequest struct {
	RequestID requestID `msgpack:requestID`
	// Path where to create the new directory
	Path string `msgpack:path`
	// FileMode is just an unsigned 32 bit integer
	Mode os.FileMode `msgpack:mode`
}

func (m *Sl_mkdirRequest) GetRequestID() requestID {
	return m.RequestID
}

func (m *Sl_mkdirRequest) GenerateErrorResponse(err error) opResponse {
	return &Sl_mkdirResponse{
		RequestID: m.RequestID,
		Error:     err,
	}
}

type Sl_mkdirResponse struct {
	opResponseBase
	RequestID requestID `msgpack:requestID`
	Error     error     `msgpack:error`
	Attr      attr      `msgpack:attr,omitempty`
}

func sl_mkdir(ctx context.Context, req opRequest, callback func()) opResponse {
	if callback != nil {
		defer callback()
	}
	mr, ok := req.(*Sl_mkdirRequest)
	if !ok {
		return &Sl_mkdirResponse{
			RequestID: req.GetRequestID(),
			Error:     errBadRequest,
		}
	}

	select {
	case <-ctx.Done():
		return mr.GenerateErrorResponse(errInterrupted)
	default:
	}

	err := os.Mkdir(mr.Path, mr.Mode)
	if err != nil {
		return mr.GenerateErrorResponse(err)
	}
	attr, err := getAttr(mr.Path)
	if err != nil {
		return mr.GenerateErrorResponse(err)
	}
	return &Sl_mkdirResponse{
		RequestID: mr.RequestID,
		Attr:      attr,
	}
}
