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

func newSecurityGetBuiltinPrivilegesFunc(t Transport) SecurityGetBuiltinPrivileges {
	return func(o ...func(*SecurityGetBuiltinPrivilegesRequest)) (*Response, error) {
		var r = SecurityGetBuiltinPrivilegesRequest{}
		for _, f := range o {
			f(&r)
		}
		return r.Do(r.ctx, t)
	}
}

// ----- API Definition -------------------------------------------------------

// SecurityGetBuiltinPrivileges -
//
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/current/security-api-get-builtin-privileges.html.
//
type SecurityGetBuiltinPrivileges func(o ...func(*SecurityGetBuiltinPrivilegesRequest)) (*Response, error)

// SecurityGetBuiltinPrivilegesRequest configures the Security Get Builtin Privileges API request.
//
type SecurityGetBuiltinPrivilegesRequest struct {
	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}

// Do executes the request and returns response or error.
//
func (r SecurityGetBuiltinPrivilegesRequest) Do(ctx context.Context, transport Transport) (*Response, error) {
	var (
		method string
		path   strings.Builder
		params map[string]string
	)

	method = "GET"

	path.Grow(len("/_security/privilege/_builtin"))
	path.WriteString("/_security/privilege/_builtin")

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
func (f SecurityGetBuiltinPrivileges) WithContext(v context.Context) func(*SecurityGetBuiltinPrivilegesRequest) {
	return func(r *SecurityGetBuiltinPrivilegesRequest) {
		r.ctx = v
	}
}

// WithPretty makes the response body pretty-printed.
//
func (f SecurityGetBuiltinPrivileges) WithPretty() func(*SecurityGetBuiltinPrivilegesRequest) {
	return func(r *SecurityGetBuiltinPrivilegesRequest) {
		r.Pretty = true
	}
}

// WithHuman makes statistical values human-readable.
//
func (f SecurityGetBuiltinPrivileges) WithHuman() func(*SecurityGetBuiltinPrivilegesRequest) {
	return func(r *SecurityGetBuiltinPrivilegesRequest) {
		r.Human = true
	}
}

// WithErrorTrace includes the stack trace for errors in the response body.
//
func (f SecurityGetBuiltinPrivileges) WithErrorTrace() func(*SecurityGetBuiltinPrivilegesRequest) {
	return func(r *SecurityGetBuiltinPrivilegesRequest) {
		r.ErrorTrace = true
	}
}

// WithFilterPath filters the properties of the response body.
//
func (f SecurityGetBuiltinPrivileges) WithFilterPath(v ...string) func(*SecurityGetBuiltinPrivilegesRequest) {
	return func(r *SecurityGetBuiltinPrivilegesRequest) {
		r.FilterPath = v
	}
}

// WithHeader adds the headers to the HTTP request.
//
func (f SecurityGetBuiltinPrivileges) WithHeader(h map[string]string) func(*SecurityGetBuiltinPrivilegesRequest) {
	return func(r *SecurityGetBuiltinPrivilegesRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		for k, v := range h {
			r.Header.Add(k, v)
		}
	}
}
