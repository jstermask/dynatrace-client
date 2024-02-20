package model


type DynatraceExtensionRequest struct {
	Name		string 
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