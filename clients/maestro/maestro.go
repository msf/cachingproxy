package maestro

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/msf/cachingproxy/model/maestro"
	"github.com/pkg/errors"
)

// Maestro is the interface for machine translation services in arch2
type Maestro interface {
	MTTask(ctx context.Context,
		serviceURL string,
		req *maestro.MTTaskRequest) (*maestro.MTTaskResponse, error)
}

type maestroClient struct {
	authUsername string
	authPassword string
	// MT & MTQE request durations are proportional to character (utf8) input size
	charsPerSecondTimeout float64
	logger                log.Logger
	httpClient            *retryablehttp.Client
}

const (
	DefaultCharsPersSecondTimeout                    = 30 // extreme lower bound of throughput for our MT engines
	MinimumRequestTimeout                            = 20 * time.Second
	machineTranslatePath                             = "v1/mt"
	machineTranslateWithQualityEstimationPath        = "v1/mt_qe"
	pivotedMachineTranslatePath                      = "v1/pivoted_mt"
	pivotedMachineTranslateWithQualityEstimationPath = "v1/pivoted_mt_qe"
	rebuildPath                                      = "v1/rebuild"
	pivotedRebuildPath                               = "v1/pivoted_rebuild"
	defaultRetryDelayMin                             = 50 * time.Millisecond // the retryablehttp default is 1 second..
	defaultRetryMax                                  = 3                     // 4 reqs in total

	metricTimingMT              = "chat2_timing_maestro_mt_secs"
	metricTimingMTWithQE        = "chat2_timing_maestro_mt_with_qe_secs"
	metricTimingPivotedMT       = "chat2_timing_maestro_pivoted_mt_secs"
	metricTimingPivotedMTWithQE = "chat2_timing_maestro_pivoted_mt_with_qe_secs"
	metricTimingRebuild         = "chat2_timing_maestro_rebuild_secs"
)

// New returns a maestroClient with request timeouts, cancellation and retry logic
func New(
	logger log.Logger,
	basicAuthUser,
	basicAuthPass string,
	charsPerSecondTimeout float64,
) *maestroClient {

	if basicAuthUser == "" || basicAuthPass == "" || charsPerSecondTimeout <= 0 {
		logger.Log("maestro", "Bad Arguments",
			"basicAuthUser", basicAuthUser,
			"basicAuthPass", basicAuthPass,
			"charsPerSecondTimeout", charsPerSecondTimeout,
		)
		return nil
	}
	m := &maestroClient{
		authUsername:          basicAuthUser,
		authPassword:          basicAuthPass,
		charsPerSecondTimeout: charsPerSecondTimeout,
		httpClient:            retryablehttp.NewClient(),
		logger:                logger,
	}
	m.httpClient.RetryWaitMin = defaultRetryDelayMin
	m.httpClient.RetryMax = defaultRetryMax
	return m
}

func (c *maestroClient) MTTask(
	ctx context.Context, serviceURL string, req *maestro.MTTaskRequest) (*maestro.MTResponse, error) {

	startTime := time.Now()

	respBuffer, err := c.callMachineTranslate(ctx, serviceURL, req, machineTranslatePath)
	if err != nil {
		return nil, err
	}

	var resp maestro.MTResponse
	err = ParseMTResponse(respBuffer, &resp)
	if err != nil {
		err = errors.Wrap(err, "maestro response parsing failed")
		c.logger.Log("maestro", req, "error", err)
		return nil, err
	}
	return &resp, nil
}

func (c *maestroClient) callMachineTranslate(
	ctx context.Context,
	MTModelURL string,
	req *maestro.MTRequest,
	path string,
) ([]byte, error) {

	body, err := json.Marshal(req)
	if err != nil {
		err = errors.Wrap(err, "maestro request json marshal failed")
		c.logger.Log("client", "maestro", "uid", req.UID, "step", "json.Marshal", "error", err)
		return nil, err
	}

	if MTModelURL == "" {
		err := errors.New("invalid argument: MTModelURL cannot be empty")
		c.logger.Log("client", "maestro", "uid", req.UID, "argument", err)
		return nil, err
	}

	timeout := getTimeoutForRequestPayload(path, req.Text, c.charsPerSecondTimeout)
	c.logger.Log("client", "maestro", "uid", req.UID, "timeout", timeout)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r, err := retryablehttp.NewRequest(
		"POST",
		MTModelURL+"/"+path,
		bytes.NewBuffer(body),
	)
	if err != nil {
		err = errors.Wrap(err, "maestro http request creation failed")
		c.logger.Log("client", "maestro", "uid", req.UID, "error", err)
		return nil, err
	}
	r.WithContext(ctx)
	r.SetBasicAuth(c.authUsername, c.authPassword)
	r.Header.Set("Content-Type", "application/json")

	httpResp, err := c.httpClient.Do(r)
	if err != nil {
		err = errors.Wrap(err, "maestro mt http request failed")
		c.logger.Log("client", "maestro", "uid", req.UID, "error", err)
		return nil, err
	}
	defer httpResp.Body.Close()

	if !isValidResponseStatusCode(httpResp) {
		return nil, errors.Errorf("received %v response from maestro mt request",
			httpResp.StatusCode)
	}

	data, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		err = errors.Wrap(err, "maestro read mt http response failed")
		c.logger.Log("client", "maestro", "uid", req.UID, "error", err)
		return nil, err
	}
	return data, nil
}

func durationForTextSize(textLen int, charsPerSecond float64) time.Duration {
	millis := float64(textLen) / charsPerSecond * 1000.0
	return time.Duration(millis) * time.Millisecond
}

func getTimeoutForRequestPayload(path, text string, expectedCharsPerSecond float64) time.Duration {
	duration := durationForTextSize(len(text), expectedCharsPerSecond)
	if strings.Contains(path, "pivot") {
		duration *= 2 // involves two MTs
	}
	if strings.Contains(path, "mt_qe") {
		duration *= 2
	}
	if duration < MinimumRequestTimeout {
		return MinimumRequestTimeout
	}
	return duration
}

func isValidResponseStatusCode(resp *http.Response) bool {
	// Responses with 2XX are an OK response
	// Responses with 3XX, 5XX get retried and if the number of retries is reached they are failures
	// Responses with 4XX are not retried and are failures
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true
	}
	return false
}

// SetHTTPClient is used by tests to mock out the http client
func (c *maestroClient) SetHTTPClient(client *http.Client) {
	c.httpClient.HTTPClient = client
}
