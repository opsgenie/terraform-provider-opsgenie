package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	BaseURL     = "https://api.opsgenie.com"
	Endpoint    = "v2/alerts"
	EndpointURL = BaseURL + "/" + Endpoint
	BadEndpoint = ":"
)

type Team struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type Log struct {
	Owner      string `json:"owner"`
	CreateDate string `json:"createdDate"`
	Log        string `json:"log"`
}

type ResultWithoutDataField struct {
	ResultMetadata
	Result string `json:"result"`
}

type aResultDoesNotWantDataFieldsToBeParsed struct {
	ResultMetadata
	Logs   []Log  `json:"logs"`
	Offset string `json:"offset"`
}

type aResultWantsDataFieldsToBeParsed struct {
	ResultMetadata
	Teams []Team `json:"data"`
}

type testRequest struct {
	BaseRequest
	MandatoryField string
	ExtraField     string
}

func (tr testRequest) Validate() error {
	if tr.MandatoryField == "" {
		return errors.New("mandatory field cannot be empty")
	}

	return nil
}

func (tr testRequest) ResourcePath() string {
	return "/an-enpoint"
}

func (tr testRequest) Method() string {
	return http.MethodPost
}

type testResult struct {
	ResultMetadata
	Data string
}

func TestParsingWithDataField(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `
			{
    "data": [
        {
            "id": "1",
            "name": "n1",
            "description": "d1"
        },
        {
            "id": "2",
            "name": "n2",
            "description": "d2"
        },
        {
            "id": "3",
            "name": "n3",
            "description": "d3"
        }
    ],
    "took": 1.08,
    "requestId": "123"
}
		`)
	}))
	defer ts.Close()

	ogClient, err := NewOpsGenieClient(&Config{})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), errors.New("API key cannot be blank.").Error())

	ogClient, err = NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})
	assert.Nil(t, err)

	request := &testRequest{MandatoryField: "afield", ExtraField: "extra"}
	result := &aResultWantsDataFieldsToBeParsed{}

	err = ogClient.Exec(nil, request, result)
	assert.Nil(t, err)
	assert.Equal(t, result.Teams[0], Team{Id: "1", Name: "n1", Description: "d1"})
	assert.Equal(t, result.Teams[1], Team{Id: "2", Name: "n2", Description: "d2"})
	assert.Equal(t, result.Teams[2], Team{Id: "3", Name: "n3", Description: "d3"})
}

func TestParsingWithoutDataField(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `
			{
    "data": {
        "offset": "123",
        "logs": [
            {
                "owner": "o1",
                "createdDate": "c1",
                "log": "l1"
            },
            {
                "owner": "o2",
                "createdDate": "c2",
                "log": "l2"
            }
        ]
    },
    "took": 0.041,
    "requestId": "123"
}
		`)
	}))
	defer ts.Close()

	ogClient, err := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})
	assert.Nil(t, err)

	request := testRequest{MandatoryField: "afield", ExtraField: "extra"}
	result := &aResultDoesNotWantDataFieldsToBeParsed{}

	err = ogClient.Exec(nil, &request, result)
	assert.Nil(t, err)
	assert.Equal(t, result.Logs[0], Log{Owner: "o1", CreateDate: "c1", Log: "l1"})
	assert.Equal(t, result.Logs[1], Log{Owner: "o2", CreateDate: "c2", Log: "l2"})
	assert.Equal(t, result.Offset, "123")
}

func TestParsingWhenApiDoesNotReturnDataField(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `
			{
				"result": "processed",
				"requestId": "123",
				"took": 0.1
			}
		`)
	}))
	defer ts.Close()

	ogClient, err := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})
	assert.Nil(t, err)

	request := testRequest{MandatoryField: "afield", ExtraField: "extra"}
	result := &ResultWithoutDataField{}

	err = ogClient.Exec(nil, &request, result)
	assert.Nil(t, err)
	assert.Equal(t, "processed", result.Result)
}

func TestExec(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
    		"Data": "processed"}`)
	}))
	defer ts.Close()

	ogClient, err := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})
	assert.Nil(t, err)

	request := &testRequest{MandatoryField: "afield", ExtraField: "extra"}
	result := &testResult{}

	err = ogClient.Exec(nil, request, result)
	assert.Nil(t, err)
	assert.Equal(t, result.Data, "processed")
}

func TestParsingErrorExec(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer ts.Close()

	ogClient, err := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})
	assert.Nil(t, err)

	request := &testRequest{MandatoryField: "afield", ExtraField: "extra"}
	result := &testResult{}

	err = ogClient.Exec(nil, request, result)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Response could not be parsed, unexpected end of JSON input")
}

func TestExecWhenRequestIsNotValid(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
    		"Data": "processed"}`)
	}))
	defer ts.Close()

	ogClient, err := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})
	assert.Nil(t, err)

	request := testRequest{ExtraField: "extra"}
	result := &testResult{}

	err = ogClient.Exec(nil, &request, result)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "mandatory field cannot be empty")
}

func TestExecWhenApiReturns422(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintln(w, `{
    "message": "Request body is not processable. Please check the errors.",
    "errors": {
        "recipients#type": "Invalid recipient type 'bb'"
    },
    "took": 0.083,
    "requestId": "Id"
}`)
	}))
	defer ts.Close()

	ogClient, err := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})
	assert.Nil(t, err)

	request := testRequest{MandatoryField: "afield", ExtraField: "extra"}
	result := &testResult{}

	err = ogClient.Exec(nil, &request, result)
	apiErr, ok := err.(*ApiError)
	assert.True(t, ok)
	assert.Equal(t, apiErr.StatusCode, 422)
	assert.Contains(t, apiErr.Error(), "422")
	assert.Contains(t, apiErr.Error(), "Invalid recipient")
}

func TestExecWhenApiReturns5XX(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{
    "message": "Internal Server Error",
    "took": 0.083,
    "requestId": "6c20ec4e-076a-4422-8d65-7b8ca92067ab"
}`)
	}))
	defer ts.Close()

	ogClient, err := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		RetryCount:     1,
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})
	assert.Nil(t, err)
	setZeroBackoff(ogClient)

	request := testRequest{MandatoryField: "afield", ExtraField: "extra"}
	result := &testResult{}

	err = ogClient.Exec(nil, &request, result)
	apiErr, ok := err.(*ApiError)
	assert.True(t, ok)
	assert.Equal(t, apiErr.StatusCode, 500)
	assert.Contains(t, apiErr.Error(), "500")
	assert.Contains(t, apiErr.Error(), "Internal Server Error")
}

func TestExecWhenApiReturnsRateLimitingDetails(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-RateLimit-State", "THROTTLED")
		w.Header().Add("X-RateLimit-Reason", "ACCOUNT")
		w.Header().Add("X-RateLimit-Period-In-Sec", "60")
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintln(w, `{
    "message": "TooManyRequests",
    "took": 1,
    "requestId": "rId"
}`)
	}))
	defer ts.Close()

	ogClient, err := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		RetryCount:     1,
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})
	assert.Nil(t, err)
	setZeroBackoff(ogClient)

	request := &testRequest{MandatoryField: "afield", ExtraField: "extra"}
	result := &testResult{}

	err = ogClient.Exec(nil, request, result)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "rId")
	assert.Equal(t, "THROTTLED", result.ResultMetadata.RateLimitState)
	assert.Equal(t, "ACCOUNT", result.ResultMetadata.RateLimitReason)
	assert.Equal(t, "60", result.ResultMetadata.RateLimitPeriod)
}

func TestSubscription(t *testing.T) {
	subscriber := MetricSubscriber{
		Process: subscriberProcessImpl,
	}
	subscriber.Register(HTTP)
	subscriber.Register(SDK)
	subscriber.Register(API)

	subscriber2 := MetricSubscriber{}
	subscriber2.Register(HTTP)

	expectedSubsMap := map[string][]MetricSubscriber{
		string(HTTP): {subscriber, subscriber2},
		string(SDK):  {subscriber},
		string(API):  {subscriber},
	}

	assert.Equal(t, len(expectedSubsMap["http"]), len(metricPublisher.SubscriberMap["http"]))
	assert.Equal(t, len(expectedSubsMap["sdk"]), len(metricPublisher.SubscriberMap["sdk"]))
	assert.Equal(t, len(expectedSubsMap["api"]), len(metricPublisher.SubscriberMap["api"]))
}

func subscriberProcessImpl(metric Metric) interface{} {
	return metric
}

func TestHttpMetric(t *testing.T) {
	var httpMetric *HttpMetric
	subscriber := MetricSubscriber{
		Process: func(metric Metric) interface{} {
			httpMetric, _ = metric.(*HttpMetric)
			return httpMetric
		},
	}
	subscriber.Register(HTTP)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{
    "message": "success",
    "took": 1,
    "requestId": "rId"
}`)
	}))
	defer ts.Close()

	ogClient, err := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		RetryCount:     1,
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})
	assert.Nil(t, err)

	request := &testRequest{MandatoryField: "afield", ExtraField: "extra"}
	result := &testResult{}

	err = ogClient.Exec(nil, request, result)
	assert.Nil(t, err)

	expectedMetric := &HttpMetric{
		RetryCount:   0,
		Error:        nil,
		ResourcePath: "/an-enpoint",
		Status:       "200 OK",
		StatusCode:   200,
	}

	assert.Equal(t, expectedMetric.StatusCode, httpMetric.StatusCode)
	assert.Equal(t, expectedMetric.Status, httpMetric.Status)
	assert.Equal(t, expectedMetric.RetryCount, httpMetric.RetryCount)
	assert.Equal(t, expectedMetric.ResourcePath, httpMetric.ResourcePath)
	assert.Nil(t, httpMetric.Error)
}

func TestHttpMetricWhenRequestRetried(t *testing.T) {
	var httpMetric *HttpMetric
	subscriber := MetricSubscriber{
		Process: func(metric Metric) interface{} {
			httpMetric, _ = metric.(*HttpMetric)
			return httpMetric
		},
	}
	subscriber.Register(HTTP)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGatewayTimeout)
		fmt.Fprintln(w, `{
    "message": "something went wrong",
}`)
	}))
	defer ts.Close()

	ogClient, err := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		RetryCount:     1,
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})
	assert.Nil(t, err)
	setZeroBackoff(ogClient)

	request := &testRequest{MandatoryField: "afield", ExtraField: "extra"}
	result := &testResult{}

	err = ogClient.Exec(nil, request, result)

	expectedMetric := &HttpMetric{
		RetryCount:   2,
		Error:        err,
		ResourcePath: "/an-enpoint",
		Status:       "504 Gateway Timeout",
		StatusCode:   504,
	}

	assert.Equal(t, expectedMetric.StatusCode, httpMetric.StatusCode)
	assert.Equal(t, expectedMetric.Status, httpMetric.Status)
	assert.Equal(t, expectedMetric.RetryCount, httpMetric.RetryCount)
	assert.Equal(t, expectedMetric.ResourcePath, httpMetric.ResourcePath)
	assert.Nil(t, httpMetric.Error)
}

func TestApiMetric(t *testing.T) {
	var apiMetric *ApiMetric
	subscriber := MetricSubscriber{
		Process: func(metric Metric) interface{} {
			apiMetric, _ = metric.(*ApiMetric)
			return apiMetric
		},
	}
	subscriber.Register(API)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-RateLimit-State", "THROTTLED")
		w.Header().Add("X-RateLimit-Reason", "ACCOUNT")
		w.Header().Add("X-RateLimit-Period-In-Sec", "60")
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintln(w, `{
    "message": "TooManyRequests",
    "took": 1,
    "requestId": "rId"
}`)
	}))
	defer ts.Close()

	ogClient, err := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		RetryCount:     1,
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})
	assert.Nil(t, err)
	setZeroBackoff(ogClient)

	request := &testRequest{MandatoryField: "afield", ExtraField: "extra"}
	result := &testResult{}

	err = ogClient.Exec(nil, request, result)
	apiErr, ok := err.(*ApiError)
	assert.True(t, ok)
	assert.Equal(t, "TooManyRequests", apiErr.Message)

	expectedMetric := ApiMetric{
		ResourcePath: "/an-enpoint",
		ResultMetadata: ResultMetadata{
			RequestId:       "rId",
			ResponseTime:    1,
			RateLimitState:  "THROTTLED",
			RateLimitReason: "ACCOUNT",
			RateLimitPeriod: "60",
			RetryCount:      2,
		},
	}

	assert.Equal(t, expectedMetric.ResourcePath, apiMetric.ResourcePath)
	assert.Equal(t, expectedMetric.ResultMetadata, apiMetric.ResultMetadata)
}

func TestSdkMetricWhenRequestIsNotValid(t *testing.T) {
	var sdkMetric *SdkMetric
	subscriber := MetricSubscriber{
		Process: func(metric Metric) interface{} {
			sdkMetric, _ = metric.(*SdkMetric)
			return sdkMetric
		},
	}
	subscriber.Register(SDK)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintln(w, `{
    "message": "invalid request",
    "took": 1,
    "requestId": "rId"
}`)
	}))
	defer ts.Close()

	ogClient, _ := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		RetryCount:     1,
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})

	request := &testRequest{ExtraField: "extra"}
	result := &testResult{}

	ogClient.Exec(nil, request, result)
	expectedMetric := &SdkMetric{
		ErrorType:         "request-validation-error",
		ErrorMessage:      "mandatory field cannot be empty",
		ResourcePath:      "/an-enpoint",
		SdkRequestDetails: request,
		SdkResultDetails:  result,
	}

	assert.Equal(t, expectedMetric.ResourcePath, sdkMetric.ResourcePath)
	assert.Equal(t, expectedMetric.ErrorType, sdkMetric.ErrorType)
	assert.Equal(t, expectedMetric.ErrorMessage, sdkMetric.ErrorMessage)
	assert.Equal(t, expectedMetric.SdkRequestDetails, sdkMetric.SdkRequestDetails)
	assert.Equal(t, expectedMetric.SdkResultDetails, sdkMetric.SdkResultDetails)
}

func TestSdkMetricWhenExecSuccessful(t *testing.T) {
	var sdkMetric *SdkMetric
	subscriber := MetricSubscriber{
		Process: func(metric Metric) interface{} {
			sdkMetric, _ = metric.(*SdkMetric)
			return sdkMetric
		},
	}
	subscriber.Register(SDK)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{
    "message": "...",
    "took": 1,
    "requestId": "rId"
}`)
	}))
	defer ts.Close()

	ogClient, _ := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		RetryCount:     1,
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})

	request := &testRequest{MandatoryField: "f1", ExtraField: "extra"}
	result := &testResult{}

	err := ogClient.Exec(nil, request, result)
	assert.Nil(t, err)

	expectedMetric := &SdkMetric{
		ErrorType:         "",
		ErrorMessage:      "",
		ResourcePath:      "/an-enpoint",
		SdkRequestDetails: request,
		SdkResultDetails:  result,
	}

	assert.Equal(t, expectedMetric.ResourcePath, sdkMetric.ResourcePath)
	assert.Equal(t, expectedMetric.ErrorType, sdkMetric.ErrorType)
	assert.Equal(t, expectedMetric.ErrorMessage, sdkMetric.ErrorMessage)
	assert.Equal(t, expectedMetric.SdkRequestDetails, sdkMetric.SdkRequestDetails)
	assert.Equal(t, expectedMetric.SdkResultDetails, sdkMetric.SdkResultDetails)
}

func TestConfiguration(t *testing.T) {

	customHttpClient := http.DefaultClient
	customHttpClient.Timeout = time.Second * 1

	customLogger := logrus.New()
	customLogger.SetLevel(logrus.DebugLevel)
	customLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		PrettyPrint:     true,
	})

	retryFunc := func(ctx context.Context, resp *http.Response, err error) (b bool, e error) {
		return false, errors.New("testError")
	}

	backOff := func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		return time.Millisecond * 1500
	}

	conf := &Config{
		ApiKey:         "apiKey",
		OpsGenieAPIURL: API_URL_EU,
		HttpClient:     customHttpClient,
		RetryCount:     7,
		RetryPolicy:    retryFunc,
		Backoff:        backOff,
		Logger:         customLogger,
		LogLevel:       logrus.ErrorLevel,
	}

	ogClient, _ := NewOpsGenieClient(conf)

	apiRequest := &testRequest{MandatoryField: "f1", ExtraField: "extra"}

	assert.Equal(t, 7, ogClient.RetryableClient.RetryMax)
	assert.Equal(t, ogClient.Config.Logger, customLogger)
	assert.Equal(t, "https://api.eu.opsgenie.com/an-enpoint", buildRequestUrl(ogClient, apiRequest, nil))
	assert.Equal(t, ogClient.RetryableClient.HTTPClient, customHttpClient)

	flag, err := ogClient.RetryableClient.CheckRetry(nil, nil, nil)
	assert.False(t, flag)
	assert.NotNil(t, err)
	assert.Equal(t, "testError", err.Error())
	assert.Equal(t, time.Millisecond*1500, ogClient.RetryableClient.Backoff(0, 0, 0, nil))
}

func TestProxyConfiguration(t *testing.T) {

	var request *http.Request

	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request = r
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, `{
    		"message": "Success",
    		"took": 1,
    		"requestId": "32c9e6ba-e5b0-4dea-aa18-e3c3352a6d96"
		}`)
	}))

	psUrl, _ := url.Parse(proxyServer.URL)

	ogClient, err := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		RetryCount:     1,
		OpsGenieAPIURL: "--",
		ProxyConfiguration: &ProxyConfiguration{
			Username: "admin",
			Password: "1234",
			Host:     psUrl.Host,
			Protocol: Http,
		},
	})
	assert.Nil(t, err)

	apiRequest := &testRequest{MandatoryField: "f1", ExtraField: "extra"}
	result := &testResult{}

	err = ogClient.Exec(nil, apiRequest, result)
	assert.Nil(t, err)

	assert.Equal(t, "GenieKey apiKey", request.Header.Get("Authorization"))
	assert.NotEmpty(t, request.Header.Get("Proxy-Authorization"))

}

func TestProxyConfigurationWhenAuthorizationIsNotRequired(t *testing.T) {

	var request *http.Request

	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request = r
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{
			"message": "",
			"took": 1,
			"requestId": "1d3993a1-9058-46f0-bbc0-93f639ad2e27"
		}`)
	}))

	psUrl, _ := url.Parse(proxyServer.URL)

	ogClient, err := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		RetryCount:     1,
		OpsGenieAPIURL: "--",
		ProxyConfiguration: &ProxyConfiguration{
			Host:     psUrl.Host,
			Protocol: Http,
		},
	})
	assert.Nil(t, err)

	apiRequest := &testRequest{MandatoryField: "f1", ExtraField: "extra"}
	result := &testResult{}

	err = ogClient.Exec(nil, apiRequest, result)
	assert.Nil(t, err)

	assert.Equal(t, "GenieKey apiKey", request.Header.Get("Authorization"))
	assert.Empty(t, request.Header.Get("Proxy-Authorization"))
}

func TestRetrieveStatusWithRetrying(t *testing.T) {

	attemptCount := 0

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Opsgenie-Errortype", "RequestNotProcessed")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, []byte(`{
			"message": "request not processed",
			"took": 1,
			"requestId": "rId"
		}`))
		attemptCount++
	}))
	defer ts.Close()

	ogClient, err := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		RetryCount:     3,
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})
	assert.Nil(t, err)
	ogClient.RetryableClient.RetryWaitMax = time.Duration(0)

	request := &testRequest{MandatoryField: "afield", ExtraField: "extra"}
	result := &testResult{}

	asyncBaseResult := AsyncBaseResult{Client: ogClient}

	start := time.Now().UnixNano()
	err = asyncBaseResult.RetrieveStatus(nil, request, result)
	end := time.Now().UnixNano()

	delta := float64(100 * time.Millisecond.Nanoseconds())
	fmt.Println("start: ", start, "\nend  : ", end, "\ndiff : ", float64(end-start), "\ndelta: ", delta)
	assert.InDelta(t, end, start, delta)

	assert.Equal(t, ogClient.Config.RetryCount+1, attemptCount)

	apiErr, ok := err.(*ApiError)

	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	assert.Equal(t, "RequestNotProcessed", apiErr.ErrorHeader)
}

func TestRetrieveStatusWithoutErrorType(t *testing.T) {

	attemptCount := 0

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, []byte(`{
			"message": "request not found",
			"took": 1,
			"requestId": "rId"
		}`))
		attemptCount++
	}))
	defer ts.Close()

	ogClient, err := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		RetryCount:     4,
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})
	assert.Nil(t, err)
	setZeroBackoff(ogClient)

	request := &testRequest{MandatoryField: "afield", ExtraField: "extra"}
	result := &testResult{}

	asyncBaseResult := AsyncBaseResult{Client: ogClient}

	err = asyncBaseResult.RetrieveStatus(nil, request, result)

	assert.Equal(t, 1, attemptCount)

	apiError, ok := err.(*ApiError)

	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiError.StatusCode)
}

func TestRetrieveStatus(t *testing.T) {

	attemptCount := 0

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{
			"Data": "Success",
			"took": 1,
			"requestId": "rId"
		}`)
		attemptCount++
	}))
	defer ts.Close()

	ogClient, err := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		RetryCount:     4,
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})
	assert.Nil(t, err)
	setZeroBackoff(ogClient)

	request := &testRequest{MandatoryField: "afield", ExtraField: "extra"}
	result := &testResult{}

	asyncBaseResult := AsyncBaseResult{Client: ogClient}

	err = asyncBaseResult.RetrieveStatus(nil, request, result)

	assert.Equal(t, 1, attemptCount)
	assert.Nil(t, err)
	assert.Equal(t, "Success", result.Data)
}

func TestRetrieveStatusContextDeadline(t *testing.T) {

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Millisecond*50)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 250)
	}))

	request := &testRequest{MandatoryField: "afield", ExtraField: "extra"}
	result := &testResult{}

	ogClient, err := NewOpsGenieClient(&Config{
		ApiKey:         "apiKey",
		RetryCount:     4,
		OpsGenieAPIURL: ApiUrl(strings.TrimPrefix(ts.URL, "http://")),
	})
	assert.Nil(t, err)

	asyncBaseResult := AsyncBaseResult{Client: ogClient}

	start := time.Now().UnixNano()
	err = asyncBaseResult.RetrieveStatus(ctx, request, result)
	end := time.Now().UnixNano()

	assert.EqualError(t, err, "context deadline exceeded")
	delta := float64(100 * time.Millisecond.Nanoseconds())
	fmt.Println("start: ", start, "\nend  : ", end, "\ndiff : ", float64(end-start), "\ndelta: ", delta)
	assert.InDelta(t, end, start, delta)
}

func setZeroBackoff(client *OpsGenieClient) {
	client.RetryableClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		return time.Duration(0)
	}
}
