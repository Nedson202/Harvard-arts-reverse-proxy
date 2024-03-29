// Licensed to Elasticsearch B.V under one or more agreements.
// Elasticsearch B.V. licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.
//
// Code generated from specification version 8.0.0: DO NOT EDIT

package esapi

import (
	"context"
	"net/http"
	"strings"
)

func newEnrichGetPolicyFunc(t Transport) EnrichGetPolicy {
	return func(o ...func(*EnrichGetPolicyRequest)) (*Response, error) {
		var r = EnrichGetPolicyRequest{}
		for _, f := range o {
			f(&r)
		}
		return r.Do(r.ctx, t)
	}
}

// ----- API Definition -------------------------------------------------------

// EnrichGetPolicy -
//
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/current/enrich-get-policy.html.
//
type EnrichGetPolicy func(o ...func(*EnrichGetPolicyRequest)) (*Response, error)

// EnrichGetPolicyRequest configures the Enrich Get Policy API request.
//
type EnrichGetPolicyRequest struct {
	Name string

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}

// Do executes the request and returns response or error.
//
func (r EnrichGetPolicyRequest) Do(ctx context.Context, transport Transport) (*Response, error) {
	var (
		method string
		path   strings.Builder
		params map[string]string
	)

	method = "GET"

	path.Grow(1 + len("_enrich") + 1 + len("policy") + 1 + len(r.Name))
	path.WriteString("/")
	path.WriteString("_enrich")
	path.WriteString("/")
	path.WriteString("policy")
	if r.Name != "" {
		path.WriteString("/")
		path.WriteString(r.Name)
	}

	params = make(map[string]string)

	if r.Pretty {
		params["pretty"] = "true"
	}

	if r.Human {
		params["human"] = "true"
	}

	if r.ErrorTrace {
		params["error_trace"] = "true"
	}

	if len(r.FilterPath) > 0 {
		params["filter_path"] = strings.Join(r.FilterPath, ",")
	}

	req, err := newRequest(method, path.String(), nil)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if len(r.Header) > 0 {
		if len(req.Header) == 0 {
			req.Header = r.Header
		} else {
			for k, vv := range r.Header {
				for _, v := range vv {
					req.Header.Add(k, v)
				}
			}
		}
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	res, err := transport.Perform(req)
	if err != nil {
		return nil, err
	}

	response := Response{
		StatusCode: res.StatusCode,
		Body:       res.Body,
		Header:     res.Header,
	}

	return &response, nil
}

// WithContext sets the request context.
//
func (f EnrichGetPolicy) WithContext(v context.Context) func(*EnrichGetPolicyRequest) {
	return func(r *EnrichGetPolicyRequest) {
		r.ctx = v
	}
}

// WithName - the name of the enrich policy.
//
func (f EnrichGetPolicy) WithName(v string) func(*EnrichGetPolicyRequest) {
	return func(r *EnrichGetPolicyRequest) {
		r.Name = v
	}
}

// WithPretty makes the response body pretty-printed.
//
func (f EnrichGetPolicy) WithPretty() func(*EnrichGetPolicyRequest) {
	return func(r *EnrichGetPolicyRequest) {
		r.Pretty = true
	}
}

// WithHuman makes statistical values human-readable.
//
func (f EnrichGetPolicy) WithHuman() func(*EnrichGetPolicyRequest) {
	return func(r *EnrichGetPolicyRequest) {
		r.Human = true
	}
}

// WithErrorTrace includes the stack trace for errors in the response body.
//
func (f EnrichGetPolicy) WithErrorTrace() func(*EnrichGetPolicyRequest) {
	return func(r *EnrichGetPolicyRequest) {
		r.ErrorTrace = true
	}
}

// WithFilterPath filters the properties of the response body.
//
func (f EnrichGetPolicy) WithFilterPath(v ...string) func(*EnrichGetPolicyRequest) {
	return func(r *EnrichGetPolicyRequest) {
		r.FilterPath = v
	}
}

// WithHeader adds the headers to the HTTP request.
//
func (f EnrichGetPolicy) WithHeader(h map[string]string) func(*EnrichGetPolicyRequest) {
	return func(r *EnrichGetPolicyRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		for k, v := range h {
			r.Header.Add(k, v)
		}
	}
}
