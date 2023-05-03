package router

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func StartRecord(req *http.Request, start time.Time) *http.Request {
	ctx := req.Context()

	v := new(Data)
	v.RequestID = uuid.New().String()

	v.Host = req.Host
	v.Endpoint = req.URL.Path
	v.TimeStart = start
	v.Device = "Web-Base"

	v.RequestMethod = req.Method
	v.RequestHeader = DumpRequest(req)

	ctx = context.WithValue(ctx, LogKey, v)

	return req.WithContext(ctx)
}
