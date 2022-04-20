package opsgenie

import (
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
)

type OpsgenieClient struct {
	client *client.OpsGenieClient
}

type Config struct {
	ApiKey          string
	ApiUrl          string
	ApiRetryCount   int
	ApiRetryWaitMin int
	ApiRetryWaitMax int
}

func (c *Config) Client() (*OpsgenieClient, error) {
	config := &client.Config{
		ApiKey:         c.ApiKey,
		RetryCount:     c.ApiRetryCount,
		OpsGenieAPIURL: client.ApiUrl(c.ApiUrl),
		Backoff: func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
			if c.ApiRetryWaitMin > 0 && c.ApiRetryWaitMax > 0 {
				return retryablehttp.DefaultBackoff(time.Duration(c.ApiRetryWaitMin)*time.Second, time.Duration(c.ApiRetryWaitMax)*time.Second, attemptNum, resp)
			} else {
				return retryablehttp.DefaultBackoff(min, max, attemptNum, resp)
			}
		},
	}
	ogCli, err := client.NewOpsGenieClient(config)
	if err != nil {
		return nil, err
	}
	ogClient := OpsgenieClient{}
	ogClient.client = ogCli
	log.Printf("[INFO] OpsGenie client configured")
	return &ogClient, nil
}
