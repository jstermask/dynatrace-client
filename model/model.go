package model

type DynatraceExtensionCreateRequest struct {
	Payload string
}

type DynatraceExtensionGetBinaryRequest struct {
	Id string
}

type DynatraceExtensionGetBinaryResponse struct {
	Id string
	Payload string
}

type DynatraceExtensionCreateResponse struct {
	Id          string
	Name        string
	Description string
	Error       *DynatraceError
}

type DynatraceError struct {
	Code    string
	Message string
}

type DynatraceExtensionMetadata struct {
	Name string `json:"name"`
}
