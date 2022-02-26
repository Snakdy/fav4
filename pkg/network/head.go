package network

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var LineMatcherV2 = regexp.MustCompile(`href="?[^"]*favicon[^"]*(png|ico)`)

var (
	ErrNotFound      = errors.New("failed to find a favicon reference in the page source")
	ErrRequestFailed = errors.New("request returned an invalid response code")
)

type HeadLoader struct {
	client *http.Client
}

func NewHeadLoader(client *http.Client) *HeadLoader {
	return &HeadLoader{
		client: client,
	}
}

func (l *HeadLoader) Get(ctx context.Context, target *url.URL) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target.String(), nil)
	if err != nil {
		log.WithError(err).WithContext(ctx).Error("failed to prepare request")
		return "", err
	}
	resp, err := l.client.Do(req)
	if err != nil {
		log.WithError(err).WithContext(ctx).Error("failed to execute request")
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.WithContext(ctx).Errorf("request failed with code: %d", resp.StatusCode)
		return "", ErrRequestFailed
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, "favicon.") {
			val := strings.TrimPrefix(strings.ReplaceAll(LineMatcherV2.FindString(line), `"`, ""), "href=")
			if !strings.HasPrefix(val, "/") {
				return val, nil
			}
			return fmt.Sprintf("https://%s%s", target.Host, val), nil
		}
	}
	return "", ErrNotFound
}
