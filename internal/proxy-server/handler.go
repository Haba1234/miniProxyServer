package proxyserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Request struct {
	Method  string            `json:"method,omitempty"`
	URL     string            `json:"url,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

type QueryResult struct {
	Status  string
	Headers http.Header
	Length  int64
}

func (s *Service) RedirectReq(w http.ResponseWriter, r *http.Request) {
	resp := NewResponse(w)
	if r.Method != http.MethodPost {
		errMessage := fmt.Sprintf("method %s is not supported for uri %s\n", r.Method, r.URL.Path)
		resp.write(http.StatusMethodNotAllowed, errMessage)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.write(http.StatusBadRequest, err.Error())
		return
	}

	req := Request{}
	err = json.Unmarshal(buf, &req)
	if err != nil {
		resp.write(http.StatusBadRequest, err.Error())
		return
	}

	if req.Method != http.MethodGet {
		resp.write(http.StatusBadRequest, "method wrong")
		return
	}

	if _, err = url.ParseRequestURI(req.URL); err != nil {
		resp.write(http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("new redirect Request (method=%s, url=%s headers=%v)",
		req.Method, req.URL, req.Headers)

	gr, err := getRequest(&req, s.Interval)
	if err != nil {
		resp.write(http.StatusInternalServerError, err.Error())
	}

	uid := uuid.NewV4()
	resp.toSend(uid, gr.Status, gr.Headers, gr.Length)
	s.toSave(uid, req, gr)
	resp.write(http.StatusOK, "")
}

func getRequest(data *Request, interval time.Duration) (QueryResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), interval)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, data.Method, data.URL, nil)
	if err != nil {
		return QueryResult{}, err
	}

	for key, value := range data.Headers {
		req.Header.Add(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return QueryResult{}, err
	}
	defer resp.Body.Close()

	return QueryResult{
		Status:  resp.Status,
		Headers: resp.Header,
		Length:  resp.ContentLength,
	}, nil
}
