package main

import "time"

const HttpMethodGet = "GET"
const HttpMethodPost = "POST"
const HttpMethodPut = "PUT"
const HttpMethodPatch = "PATCH"
const HttpMethodDelete = "DELETE"

type (
	// Key context
	Key int
)

const (
	services = "Internal Service"

	LogKey = Key(48)
)

// Data is data standard output
type Data struct {
	RequestID     string    `json:"RequestID"`
	TimeStart     time.Time `json:"TimeStart"`
	UserCode      string    `json:"UserCode"`
	Device        string    `json:"Device"`
	Host          string    `json:"Host"`
	Endpoint      string    `json:"Endpoint"`
	RequestMethod string    `json:"RequestMethod"`
	RequestHeader string    `json:"RequestHeader"`
	StatusCode    int       `json:"StatusCode"`
	Response      string    `json:"Response"`
	ExecTime      float64   `json:"ExecutionTime"`
	Messages      []string  `json:"Messages"`
}

type Responseservice struct {
	Status       int         `json:"status"`
	ErrorMessage string      `json:"error_message"`
	Data         interface{} `json:"data"`
	Pagination   interface{} `json:"pagination"`
	Message      string      `json:"message"`
}
