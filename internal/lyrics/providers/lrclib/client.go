package lrclib

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	errs "lrcsnc/internal/lyrics/errors"
	"lrcsnc/internal/pkg/log"
)

var httpClient = http.Client{
	Timeout: 10 * time.Second,
}

func requestLyrics(title string, artist string) ([]byte, error) {
	urlReqPath := "https://lrclib.net/api/search?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v", title, artist))
	_, err := url.Parse(urlReqPath)
	if err != nil {
		log.Fatal("lyrics/providers/lrclib/client", fmt.Sprintf("Failed to parse string (%v) to URL; please, report this issue to GitHub. More:\n%v", urlReqPath, err))
	}

	resp, err := httpClient.Get(urlReqPath)
	if os.IsTimeout(err) {
		return nil, errs.ServerTimeout
	}
	if err != nil || resp.StatusCode != 200 {
		return nil, errs.ServerError
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, errs.BodyReadFail
	}
	return body, nil
}
