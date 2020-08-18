package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jaegertracing/jaeger/plugin/storage/haystack/store/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewHttpClient(t *testing.T) {
	var (
		v   = viper.New()
		cfg = config.HaystackConfig{}
	)
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err)
	cfg.InitFromViper(v)
	httpClient := NewHTTPClient(cfg, logger)
	assert.NotEmpty(t, httpClient)
	assert.NotEmpty(t, httpClient.client)
	assert.Equal(t, httpClient.authToken, "")
	assert.Equal(t, httpClient.endpoint, "")
}

func TestPost(t *testing.T) {
	var (
		v   = viper.New()
		cfg = config.HaystackConfig{}
	)
	logger, e := zap.NewDevelopment()
	assert.Nil(t, e)
	cfg.InitFromViper(v)
	httpClient := NewHTTPClient(cfg, logger)

	t.Log("Test endpoint empty error case")
	e = httpClient.Post(nil)
	assert.NotNil(t, e)

	t.Log("Test endpoint returns non success response status code")
	failureJobHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	server := startMockESCluster(failureJobHandler)
	httpClient.endpoint = server.URL
	e = httpClient.Post(nil)
	assert.NotNil(t, e)

	t.Log("Test endpoint returns success response status code")
	successJobHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
	server = startMockESCluster(successJobHandler)
	httpClient.endpoint = server.URL
	e = httpClient.Post(nil)
	assert.Nil(t, e)
}

func startMockESCluster(jobHandler func(_ http.ResponseWriter, r *http.Request)) *httptest.Server {
	server := httptest.NewUnstartedServer(http.HandlerFunc(jobHandler))
	server.Start()
	return server
}
