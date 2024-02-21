package model


type DynatraceExtensionRequest struct {
	Payload     string
}

type DynatraceExtensionResponse struct {
	Id string
	Name string
	Description string
	Error *DynatraceError
}

type DynatraceError struct {
	Code string
	Message string
}

type DynatraceExtensionMetadata struct {
	Name string `json:"name"`
}

