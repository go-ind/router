package goindrouter

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
)

func ResponseJSON(w http.ResponseWriter, ctx context.Context, code int, status bool, message string, rs, pagination interface{}) {
	resservice := Responseservice{}
	resservice.Status = code
	if status {
		resservice.Data = rs
		resservice.Pagination = pagination
	} else {
		resservice.ErrorMessage = "Error"
	}

	var input []byte

	resservice.Message = message
	switch rs.(type) {
	case string:
		input = []byte(rs.(string))
	case []byte:
		input = rs.([]byte)
	default:
		input, _ = JSONMarshal(rs)
	}
	if ctx == nil {
		// Handle For CTX if is null
		ctx = context.TODO()
	}
	Logger(ctx, string(input), code)
	origin := "*"

	v, ok := ctx.Value(LogKey).(*Data)
	if ok {
		words := strings.Fields(v.RequestHeader)
		for i := 0; i < len(words); i++ {
			if words[i] == "Origin:" {
				origin = words[i+1]
				break
			}
		}

	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Strict-Transport-Security", "max-age=15552000; includeSubDomains")
	w.Header().Set("X-DNS-Prefetch-Control", "off")
	w.Header().Set("Vary", "X-HTTP-Method-Override")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resservice)
}

func ResponseXML(w http.ResponseWriter, ctx context.Context, code int, status bool, message string, rs, pagination interface{}) {
	resservice := Responseservice{}
	resservice.Status = code
	if status {
		resservice.Data = rs
		resservice.Pagination = pagination
	} else {
		resservice.ErrorMessage = "Error"
	}

	var input []byte

	resservice.Message = message
	switch rs.(type) {
	case string:
		input = []byte(rs.(string))
	case []byte:
		input = rs.([]byte)
	default:
		input, _ = JSONMarshal(rs)
	}
	if ctx == nil {
		// Handle For CTX if is null
		ctx = context.TODO()
	}
	Logger(ctx, string(input), code)
	origin := "*"

	v, ok := ctx.Value(LogKey).(*Data)
	if ok {
		words := strings.Fields(v.RequestHeader)
		for i := 0; i < len(words); i++ {
			if words[i] == "Origin:" {
				origin = words[i+1]
				break
			}
		}

	}
	x, err := xml.MarshalIndent(resservice, "", "	")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// json.NewEncoder(w).Encode(resservice)

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Strict-Transport-Security", "max-age=15552000; includeSubDomains")
	w.Header().Set("X-DNS-Prefetch-Control", "off")
	w.Header().Set("Vary", "X-HTTP-Method-Override")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	w.WriteHeader(code)
	w.Write(x)
}
