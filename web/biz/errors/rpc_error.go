package errors

const RPCErrorMsg = "RPC error"

type RPCError struct {
	cause error
}

func NewRPCError(cause error) RPCError {
	return RPCError{cause: cause}
}

func (r RPCError) Error() string {
	return RPCErrorMsg + ":" + r.cause.Error()
}
