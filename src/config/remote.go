package config

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type FetchConfig struct {
	URL        string
	AuthHeader string
	Insecure   bool
}

func FetchRemote(cfg FetchConfig) ([]byte, error) {
	if cfg.URL == "" {
		return nil, errors.New("empty url")
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.Insecure,
			},
		},
	}

	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		req, err := http.NewRequest("GET", cfg.URL, nil)
		if err != nil {
			return nil, err
		}

		if cfg.AuthHeader != "" {
			parts := strings.SplitN(cfg.AuthHeader, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				req.Header.Set(key, value)
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			if attempt < 3 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = errors.New("bad status")
			if attempt < 3 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			if attempt < 3 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		return data, nil
	}

	return nil, lastErr
}