package browser

// InterruptRequest is used to interrupt a running
// operation. It should contain the requestID of the
// operation to be interrupted. InterruptRequests do
// not generate a response. The target operation will
// return errInterrupted if interruption takes place
// on time.
type InterruptRequest struct {
	RequestID requestID `msgpack:requestID`
	// IntrID is theID of the operation
	// which is to be interrupted.
	// If the request reaches the server after
	// the operation has executed then it is
	// ignored. Handle this as appropriate when
	// writing FUSE bindings.
	IntrID requestID `msgpack:intrID`
}

func (i *InterruptRequest) GetRequestID() requestID {
	return i.RequestID
}

func (i *InterruptRequest) GenerateErrorResponse(err error) opResponse {
	return nil
}

func op_interrupt(intr *InterruptRequest) {
	localRegistry.callCancelFunc(intr.IntrID)
}
