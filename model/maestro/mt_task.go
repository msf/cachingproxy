package maestro

type MTTaskRequest struct {
	SourceLanguage string
	TargetSegment  string
	UID            string
	Origin         string
	ClientBrand    string
	Tone           string
	ContentType    string
	MaestroNuggets []Nugget
}

type TranslatedNuggets struct {
	JobEngine             string
	JobEngineModelName    string
	JobEngineModelVersion string
	Nuggets               []Nugget
}

type MTTaskResponse struct {
	UID               string
	TranslatedNuggets TranslatedNuggets
}
