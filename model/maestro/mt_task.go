package maestro

type MTTaskRequest struct {
	SourceLanguage string   `json:"source_language,omitempty"`
	TargetSegment  string   `json:"target_segment,omitempty"`
	UID            string   `json:"uid,omitempty"`
	Origin         string   `json:"origin,omitempty"`
	ClientBrand    string   `json:"client_brand,omitempty"`
	Tone           string   `json:"tone,omitempty"`
	ContentType    string   `json:"content_type,omitempty"`
	MaestroNuggets []Nugget `json:"maestro_nuggets,omitempty"`
}

type TranslatedNuggets struct {
	JobEngine             string   `json:"job_engine,omitempty"`
	JobEngineModelName    string   `json:"job_engine_model_name,omitempty"`
	JobEngineModelVersion string   `json:"job_engine_model_version,omitempty"`
	Nuggets               []Nugget `json:"nuggets,omitempty"`
}

type MTTaskResponse struct {
	UID               string            `json:"uid,omitempty"`
	TranslatedNuggets TranslatedNuggets `json:"translated_nuggets,omitempty"`
}
