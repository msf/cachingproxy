package model

type MachineTranslationRequest struct {
	ID       string            `json:"id,omitempty"`
	Segments []string          `json:"segments,omitempty"`
	Metadata MTRequestMetadata `json:"metadata,omitempty"`
}

func (m MachineTranslationRequest) HasError() error {
	return nil
}

type MTRequestMetadata struct {
	SourceLang string            `json:"source_lang,omitempty"`
	TargetLang string            `json:"target_lang,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type MachineTranslationResponse struct {
	RequestID       string            `json:"request_id,omitempty"`
	TargetSegments  []string          `json:"target_segments,omitempty"`
	RequestMetadata MTRequestMetadata `json:"request_metadata,omitempty"`
}
