package check

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	MAX_TIMEOUT = 5 * time.Minute

	DEFAULT_REQUEST_TIMEOUT = 5 * time.Second
	DEFAULT_INTERVAL        = 1 * time.Second
)

type HttpGetReadyCheck struct {
	client         *http.Client  `json:"-"`
	Url            *url.URL      `json:"url"`
	StatusCode     int           `json:"status"`
	Timeout        time.Duration `json:"timeout"`
	RequestTimeout time.Duration `json:"request_timeout"`
	Interval       time.Duration `json:"interval"`
}

func (h *HttpGetReadyCheck) UnmarshalJSON(data []byte) error {
	type Alias HttpGetReadyCheck
	aux := &struct {
		*Alias
		Url string `json:"url"`
	}{Alias: (*Alias)(h)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	h.Url, err = url.Parse(aux.Url)
	if err != nil {
		return err
	}

	if h.Interval == 0 {
		logrus.Debugf("Interval not set, using default %s", DEFAULT_INTERVAL)
		h.Interval = DEFAULT_INTERVAL
	}
	if h.Timeout == 0 {
		logrus.Debugf("Timeout not set, using default %s", MAX_TIMEOUT)
		h.Timeout = MAX_TIMEOUT
	}
	if h.Timeout > MAX_TIMEOUT {
		logrus.Warnf("Timeout %s is greater than max timeout %s, using max timeout", h.Timeout, MAX_TIMEOUT)
		h.Timeout = MAX_TIMEOUT
	}
	if h.RequestTimeout == 0 {
		logrus.Debugf("RequestTimeout not set, using default %s", DEFAULT_REQUEST_TIMEOUT)
		h.RequestTimeout = DEFAULT_REQUEST_TIMEOUT
	}

	h.client = &http.Client{
		Timeout: h.RequestTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return nil
}

func (h *HttpGetReadyCheck) RenderUrl(context interface{}) (string, error) {
	tmpl, err := template.New("url").Parse(h.Url.String())
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, context)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Validate ensures the setup is appropriate for the test context
// In this case, httpGet has no dependencies, so it's always valid
func (h *HttpGetReadyCheck) Validate() error {
	return nil
}

func (h *HttpGetReadyCheck) WaitForReady(context interface{}) error {
	renderedUrl, err := h.RenderUrl(context)
	if err != nil {
		return err
	}

	endTime := time.Now().Add(h.Timeout)
	for time.Now().Before(endTime) {
		resp, err := h.client.Get(renderedUrl)
		if err != nil {
			return err
		}
		if resp.StatusCode == h.StatusCode {
			return nil
		}
		time.Sleep(h.Interval)
	}
	return fmt.Errorf("timed out waiting for %s to return %d", h.Url.String(), h.StatusCode)
}
