package client

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jaegertracing/jaeger/plugin/storage/haystack/store/config"
	"go.uber.org/zap"
)

var (
	errReceivedNon204StatusCode = errors.New("received non 204 success code")
)

// HTTPClient struct defines http.Client
type HTTPClient struct {
	client    *http.Client
	authToken string
	endpoint  string
	logger    *zap.Logger
}

// NewHTTPClient constructor
func NewHTTPClient(config config.HaystackConfig, logger *zap.Logger) *HTTPClient {
	var defaultTransport http.RoundTripper = &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		MaxIdleConns:        config.HTTPMaxIdleConns,
		MaxIdleConnsPerHost: config.HTTPMaxIdleConnsPerHost,
		IdleConnTimeout:     time.Duration(config.HTTPRequestTimeout) * time.Second,
	}

	client := &http.Client{Transport: defaultTransport}

	return &HTTPClient{
		client:    client,
		authToken: config.AuthToken,
		endpoint:  config.ProxyURL,
		logger:    logger,
	}
}

// SetEndpoint sets the endpoint
func (c *HTTPClient) SetEndpoint(endpoint string) {
	c.endpoint = endpoint
}

// Post sends request to endpoint
func (c *HTTPClient) Post(batch []byte) error {
	req, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewBuffer(batch))
	if err != nil {
		c.logger.Error("failed to create new batch request", zap.String("error", err.Error()))
		return err
	}
	// set headers
	req.Header.Set("x-auth-token", c.authToken)
	req.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("failed to send request", zap.String("error", err.Error()))
		return err
	}
	if response == nil {
		c.logger.Error("failed to receive response object")
		return nil
	}

	if response.StatusCode != http.StatusNoContent {
		var responseMsg string
		if response.Body != nil {
			resp, err := ioutil.ReadAll(response.Body)
			if err == nil {
				responseMsg = string(resp)
			} else {
				c.logger.Warn("Failed to read response body", zap.Int("statusCode", response.StatusCode), zap.String("error", err.Error()))
			}
		}
		c.logger.Warn("Received non 204 response status code", zap.Int("statusCode", response.StatusCode), zap.String("response", responseMsg))
		return errReceivedNon204StatusCode
	}
	return nil
}
