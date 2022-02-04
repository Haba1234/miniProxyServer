package proxyserver

import (
	"encoding/json"
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type Data struct {
	ID      uuid.UUID   `json:"id"`
	Status  string      `json:"status"`
	Headers http.Header `json:"headers"`
	Length  int64       `json:"length"`
}

type Body struct {
	Data  `json:"data"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type Response struct {
	writer http.ResponseWriter
	header http.Header
	Body
}

func NewResponse(w http.ResponseWriter) (r *Response) {
	return &Response{writer: w,
		header: w.Header()}
}

func (r *Response) setStatusCode(code int) {
	r.writer.WriteHeader(code)
}

func (r *Response) write(statusCode int, errMessage string) {
	r.header.Set("Content-Type", "application/json")

	r.Error.Message = errMessage
	respBuf, err := json.Marshal(r.Body)
	if err != nil {
		log.Printf("body marshal error: %s", err)
		r.setStatusCode(http.StatusInternalServerError)
	}

	r.setStatusCode(statusCode)
	_, err = r.writer.Write(respBuf)
	if err != nil {
		log.Printf("write error: %s", err)
	}
}

func (r *Response) toSend(id uuid.UUID, status string, headers http.Header, length int64) {
	r.ID = id
	r.Status = status
	r.Body.Headers = headers
	r.Length = length
}
