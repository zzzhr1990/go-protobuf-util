package httprpc

/*
{
    "message": "missing param identity",
    "reference": "INVALID_ARGUMENT",
    "rpc": "InvalidArgument",
    "status": 400,
    "success": false
}
*/

// ErrorResponse comon error
type ErrorResponse struct {
	Message   string `json:"message,omitempty"`
	Reference string `json:"reference,omitempty"`
	Status    int32  `json:"status,omitempty"`
	Success   bool   `json:"success,omitempty"`
}
