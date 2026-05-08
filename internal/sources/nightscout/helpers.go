package nightscout

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
	"zgluco/internal/types"
)

func (n *Nightscout) doApiCall(endPoint, findCriteria, sortCriteria string, timeRange types.TimeRange) ([]byte, error) {
	var all []json.RawMessage

	for t := timeRange.From; t.Before(timeRange.To); t = t.Add(24 * time.Hour) {
		end := t.Add(24 * time.Hour)
		upperOp := "$lt"
		if end.After(timeRange.To) {
			end = timeRange.To
			upperOp = "$lte"
		}

		chunk, err := n.doSinglePage(endPoint, findCriteria, sortCriteria, upperOp, types.TimeRange{From: t, To: end})
		if err != nil {
			return nil, err
		}

		var items []json.RawMessage
		if err = json.Unmarshal(chunk, &items); err != nil {
			return nil, err
		}

		all = append(all, items...)
	}

	return json.Marshal(all)
}

func (n *Nightscout) doSinglePage(endPoint, findCriteria, sortCriteria, upperOp string, timeRange types.TimeRange) ([]byte, error) {
	var params = make(url.Values)

	params.Add(fmt.Sprintf("find[%s][$gte]", findCriteria), timeRange.From.Format(time.RFC3339Nano))
	params.Add(fmt.Sprintf("find[%s][%s]", findCriteria, upperOp), timeRange.To.Format(time.RFC3339Nano))

	if sortCriteria != "" {
		params.Add(fmt.Sprintf("sort[%s]", sortCriteria), "1")
	}

	params.Add("count", "1000")

	apiUrl := fmt.Sprintf("%s%s?%s", n.nightscoutUrl, endPoint, params.Encode())

	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	if n.apiKey != "" {
		req.Header.Add("api-secret", fmt.Sprintf("%x", sha1.Sum([]byte(n.apiKey))))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		var nsErr struct {
			Message     string `json:"message"`
			Description string `json:"description"`
		}
		if json.Unmarshal(b, &nsErr) == nil && nsErr.Message != "" {
			if nsErr.Description != "" {
				return nil, fmt.Errorf("%s (%d): %s", nsErr.Message, res.StatusCode, nsErr.Description)
			}
			return nil, fmt.Errorf("%s (%d)", nsErr.Message, res.StatusCode)
		}
		return nil, fmt.Errorf("HTTP %d: %s", res.StatusCode, b)
	}

	return b, nil
}
