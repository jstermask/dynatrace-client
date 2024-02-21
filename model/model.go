package model


type DynatraceExtensionRequest struct {
	Payload     string
}

type DynatraceExtensionResponse struct {
	Id string
	Name string
	Description string
}

type DynatraceExtensionMetadata struct {
	Name string `json:"name"`
}