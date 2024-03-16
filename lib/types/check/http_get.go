package check

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	Url            string        `json:"url"`
	StatusCode     int           `json:"status"`
	Timeout        time.Duration `json:"timeout"`
	RequestTimeout time.Duration `json:"requestTimeout"`
	Interval       time.Duration `json:"interval"`
}

func (h *HttpGetReadyCheck) UnmarshalJSON(data []byte) error {
	type auxH struct {
		Url            string `json:"url"`
		StatusCode     int    `json:"status"`
		Timeout        string `json:"timeout"`
		Interval       string `json:"interval"`
		RequestTimeout string `json:"requestTimeout"`
	}
	var aux auxH
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	h.Url = aux.Url
	h.StatusCode = aux.StatusCode
	if aux.Interval == "" {
		aux.Interval = DEFAULT_INTERVAL.String()
	}
	h.Interval, err = time.ParseDuration(aux.Interval)
	if err != nil {
		return err
	}
	if aux.Timeout == "" {
		aux.Timeout = MAX_TIMEOUT.String()
	}
	h.Timeout, err = time.ParseDuration(aux.Timeout)
	if err != nil {
		return err
	}
	if aux.RequestTimeout == "" {
		aux.RequestTimeout = DEFAULT_REQUEST_TIMEOUT.String()
	}
	h.RequestTimeout, err = time.ParseDuration(aux.RequestTimeout)
	if err != nil {
		return err
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
	tmpl, err := template.New("url").Parse(h.Url)
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
func (h *HttpGetReadyCheck) Validate(context interface{}, logFile io.Writer) error {
	return nil
}

func (h *HttpGetReadyCheck) WaitForReady(context interface{}, logFile io.Writer) error {
	renderedUrl, err := h.RenderUrl(context)
	if err != nil {
		return err
	}

	endTime := time.Now().Add(h.Timeout)
	logrus.Debugf("Waiting for %s to return %d\n", renderedUrl, h.StatusCode)
	for time.Now().Before(endTime) {
		resp, err := h.client.Get(renderedUrl)
		if err != nil {
			logrus.WithError(err).Debugf("Error getting %s", renderedUrl)
		} else if resp.StatusCode == h.StatusCode {
			logrus.Infof("Got %s with status %d", renderedUrl, h.StatusCode)
			return nil
		}
		time.Sleep(h.Interval)
	}
	return fmt.Errorf("timed out waiting for %s to return %d", renderedUrl, h.StatusCode)
}
