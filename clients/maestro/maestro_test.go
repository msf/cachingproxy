//go:build unit
// +build unit

package maestro

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/Unbabel/architecture-v2/flowrunner/common/go-kit/log"
)

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newHTTPClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(fn),
	}
}

const MTResponseStringForTest string = `
{
    "debug_info": {
        "build_output": null,
        "metrics": [
            {
                "end_time": 1586126622.3406258,
                "name": "maestro.chat_machine_translation.run",
                "start_time": 1586126621.6208394
            },
            {
                "end_time": 1586126636.7249038,
                "name": "maestro.chat_machine_translation.rebuild",
                "start_time": 1586126636.7129135
            }
        ],
        "metrics_build": [
            {
                "name": "maestro.build.input.text.size.chars",
                "timestamp": 1586126635.958016,
                "value": 73
            },
            {
                "end_time": 1586126635.9875088,
                "name": "maestro.build.extract_source_nuggets_for_mt",
                "start_time": 1586126635.9873843
            }
        ],
        "metrics_rebuild": [
            {
                "end_time": 1586126636.7177303,
                "name": "maestro.rebuild_task.process_result_from_mt",
                "start_time": 1586126636.7129586
            },
            {
                "name": "maestro.rebuild.translated_content.size.chars",
                "timestamp": 1586126636.7248845,
                "value": 72
            }
        ],
        "metrics_translate": [
            {
                "name": "maestro.machine_translate.input.source_nuggets.count",
                "timestamp": 1586126635.9875798,
                "value": 2
            },
            {
                "end_time": 1586126636.712757,
                "name": "maestro.machine_translate.translate",
                "start_time": 1586126636.2154226
            }
        ],
        "rebuild_output": null,
        "translate_output": null
    },
    "text": "Please translate this simple sentence.\nTranslate another sentence please.",
    "translated_content": "Por favor, traduzir esta frase simples.\nTraduzir outra frase, por favor.",
    "translated_data": {
        "job_engine": "unbabel-nmt",
        "job_engine_model_name": "chat",
        "job_engine_model_version": "2019-01-08T22:55:48Z",
        "nuggets": [
            {
							  "anonymization_uuid": "836cb66206fa486888151dc8d168baa6",
                "chunk": "5e8a5f2b0b3d2b0008b2a3b0",
                "chunk_name": null,
                "id": "5e8a5f2b0b3d2b0008b2a3b2",
                "meta_attributes":{ "handler":"client"},
                "mt_annotations": {},
                "mt_engine": "unbabel-nmt",
                "mt_markup": [],
                "mt_num_words": 6,
                "mt_text": "Por favor, traduzir esta frase simples.",
                "num_words": 5,
                "position": 0,
                "qe_alerts": null,
                "qe_score": -1,
                "rules": null,
                "text": "Please translate this simple sentence.",
                "text_annotations": {},
                "text_markup": [],
                "tm_entry_id": null,
                "tm_translation_id": null,
                "type": "text"
            },
            {
							  "anonymization_uuid": "836cb66206fa486888151dc8d168baa7",
                "chunk": "5e8a5f2b0b3d2b0008b2a3b1",
                "chunk_name": null,
                "id": "5e8a5f2b0b3d2b0008b2a3b3",
                "mt_annotations": {},
                "mt_engine": "unbabel-nmt",
                "mt_markup": [],
                "mt_num_words": 5,
                "mt_text": "Traduzir outra frase, por favor.",
                "num_words": 4,
                "position": 1,
                "qe_alerts": null,
                "qe_score": -1,
                "rules": null,
                "text": "Translate another sentence please.",
                "text_annotations": {},
                "text_markup": [],
                "tm_entry_id": null,
                "tm_translation_id": null,
                "type": "text"
            }
				]
    },
    "uid": "unbabel_machine_translation_flow"
}
`

const rebuildResponseStringForTest string = `
{
    "debug_info": {
        "build_output": null,
        "metrics": [
            {
                "end_time": 1586126636.7249038,
                "name": "maestro.chat_machine_translation.rebuild",
                "start_time": 1586126636.7129135
            }
        ],
        "metrics_build": [],
        "metrics_rebuild": [
            {
                "end_time": 1586126636.7177303,
                "name": "maestro.rebuild_task.process_result_from_mt",
                "start_time": 1586126636.7129586
            },
            {
                "name": "maestro.rebuild.translated_content.size.chars",
                "timestamp": 1586126636.7248845,
                "value": 72
            }
        ],
        "metrics_translate": [],
        "rebuild_output": null,
        "translate_output": null
    },
    "uid": "unbabel_rebuild_flow",
    "text": "Please translate this simple sentence.\nTranslate another sentence please.",
    "translated_content": "Por favor, traduzir esta frase simples.\nTraduzir outra frase, por favor.",
    "translated_data": {
        "job_engine": "unbabel-nmt",
        "job_engine_model_name": "chat",
        "job_engine_model_version": "2019-01-08T22:55:48Z",
        "nuggets": [
            {
							  "anonymization_uuid": "836cb66206fa486888151dc8d168baa6",
                "chunk": "5e8a5f2b0b3d2b0008b2a3b0",
                "chunk_name": null,
                "id": "5e8a5f2b0b3d2b0008b2a3b2",
                "meta_attributes":{ "handler":"client"},
                "mt_annotations": {},
                "mt_engine": "unbabel-nmt",
                "mt_markup": [],
                "mt_num_words": 6    ,
                "mt_text": "Por favor, traduzir esta frase simples.",
                "num_words": 5,
                "position": 0,
                "qe_alerts": null,
                "qe_score": -1,
                "rules": null,
                "text": "Please translate this simple sentence.",
                "text_annotations": {
                    "anonymization": {},
                    "notranslate": {}
                },
                "text_markup": [],
                "tm_entry_id": null,
                "tm_translation_id": null,
                "type": "text"
            },
            {
							  "anonymization_uuid": "836cb66206fa486888151dc8d168baa7",
                "chunk": "5e8a5f2b0b3d2b0008b2a3b0",
                "chunk_name": null,
                "id": "5e8a5f2b0b3d2b0008b2a3b3",
                "mt_annotations": {},
                "mt_engine": "unbabel-nmt",
                "mt_markup": [],
                "mt_num_words": 5,
                "mt_text": "Traduzir outra frase, por favor.",
                "num_words": 4,
                "position": 1,
                "qe_alerts": null,
                "qe_score": -1,
                "rules": null,
                "text": "Translate another sentence please.",
                "text_annotations": {
                    "anonymization": {},
                    "notranslate": {}
                },
                "text_markup": [],
                "tm_entry_id": null,
                "tm_translation_id": null,
                "type": "text"
            }
        ],
        "skeleton": "<ubid>5e8a5f2b0b3d2b0008b2a3b0</ubid>"
    }
}
`

func TestMachineTranslateHappyCase(t *testing.T) {
	responseBody := MTResponseStringForTest
	expectedRequestBody := `{
  "uid": "1245",
  "content_type": "chat",
  "source_language": "en",
  "target_language": "pt",
  "text": "Please translate this simple sentence.\nTranslate another sentence please.",
  "origin": "origin",
  "client_username": "username",
  "client_brand": "brand",
  "glossary_id": "123",
  "text_format": "text",
  "build_rebuild_config_json": "{\"annotation_categories\": [\"categories\"]}"
}`

	maestroClient := New(log.NewNopLogger(), "user", "pass", DefaultCharsPersSecondTimeout)

	// replace the *http.Client w/ one with overriden Transport
	maestroClient.httpClient.HTTPClient = newHTTPClient(
		func(req *http.Request) *http.Response {
			buf, _ := ioutil.ReadAll(req.Body)
			var prettyJSON bytes.Buffer
			_ = json.Indent(&prettyJSON, buf, "", "  ")
			assert.Equal(t, expectedRequestBody, prettyJSON.String())
			assert.Equal(t, req.URL.String(), "http://chat-mt-solo-google.maestro.svc.cluster.local/v1/mt")
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
				Header:     make(http.Header),
			}
		},
	)

	resp, err := maestroClient.MachineTranslate(context.TODO(),
		"http://chat-mt-solo-google.maestro.svc.cluster.local",
		&MTRequest{
			Text:                   "Please translate this simple sentence.\nTranslate another sentence please.",
			TextFormat:             "text",
			ContentType:            "chat",
			ClientBrand:            "brand",
			GlossaryID:             "123",
			ClientUsername:         "username",
			SourceLanguage:         "en",
			TargetLanguage:         "pt",
			Origin:                 "origin",
			UID:                    "1245",
			BuildRebuildConfigJSON: "{\"annotation_categories\": [\"categories\"]}",
		},
	)
	assert.Nil(t, err)
	assert.Equal(t, "unbabel_machine_translation_flow", resp.UID)
	assert.Equal(t, "Please translate this simple sentence.\nTranslate another sentence please.", resp.Text)
	assert.Equal(t, 2, len(resp.TranslatedData.Nuggets))
	assert.Equal(t, "client", resp.TranslatedData.Nuggets[0].MetaAttributes.Handler)
}

func TestMachineTranslateTicketHappyCase(t *testing.T) {
	responseBody := MTResponseStringForTest
	expectedRequestBody := `{
  "uid": "1245",
  "content_type": "ticket",
  "source_language": "en",
  "target_language": "pt",
  "text": "Please translate this simple sentence.\nTranslate another sentence please.",
  "origin": "origin",
  "client_username": "username",
  "client_brand": "brand",
  "glossary_id": "123",
  "text_format": "text",
  "build_rebuild_config_json": "{\"annotation_categories\": [\"categories\"]}"
}`

	ticketMTClient := New(log.NewNopLogger(), "user", "pass", DefaultCharsPersSecondTimeout)

	// replace the *http.Client w/ one with overriden Transport
	ticketMTClient.httpClient.HTTPClient = newHTTPClient(
		func(req *http.Request) *http.Response {
			buf, _ := ioutil.ReadAll(req.Body)
			var prettyJSON bytes.Buffer
			_ = json.Indent(&prettyJSON, buf, "", "  ")
			assert.Equal(t, expectedRequestBody, prettyJSON.String())
			assert.Equal(t, req.URL.String(), "http://ticket-mt-solo-google.maestro.svc.cluster.local/v1/mt")
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
				Header:     make(http.Header),
			}
		},
	)

	resp, err := ticketMTClient.MachineTranslate(context.TODO(),
		"http://ticket-mt-solo-google.maestro.svc.cluster.local",
		&MTRequest{
			Text:                   "Please translate this simple sentence.\nTranslate another sentence please.",
			TextFormat:             "text",
			ContentType:            "ticket",
			ClientBrand:            "brand",
			GlossaryID:             "123",
			ClientUsername:         "username",
			SourceLanguage:         "en",
			TargetLanguage:         "pt",
			Origin:                 "origin",
			UID:                    "1245",
			BuildRebuildConfigJSON: "{\"annotation_categories\": [\"categories\"]}",
		},
	)
	assert.Nil(t, err)
	assert.Equal(t, "unbabel_machine_translation_flow", resp.UID)
	assert.Equal(t, "Please translate this simple sentence.\nTranslate another sentence please.", resp.Text)
	assert.Equal(t, 2, len(resp.TranslatedData.Nuggets))
}

func TestMachineTranslateWithInvalidRequest(t *testing.T) {
	responseBody := `{}`
	expectedRequestBody := `{
  "uid": "1245",
  "content_type": "ticket",
  "source_language": "en",
  "target_language": "pt",
  "text": "",
  "origin": "origin",
  "client_username": "username",
  "client_brand": "brand",
  "glossary_id": "123",
  "text_format": "text",
  "build_rebuild_config_json": "{\"annotation_categories\": [\"categories\"]}"
}`

	ticketMTClient := New(log.NewNopLogger(), "user", "pass", DefaultCharsPersSecondTimeout)

	// replace the *http.Client w/ one with overriden Transport
	ticketMTClient.httpClient.HTTPClient = newHTTPClient(
		func(req *http.Request) *http.Response {
			buf, _ := ioutil.ReadAll(req.Body)
			var prettyJSON bytes.Buffer
			_ = json.Indent(&prettyJSON, buf, "", "  ")
			assert.Equal(t, expectedRequestBody, prettyJSON.String())
			assert.Equal(t, req.URL.String(), "http://ticket-mt-solo-google.maestro.svc.cluster.local/v1/mt")
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			return &http.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
				Header:     make(http.Header),
			}
		},
	)

	resp, err := ticketMTClient.MachineTranslate(context.TODO(),
		"http://ticket-mt-solo-google.maestro.svc.cluster.local",
		&MTRequest{
			TextFormat:             "text",
			ContentType:            "ticket",
			ClientBrand:            "brand",
			GlossaryID:             "123",
			ClientUsername:         "username",
			SourceLanguage:         "en",
			TargetLanguage:         "pt",
			Origin:                 "origin",
			UID:                    "1245",
			BuildRebuildConfigJSON: "{\"annotation_categories\": [\"categories\"]}",
		},
	)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, "received 422 response from maestro mt request", err.Error())
}

func TestMachineTranslateRetryLogic(t *testing.T) {
	maestroClient := New(log.NewNopLogger(), "user", "pass", DefaultCharsPersSecondTimeout)

	countRequests := 0
	// replace the *http.Client w/ one with overriden Transport
	maestroClient.httpClient.HTTPClient = newHTTPClient(
		func(req *http.Request) *http.Response {
			countRequests += 1
			return &http.Response{
				StatusCode: 503, // server is having a bad day...
				Body:       ioutil.NopCloser(bytes.NewBufferString("very busy")),
				Header:     make(http.Header),
			}
		},
	)

	resp, err := maestroClient.MachineTranslate(context.TODO(),
		"http://ticket-mt-solo-paypal-en-pt.maestro.svc.cluster.local",
		&MTRequest{
			Text:           "This is my text to translate",
			TextFormat:     "text",
			ContentType:    "chat",
			ClientBrand:    "brand",
			GlossaryID:     "123",
			ClientUsername: "username",
			SourceLanguage: "en",
			TargetLanguage: "pt",
			Origin:         "origin",
			UID:            "1245",
		},
	)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, countRequests-1, maestroClient.httpClient.RetryMax)
}

func TestRebuildHappyCase(t *testing.T) {
	expectedRequestBody := `{
  "uid": "unbabel_rebuild_flow",
  "source_language": "en",
  "target_language": "pt",
  "text": "Please translate this simple sentence.\nTranslate another sentence please.",
  "translated_data": {
    "job_engine": "unbabel-nmt",
    "job_engine_model_name": "chat",
    "job_engine_model_version": "2019-01-08T22:55:48Z",
    "nuggets": [
      {
        "anonymization_uuid": "836cb66206fa486888151dc8d168baa6",
        "chunk": "5e8a5f2b0b3d2b0008b2a3b0",
        "chunk_name": "",
        "id": "5e8a5f2b0b3d2b0008b2a3b2",
        "human_edition_metadata": {},
        "meta_attributes": {
          "handler": "client"
        },
        "mt_annotations": {
          "anonymization": {},
          "notranslate": {}
        },
        "mt_engine": "unbabel-nmt",
        "mt_markup": [
          {
            "tid": 0,
            "start": 0,
            "text": "\u003cp\u003e"
          }
        ],
        "mt_num_words": 6,
        "mt_text": "Por favor, traduzir esta frase simples.",
        "mt_text_tm": "Por favor, traduzir esta frase simples.",
        "num_words": 5,
        "position": 0,
        "qe_alerts": null,
        "qe_score": -1,
        "rules": "",
        "text": "Please translate this simple sentence.",
        "text_annotations": {
          "anonymization": {},
          "notranslate": {}
        },
        "text_markup": [
          {
            "tid": 0,
            "start": 0,
            "text": "\u003cp\u003e"
          }
        ],
        "text_tm": "Please translate this simple sentence.",
        "tm_confidence_score": 0,
        "tm_curated": false,
        "tm_entry_id": "",
        "tm_is_blocked_for_editors": false,
        "tm_is_visible_for_editors": false,
        "tm_match_by_brand": false,
        "tm_match_by_client": false,
        "tm_match_by_content_type": false,
        "tm_match_by_origin": false,
        "tm_translation_id": "",
        "type": "text",
        "tm_uses_placeholders_feature": false
      },
      {
        "anonymization_uuid": "836cb66206fa486888151dc8d168baa7",
        "chunk": "5e8a5f2b0b3d2b0008b2a3b0",
        "chunk_name": "",
        "id": "5e8a5f2b0b3d2b0008b2a3b3",
        "human_edition_metadata": {},
        "mt_annotations": {
          "anonymization": {},
          "notranslate": {}
        },
        "mt_engine": "unbabel-nmt",
        "mt_markup": [],
        "mt_num_words": 5,
        "mt_text": "Traduzir outra frase, por favor.",
        "mt_text_tm": "Traduzir outra frase, por favor.",
        "num_words": 4,
        "position": 1,
        "qe_alerts": null,
        "qe_score": -1,
        "rules": "",
        "text": "Translate another sentence please.",
        "text_annotations": {
          "anonymization": {},
          "notranslate": {}
        },
        "text_markup": [],
        "text_tm": "Translate another sentence please.",
        "tm_confidence_score": 0,
        "tm_curated": false,
        "tm_entry_id": "",
        "tm_is_blocked_for_editors": false,
        "tm_is_visible_for_editors": false,
        "tm_match_by_brand": false,
        "tm_match_by_client": false,
        "tm_match_by_content_type": false,
        "tm_match_by_origin": false,
        "tm_translation_id": "",
        "type": "text",
        "tm_uses_placeholders_feature": false
      }
    ],
    "skeleton": "\u003cubid\u003e5e8a5f2b0b3d2b0008b2a3b0\u003c/ubid\u003e",
    "source_num_words": 9,
    "target_num_words": 11,
    "maestro_version": ""
  }
}`
	responseBody := rebuildResponseStringForTest

	maestroClient := New(log.NewNopLogger(), "user", "pass", DefaultCharsPersSecondTimeout)

	// replace the *http.Client w/ one with overriden Transport
	maestroClient.httpClient.HTTPClient = newHTTPClient(
		func(req *http.Request) *http.Response {
			buf, _ := ioutil.ReadAll(req.Body)
			var prettyJSON bytes.Buffer
			_ = json.Indent(&prettyJSON, buf, "", "  ")
			assert.Equal(t, expectedRequestBody, prettyJSON.String())
			assert.Equal(t, req.URL.String(), "http://chat-mt-solo-test-chat.maestro.svc.cluster.local/v1/rebuild")
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
				Header:     make(http.Header),
			}
		},
	)

	position0, position1 := 0, 1
	resp, err := maestroClient.Rebuild(context.TODO(),
		"http://chat-mt-solo-test-chat.maestro.svc.cluster.local",
		&RebuildRequest{
			Text:           "Please translate this simple sentence.\nTranslate another sentence please.",
			SourceLanguage: "en",
			TargetLanguage: "pt",
			UID:            "unbabel_rebuild_flow",
			ContentType:    "chat",
			TranslatedData: TranslatedData{
				JobEngine:             "unbabel-nmt",
				JobEngineModelName:    "chat",
				JobEngineModelVersion: "2019-01-08T22:55:48Z",
				SourceNumWords:        9,
				TargetNumWords:        11,
				Nuggets: []Nugget{
					Nugget{
						AnonymizationUUID: "836cb66206fa486888151dc8d168baa6",
						Chunk:             "5e8a5f2b0b3d2b0008b2a3b0",
						ID:                "5e8a5f2b0b3d2b0008b2a3b2",
						MTEngine:          "unbabel-nmt",
						MetaAttributes:    &MetaAttributes{Handler: "client"},
						MTAnnotations: Annotations{
							"anonymization": map[string]Annotation{},
							"notranslate":   map[string]Annotation{},
						},
						MTMarkup: []MarkupTag{
							MarkupTag{
								TID:   0,
								Start: 0,
								Text:  "<p>",
							},
						},
						MTNumWords: 6,
						MTText:     "Por favor, traduzir esta frase simples.",
						MTTextTM:   "Por favor, traduzir esta frase simples.",
						NumWords:   5,
						Position:   position0,
						QEScore:    -1,
						Rules:      "",
						Text:       "Please translate this simple sentence.",
						TextAnnotations: Annotations{
							"anonymization": map[string]Annotation{},
							"notranslate":   map[string]Annotation{},
						},
						TextMarkup: []MarkupTag{
							MarkupTag{
								TID:   0,
								Start: 0,
								Text:  "<p>",
							},
						},
						TextTM:               "Please translate this simple sentence.",
						TMEntryID:            "",
						TMTranslationID:      "",
						Type:                 "text",
						HumanEditionMetadata: HumanEditionMetadata{},
					},
					Nugget{
						AnonymizationUUID: "836cb66206fa486888151dc8d168baa7",
						Chunk:             "5e8a5f2b0b3d2b0008b2a3b0",
						ID:                "5e8a5f2b0b3d2b0008b2a3b3",
						MTEngine:          "unbabel-nmt",
						MTAnnotations: Annotations{
							"anonymization": map[string]Annotation{},
							"notranslate":   map[string]Annotation{},
						},
						MTMarkup:   []MarkupTag{},
						MTNumWords: 5,
						MTText:     "Traduzir outra frase, por favor.",
						MTTextTM:   "Traduzir outra frase, por favor.",
						NumWords:   4,
						Position:   position1,
						QEScore:    -1,
						Rules:      "",
						Text:       "Translate another sentence please.",
						TextAnnotations: Annotations{
							"anonymization": map[string]Annotation{},
							"notranslate":   map[string]Annotation{},
						},
						TextMarkup:           []MarkupTag{},
						TextTM:               "Translate another sentence please.",
						TMEntryID:            "",
						TMTranslationID:      "",
						Type:                 "text",
						HumanEditionMetadata: HumanEditionMetadata{},
					},
				},
				Skeleton: "<ubid>5e8a5f2b0b3d2b0008b2a3b0</ubid>",
			},
		},
	)
	assert.Nil(t, err)
	assert.Equal(t, "unbabel_rebuild_flow", resp.UID)
	assert.Equal(t, "Please translate this simple sentence.\nTranslate another sentence please.", resp.Text)
	assert.Equal(t, "Por favor, traduzir esta frase simples.\nTraduzir outra frase, por favor.", resp.TranslatedContent)
	assert.Equal(t, "<ubid>5e8a5f2b0b3d2b0008b2a3b0</ubid>", resp.TranslatedData.Skeleton)
	assert.Equal(t, 2, len(resp.TranslatedData.Nuggets))
}

func TestRebuildWithInvalidRequest(t *testing.T) {
	expectedRequestBody := `{
  "uid": "unbabel_rebuild_flow",
  "source_language": "en",
  "target_language": "pt",
  "text": "Please translate this simple sentence.\nTranslate another sentence please.",
  "translated_data": {
    "job_engine": "unbabel-nmt",
    "job_engine_model_name": "chat",
    "job_engine_model_version": "2019-01-08T22:55:48Z",
    "nuggets": [
      {
        "anonymization_uuid": "836cb66206fa486888151dc8d168baa6",
        "chunk": "5e8a5f2b0b3d2b0008b2a3b0",
        "chunk_name": "",
        "id": "5e8a5f2b0b3d2b0008b2a3b2",
        "human_edition_metadata": {},
        "meta_attributes": {
          "handler": "client"
        },
        "mt_annotations": {
          "anonymization": {},
          "notranslate": {}
        },
        "mt_engine": "unbabel-nmt",
        "mt_markup": [
          {
            "tid": 0,
            "start": 0,
            "text": "\u003cp\u003e"
          }
        ],
        "mt_num_words": 6,
        "mt_text": "",
        "mt_text_tm": "Por favor, traduzir esta frase simples.",
        "num_words": 5,
        "position": 0,
        "qe_alerts": null,
        "qe_score": -1,
        "rules": "",
        "text": "Please translate this simple sentence.",
        "text_annotations": {
          "anonymization": {},
          "notranslate": {}
        },
        "text_markup": [
          {
            "tid": 0,
            "start": 0,
            "text": "\u003cp\u003e"
          }
        ],
        "text_tm": "Please translate this simple sentence.",
        "tm_confidence_score": 0,
        "tm_curated": false,
        "tm_entry_id": "",
        "tm_is_blocked_for_editors": false,
        "tm_is_visible_for_editors": false,
        "tm_match_by_brand": false,
        "tm_match_by_client": false,
        "tm_match_by_content_type": false,
        "tm_match_by_origin": false,
        "tm_translation_id": "",
        "type": "text",
        "tm_uses_placeholders_feature": false
      },
      {
        "anonymization_uuid": "836cb66206fa486888151dc8d168baa7",
        "chunk": "5e8a5f2b0b3d2b0008b2a3b0",
        "chunk_name": "",
        "id": "5e8a5f2b0b3d2b0008b2a3b3",
        "human_edition_metadata": {},
        "mt_annotations": {
          "anonymization": {},
          "notranslate": {}
        },
        "mt_engine": "unbabel-nmt",
        "mt_markup": [],
        "mt_num_words": 5,
        "mt_text": "Traduzir outra frase, por favor.",
        "mt_text_tm": "Traduzir outra frase, por favor.",
        "num_words": 4,
        "position": 1,
        "qe_alerts": null,
        "qe_score": -1,
        "rules": "",
        "text": "Translate another sentence please.",
        "text_annotations": {
          "anonymization": {},
          "notranslate": {}
        },
        "text_markup": [],
        "text_tm": "Translate another sentence please.",
        "tm_confidence_score": 0,
        "tm_curated": false,
        "tm_entry_id": "",
        "tm_is_blocked_for_editors": false,
        "tm_is_visible_for_editors": false,
        "tm_match_by_brand": false,
        "tm_match_by_client": false,
        "tm_match_by_content_type": false,
        "tm_match_by_origin": false,
        "tm_translation_id": "",
        "type": "text",
        "tm_uses_placeholders_feature": false
      }
    ],
    "skeleton": "\u003cubid\u003e5e8a5f2b0b3d2b0008b2a3b0\u003c/ubid\u003e",
    "source_num_words": 9,
    "target_num_words": 11,
    "maestro_version": ""
  }
}`
	responseBody := `{}`

	maestroClient := New(log.NewNopLogger(), "user", "pass", DefaultCharsPersSecondTimeout)

	// replace the *http.Client w/ one with overriden Transport
	maestroClient.httpClient.HTTPClient = newHTTPClient(
		func(req *http.Request) *http.Response {
			buf, _ := ioutil.ReadAll(req.Body)
			var prettyJSON bytes.Buffer
			_ = json.Indent(&prettyJSON, buf, "", "  ")
			assert.Equal(t, expectedRequestBody, prettyJSON.String())
			assert.Equal(t, req.URL.String(), "http://chat-mt-solo-test-chat.maestro.svc.cluster.local/v1/rebuild")
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			return &http.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
				Header:     make(http.Header),
			}
		},
	)

	position0, position1 := 0, 1
	resp, err := maestroClient.Rebuild(context.TODO(),
		"http://chat-mt-solo-test-chat.maestro.svc.cluster.local",
		&RebuildRequest{
			Text:           "Please translate this simple sentence.\nTranslate another sentence please.",
			SourceLanguage: "en",
			TargetLanguage: "pt",
			UID:            "unbabel_rebuild_flow",
			ContentType:    "chat",
			TranslatedData: TranslatedData{
				JobEngine:             "unbabel-nmt",
				JobEngineModelName:    "chat",
				JobEngineModelVersion: "2019-01-08T22:55:48Z",
				SourceNumWords:        9,
				TargetNumWords:        11,
				Nuggets: []Nugget{
					Nugget{
						AnonymizationUUID: "836cb66206fa486888151dc8d168baa6",
						Chunk:             "5e8a5f2b0b3d2b0008b2a3b0",
						ID:                "5e8a5f2b0b3d2b0008b2a3b2",
						MTEngine:          "unbabel-nmt",
						MetaAttributes:    &MetaAttributes{Handler: "client"},
						MTAnnotations: Annotations{
							"anonymization": map[string]Annotation{},
							"notranslate":   map[string]Annotation{},
						},
						MTMarkup: []MarkupTag{
							MarkupTag{
								TID:   0,
								Start: 0,
								Text:  "<p>",
							},
						},
						MTNumWords: 6,
						MTText:     "",
						MTTextTM:   "Por favor, traduzir esta frase simples.",
						NumWords:   5,
						Position:   position0,
						QEScore:    -1,
						Rules:      "",
						Text:       "Please translate this simple sentence.",
						TextAnnotations: Annotations{
							"anonymization": map[string]Annotation{},
							"notranslate":   map[string]Annotation{},
						},
						TextMarkup: []MarkupTag{
							MarkupTag{
								TID:   0,
								Start: 0,
								Text:  "<p>",
							},
						},
						TextTM:               "Please translate this simple sentence.",
						TMEntryID:            "",
						TMTranslationID:      "",
						Type:                 "text",
						HumanEditionMetadata: HumanEditionMetadata{},
					},
					Nugget{
						AnonymizationUUID: "836cb66206fa486888151dc8d168baa7",
						Chunk:             "5e8a5f2b0b3d2b0008b2a3b0",
						ID:                "5e8a5f2b0b3d2b0008b2a3b3",
						MTEngine:          "unbabel-nmt",
						MTAnnotations: Annotations{
							"anonymization": map[string]Annotation{},
							"notranslate":   map[string]Annotation{},
						},
						MTMarkup:   []MarkupTag{},
						MTNumWords: 5,
						MTText:     "Traduzir outra frase, por favor.",
						MTTextTM:   "Traduzir outra frase, por favor.",
						NumWords:   4,
						Position:   position1,
						QEScore:    -1,
						Rules:      "",
						Text:       "Translate another sentence please.",
						TextAnnotations: Annotations{
							"anonymization": map[string]Annotation{},
							"notranslate":   map[string]Annotation{},
						},
						TextMarkup:           []MarkupTag{},
						TextTM:               "Translate another sentence please.",
						TMEntryID:            "",
						TMTranslationID:      "",
						Type:                 "text",
						HumanEditionMetadata: HumanEditionMetadata{},
					},
				},
				Skeleton: "<ubid>5e8a5f2b0b3d2b0008b2a3b0</ubid>",
			},
		},
	)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, "received 422 response from maestro rebuild request", err.Error())
}

func TestDurationPerTextSize(t *testing.T) {
	assert.Equal(t, 1*time.Second, durationForTextSize(10, 10.0))
	assert.Equal(t, 1*time.Millisecond, durationForTextSize(1, 1000.0))
	assert.Equal(t, 500*time.Millisecond, durationForTextSize(20, 40.0))
	assert.Equal(t, 50*time.Second, durationForTextSize(2000, 40.0))
}

func TestProportionalRequestTimeout(t *testing.T) {
	type test struct {
		path        string
		text        string
		charsPerSec float64
		expected    time.Duration
	}
	text := "this is my text that needs to be translated"
	shortText := "hi"
	duration := durationForTextSize(len(text), 1.0)
	tests := []test{
		{text: text, charsPerSec: 1.0, path: machineTranslatePath, expected: duration},
		{text: text, charsPerSec: 1.0, path: machineTranslateWithQualityEstimationPath, expected: 2 * duration},
		{text: text, charsPerSec: 1.0, path: pivotedMachineTranslatePath, expected: 2 * duration},
		{text: text, charsPerSec: 1.0, path: pivotedMachineTranslateWithQualityEstimationPath, expected: 4 * duration},
		{text: shortText, charsPerSec: 1.0, path: machineTranslatePath, expected: MinimumRequestTimeout},
		{text: shortText, charsPerSec: 1.0, path: machineTranslateWithQualityEstimationPath, expected: MinimumRequestTimeout},
		{text: shortText, charsPerSec: 1.0, path: pivotedMachineTranslatePath, expected: MinimumRequestTimeout},
		{text: shortText, charsPerSec: 1.0, path: pivotedMachineTranslateWithQualityEstimationPath, expected: MinimumRequestTimeout},
		{text: "", charsPerSec: 1.0, path: machineTranslatePath, expected: MinimumRequestTimeout},
	}

	for _, i := range tests {
		timeout := getTimeoutForRequestPayload(i.path, i.text, i.charsPerSec)
		assert.Equal(t, i.expected, timeout, fmt.Sprintf("%+v", i))

	}
}
