package maestro

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

type MTRequest struct {
	UID                    string  `json:"uid" bson:"uid"`
	ContentType            string  `json:"content_type,omitempty" bson:"content_type"`
	SourceLanguage         string  `json:"source_language" bson:"source_language"`
	TargetLanguage         string  `json:"target_language" bson:"target_language"`
	Text                   string  `json:"text" bson:"text"`
	Origin                 string  `json:"origin" bson:"origin"`
	ClientUsername         string  `json:"client_username" bson:"client_username"`
	ClientBrand            string  `json:"client_brand,omitempty" bson:"client_brand"`
	GlossaryID             string  `json:"glossary_id,omitempty" bson:"glossary_id"`
	Tone                   string  `json:"tone,omitempty" bson:"tone"`
	TextFormat             string  `json:"text_format" bson:"text_format"`
	BuildRebuildConfigJSON string  `json:"build_rebuild_config_json" bson:"build_rebuild_config_json"`
	QualitySkipThreshold   float64 `json:"quality_skip_threshold,omitempty" bson:"quality_skip_threshold,omitempty"`
	QualityServiceEndpoint string  `json:"quality_service_endpoint,omitempty" bson:"quality_service_endpoint,omitempty"`
}

type GenericMTResponse interface {
	GetTranslatedContent() string
	GetTranslatedData() []TranslatedData
	GetDebugInfo() DebugInfo
}

type MTResponse struct {
	UID                       string         `json:"uid" bson:"uid"`
	Text                      string         `json:"text,omitempty" bson:"text"`
	TranslatedContent         string         `json:"translated_content" bson:"translated_content"`
	TranslatedData            TranslatedData `json:"translated_data" bson:"translated_data"`
	DebugInfo                 DebugInfo      `json:"debug_info,omitempty" bson:"debug_info"`
	QualityScore              float64        `json:"quality_score,omitempty" bson:"quality_score,omitempty"`
	CanSkipHumanEdition       bool           `json:"can_skip_human_edition,omitempty" bson:"can_skip_human_edition,omitempty"`
	CanSkipHumanEditionReason string         `json:"can_skip_reason,omitempty" bson:"can_skip_reason,omitempty"`
}

type Failure struct {
	Context  []FailureContext `json:"context,omitempty" bson:"context,omitempty"`
	Category string           `json:"category,omitempty" bson:"category,omitempty"`
}

type FailureContext struct {
	Reason string `json:"reason" bson:"reason"`
}
type TranslatedData struct {
	JobEngine             string   `json:"job_engine" bson:"job_engine"`
	JobEngineModelName    string   `json:"job_engine_model_name" bson:"job_engine_model_name"`
	JobEngineModelVersion string   `json:"job_engine_model_version" bson:"job_engine_model_version"`
	Nuggets               []Nugget `json:"nuggets" bson:"nuggets"`
	Skeleton              string   `json:"skeleton" bson:"skeleton"`
	SourceNumWords        int      `json:"source_num_words" bson:"source_num_words"`
	TargetNumWords        int      `json:"target_num_words" bson:"target_num_words"`
	HTMLFlow              string   `json:"html_flow,omitempty" bson:"html_flow,omitempty"`
	MaestroVersion        string   `json:"maestro_version" bson:"maestro_version"`
}

type HumanEditionMetadata struct {
	TranslationNeedsContext *bool `json:"translation_needs_context,omitempty" bson:"translation_needs_context,omitempty"`
}

type Nugget struct {
	AnonymizationUUID         string                 `json:"anonymization_uuid,omitempty" bson:"anonymization_uuid"`
	Chunk                     string                 `json:"chunk,omitempty" bson:"chunk"`
	ChunkSkeleton             string                 `json:"chunk_skeleton,omitempty" bson:"chunk_skeleton"`
	ChunkName                 string                 `json:"chunk_name" bson:"chunk_name"`
	ChunkExtra                map[string]interface{} `json:"chunk_extra,omitempty" bson:"chunk_extra"`
	ID                        string                 `json:"id,omitempty" bson:"id"`
	HumanAnnotations          Annotations            `json:"human_annotations,omitempty" bson:"human_annotations"`
	HumanEditionMetadata      HumanEditionMetadata   `json:"human_edition_metadata,omitempty" bson:"human_edition_metadata"`
	HumanMarkup               []MarkupTag            `json:"human_markup,omitempty" bson:"human_markup"`
	HumanNumWords             int                    `json:"human_num_words,omitempty" bson:"human_num_words"`
	HumanText                 string                 `json:"human_text,omitempty" bson:"human_text"`
	HumanTextTM               string                 `json:"human_text_tm,omitempty" bson:"human_text_tm"`
	MetaAttributes            *MetaAttributes        `json:"meta_attributes,omitempty" bson:"meta_attributes"`
	MTAnnotations             Annotations            `json:"mt_annotations" bson:"mt_annotations"`
	MTEngine                  string                 `json:"mt_engine" bson:"mt_engine"`
	MTMarkup                  []MarkupTag            `json:"mt_markup" bson:"mt_markup"`
	MTNumWords                int                    `json:"mt_num_words" bson:"mt_num_words"`
	MTText                    string                 `json:"mt_text" bson:"mt_text"`
	MTTextTM                  string                 `json:"mt_text_tm" bson:"mt_text_tm"`
	NumWords                  int                    `json:"num_words" bson:"num_words"`
	Position                  int                    `json:"position" bson:"position"`
	QEAlerts                  []QEAlert              `json:"qe_alerts" bson:"qe_alerts"`
	QEScore                   float64                `json:"qe_score" bson:"qe_score"`
	Rules                     string                 `json:"rules" bson:"rules"`
	Text                      string                 `json:"text" bson:"text"`
	TextNoRespace             string                 `json:"text_no_respace,omitempty" bson:"text_no_respace"`
	TextAnnotations           Annotations            `json:"text_annotations" bson:"text_annotations"`
	TextMarkup                []MarkupTag            `json:"text_markup" bson:"text_markup"`
	TextMarkupNoRespace       []MarkupTag            `json:"text_markup_no_respace,omitempty" bson:"text_markup_no_respace"`
	TextTM                    string                 `json:"text_tm" bson:"text_tm"`
	TMConfidenceScore         float64                `json:"tm_confidence_score" bson:"tm_confidence_score"`
	TMCurated                 bool                   `json:"tm_curated" bson:"tm_curated"`
	TMEntryID                 string                 `json:"tm_entry_id" bson:"tm_entry_id"`
	TMIsBlockedForEditors     bool                   `json:"tm_is_blocked_for_editors" bson:"tm_is_blocked_for_editors"`
	TMIsVisibleForEditors     bool                   `json:"tm_is_visible_for_editors" bson:"tm_is_visible_for_editors"`
	TMMatchByBrand            bool                   `json:"tm_match_by_brand" bson:"tm_match_by_brand"`
	TMMatchByClient           bool                   `json:"tm_match_by_client" bson:"tm_match_by_client"`
	TMMatchByContentType      bool                   `json:"tm_match_by_content_type" bson:"tm_match_by_content_type"`
	TMMatchByOrigin           bool                   `json:"tm_match_by_origin" bson:"tm_match_by_origin"`
	TMTranslationID           string                 `json:"tm_translation_id" bson:"tm_translation_id"`
	Type                      string                 `json:"type" bson:"type"`
	TMUsesPlaceholdersFeature bool                   `json:"tm_uses_placeholders_feature" bson:"tm_uses_placeholders_feature"`
}

// Annotations is a nested dictionary with dict with:
//  - key=annotation name (for example: "notranslate", "number", "glossary", "anonymization")
//  - value=dict with:
//      - key=int, start position in the text this annotation refers to
//      - value=Annotation, where not all fields are present for all annotation types
type Annotations map[string]map[string]Annotation

// maestro annotations are a polymorphic type that is used for all types of annotations
// not all instances use all the fields fields
type Annotation struct {
	ID     string `json:"id" bson:"id"`
	Start  int    `json:"start" bson:"start"`
	End    int    `json:"end" bson:"end"`
	String string `json:"string" bson:"string"`

	// used by anonymization and  notranslate annotations
	Placeholder string `json:"placeholder,omitempty" bson:"placeholder,omitempty"`

	// used by anonymization annotations
	AnnotationType string `json:"annotation_type,omitempty" bson:"annotation_type,omitempty"`

	// used by glossary annotations
	Translation string `json:"translation,omitempty" bson:"translation,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}

type MarkupTag struct {
	TID   int    `json:"tid" bson:"tid"`
	Start int    `json:"start" bson:"start"`
	Text  string `json:"text" bson:"text"`
}

// MetaAttributeHandler is the only MetaAttribute known
// identifies if a nugget came from User or Agent on a chat
const MetaAttributeHandler = "handler"

type MetaAttributes struct {
	Handler string `json:"handler,omitempty" bson:"handler"`
}

type QEAlert struct {
	Type     string `json:"type" bson:"type"`
	Token    string `json:"token" bson:"token"`
	Position []int  `json:"position" bson:"position"`
}

type DebugInfo struct {
	MetricsBuild     []MTMetric  `json:"metrics_build,omitempty" bson:"metrics_build"`
	MetricsTranslate []MTMetric  `json:"metrics_translate,omitempty" bson:"metrics_translate"`
	MetricsRebuild   []MTMetric  `json:"metrics_rebuild,omitempty" bson:"metrics_rebuild"`
	Metrics          []MTMetric  `json:"metrics,omitempty" bson:"metrics"`
	BuildOutput      interface{} `json:"build_output,omitempty" bson:"build_output"`
	TranslateOutput  interface{} `json:"translate_output,omitempty" bson:"translate_output"`
	RebuildOutput    interface{} `json:"rebuild_output,omitempty" bson:"rebuild_output"`
}

type PivotedDebugInfo struct {
	MetricsBuild           []MTMetric `json:"metrics_build,omitempty" bson:"metrics_build"`
	MetricsFirstTranslate  []MTMetric `json:"metrics_first_translate,omitempty" bson:"metrics_first_translate"`
	MetricsPivot           []MTMetric `json:"metrics_pivot,omitempty" bson:"metrics_pivot"`
	MetricsSecondTranslate []MTMetric `json:"metrics_second_translate,omitempty" bson:"metrics_second_translate"`
	MetricsRebuild         []MTMetric `json:"metrics_rebuild,omitempty" bson:"metrics_rebuild"`
	Metrics                []MTMetric `json:"metrics,omitempty" bson:"metrics"`
}

type MTMetric struct {
	Name          string  `json:"name"`
	StartTime     float64 `json:"start_time"`
	EndTime       float64 `json:"end_time"`
	ElapsedMillis int     `json:"elapsed_millis"`
}

// ParseMTResponse can decode the complex mt_response
// TODO(msf): handle the nested dicts inside Nuggets correctly
func ParseMTResponse(respBuffer []byte, resp *MTResponse) error {
	err := json.NewDecoder(bytes.NewReader(respBuffer)).Decode(&resp)
	if err != nil {
		return errors.Wrapf(err, "original payload: %v", string(respBuffer))
	}
	return nil
}

// ParsePivotedMTResponse can decode the complex pivoted mt_response
// TODO(msf): handle the nested dicts inside Nuggets correctly
func ParsePivotedMTResponse(respBuffer []byte, resp *PivotedMTResponse) error {
	err := json.NewDecoder(bytes.NewReader(respBuffer)).Decode(&resp)
	if err != nil {
		return errors.Wrapf(err, "json decode error, payload: %v", string(respBuffer))
	}
	return nil
}

// ParseRebuildResponse can decode the complex rebuild response
// TODO(msf): handle the nested dicts inside Nuggets correctly
func ParseRebuildResponse(respBuffer []byte, resp *RebuildResponse) error {
	err := json.NewDecoder(bytes.NewReader(respBuffer)).Decode(&resp)
	if err != nil {
		return errors.Wrapf(err, "json decode error, payload: %v", string(respBuffer))
	}
	return nil
}

func (data TranslatedData) QEValue() float64 {
	// TODO(ak): Use document level QE score from Maestro, do not recompute
	translatedNuggets := data.Nuggets
	if len(translatedNuggets) < 1 {
		return 0.0
	}
	qeScore := translatedNuggets[0].QEScore
	// NOTE(ak): Document level QE score is minium (and not average) of nuggets' scores
	for _, nu := range translatedNuggets {
		if nu.QEScore < qeScore {
			qeScore = nu.QEScore
		}
	}
	if qeScore < 0 {
		return 0.0
	}
	return qeScore
}

func (data TranslatedData) IsTMOrUntranslatable() bool {
	translatedNuggets := data.Nuggets

	if len(translatedNuggets) < 1 {
		return false
	}

	for _, nu := range translatedNuggets {
		if nu.MTEngine != "tm" && nu.MTEngine != "untranslatable" {
			return false
		}
	}
	return true
}

func (data TranslatedData) NumberWords() int {
	count := 0
	for _, nu := range data.Nuggets {
		count += nu.NumWords
	}
	return count
}

func (data TranslatedData) HasTMs() bool {
	translatedNuggets := data.Nuggets
	if len(translatedNuggets) < 1 {
		return false
	}
	for _, nu := range translatedNuggets {
		if nu.MTEngine != "tm" {
			return false
		}
	}
	return true
}

func (data TranslatedData) IsMarkupHumanEditable() bool {
	return data.HTMLFlow == "inline" || data.HTMLFlow == "markup_aligner"
}

func (data TranslatedData) GetTranslationType() string {
	if data.HasTMs() {
		return string(translationtype.TM)
	}
	return string(translationtype.MT)
}
