// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/cloustone/pandas"
	"github.com/cloustone/pandas/bootstrap"
	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/pkg/errors"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-zoo/bone"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	contentType  = "application/json"
	maxLimit     = 100
	defaultLimit = 10
)

var (
	errUnsupportedContentType = errors.New("unsupported content type")
	errInvalidQueryParams     = errors.New("invalid query params")
	errInvalidLimitParam      = errors.New("invalid limit query param")
	errInvalidOffsetParam     = errors.New("invalid offset query param")
	fullMatch                 = []string{"state", "external_id", "mainflux_id", "mainflux_key"}
	partialMatch              = []string{"name"}
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc bootstrap.Service, reader bootstrap.ConfigReader) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}
	r := bone.New()

	r.Post("/things/configs", kithttp.NewServer(
		addEndpoint(svc),
		decodeAddRequest,
		encodeResponse,
		opts...))

	r.Get("/things/configs/:id", kithttp.NewServer(
		viewEndpoint(svc),
		decodeEntityRequest,
		encodeResponse,
		opts...))

	r.Put("/things/configs/:id", kithttp.NewServer(
		updateEndpoint(svc),
		decodeUpdateRequest,
		encodeResponse,
		opts...))

	r.Patch("/things/configs/certs/:id", kithttp.NewServer(
		updateCertEndpoint(svc),
		decodeUpdateCertRequest,
		encodeResponse,
		opts...))

	r.Put("/things/configs/connections/:id", kithttp.NewServer(
		updateConnEndpoint(svc),
		decodeUpdateConnRequest,
		encodeResponse,
		opts...))

	r.Get("/things/configs", kithttp.NewServer(
		listEndpoint(svc),
		decodeListRequest,
		encodeResponse,
		opts...))

	r.Get("/things/unknown/configs", kithttp.NewServer(
		listEndpoint(svc),
		decodeUnknownRequest,
		encodeResponse,
		opts...))

	r.Get("/things/bootstrap/:external_id", kithttp.NewServer(
		bootstrapEndpoint(svc, reader, false),
		decodeBootstrapRequest,
		encodeResponse,
		opts...))

	r.Get("/things/bootstrap/secure/:external_id", kithttp.NewServer(
		bootstrapEndpoint(svc, reader, true),
		decodeBootstrapRequest,
		encodeSecureRes,
		opts...))

	r.Put("/things/state/:id", kithttp.NewServer(
		stateEndpoint(svc),
		decodeStateRequest,
		encodeResponse,
		opts...))

	r.Delete("/things/configs/:id", kithttp.NewServer(
		removeEndpoint(svc),
		decodeEntityRequest,
		encodeResponse,
		opts...))

	r.GetFunc("/version", pandas.Version("bootstrap"))
	r.Handle("/metrics", promhttp.Handler())

	return SetupMiddleware(r)
}

func assetFS() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:    bootstrap.Asset,
		AssetDir: bootstrap.AssetDir,
		Prefix:   "dist",
	}
}

//SSwagger s_swagger
func SSwagger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/swagger" || r.URL.Path == "/" {
			http.Redirect(w, r, "/swagger/", http.StatusFound)
			return
		}

		if strings.Index(r.URL.Path, "/swagger/") == 0 {
			http.StripPrefix("/swagger/", http.FileServer(assetFS())).ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

//RedocUI docs to show redoc ui
func RedocUI(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opts := middleware.RedocOpts{
			Path:     "docs",
			SpecURL:  r.URL.Host + "/swagger/static/swagger/swagger.yaml",
			RedocURL: r.URL.Host + "/swagger/static/js/redoc.standalone.js",
			Title:    "swagger api",
		}

		middleware.Redoc(opts, handler).ServeHTTP(w, r)
		return
	})
}

//SetupMiddleware setupmiddleware
func SetupMiddleware(handler http.Handler) http.Handler {
	return SSwagger(
		RedocUI(handler),
	)
}
func decodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := addReq{key: r.Header.Get("Authorization")}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(bootstrap.ErrMalformedEntity, err)
	}

	return req, nil
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := updateReq{key: r.Header.Get("Authorization")}
	req.id = bone.GetValue(r, "id")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(bootstrap.ErrMalformedEntity, err)
	}

	return req, nil
}

func decodeUpdateCertRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := updateCertReq{
		key:     r.Header.Get("Authorization"),
		thingID: bone.GetValue(r, "id"),
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(bootstrap.ErrMalformedEntity, err)
	}

	return req, nil
}

func decodeUpdateConnRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := updateConnReq{key: r.Header.Get("Authorization")}
	req.id = bone.GetValue(r, "id")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(bootstrap.ErrMalformedEntity, err)
	}

	return req, nil
}

func decodeUnknownRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, errInvalidQueryParams
	}

	offset, limit, err := parsePagePrams(q)
	if err != nil {
		return nil, err
	}

	req := listReq{
		key:    r.Header.Get("Authorization"),
		filter: bootstrap.Filter{Unknown: true},
		offset: offset,
		limit:  limit,
	}

	return req, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, errInvalidQueryParams
	}

	offset, limit, err := parsePagePrams(q)
	if err != nil {
		return nil, err
	}

	filter := parseFilter(q)

	req := listReq{
		key:    r.Header.Get("Authorization"),
		filter: filter,
		offset: offset,
		limit:  limit,
	}

	return req, nil
}

func decodeBootstrapRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := bootstrapReq{
		id:  bone.GetValue(r, "external_id"),
		key: r.Header.Get("Authorization"),
	}

	return req, nil
}

func decodeStateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := changeStateReq{key: r.Header.Get("Authorization")}
	req.id = bone.GetValue(r, "id")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(bootstrap.ErrMalformedEntity, err)
	}

	return req, nil
}

func decodeEntityRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := entityReq{
		key: r.Header.Get("Authorization"),
		id:  bone.GetValue(r, "id"),
	}

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)
	if ar, ok := response.(mainflux.Response); ok {
		for k, v := range ar.Headers() {
			w.Header().Set(k, v)
		}

		w.WriteHeader(ar.Code())

		if ar.Empty() {
			return nil
		}
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeSecureRes(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	if b, ok := response.([]byte); ok {
		if _, err := w.Write(b); err != nil {
			return err
		}
	}
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch errorVal := err.(type) {
	case errors.Error:
		w.Header().Set("Content-Type", contentType)
		switch {
		case errors.Contains(errorVal, errUnsupportedContentType):
			w.WriteHeader(http.StatusUnsupportedMediaType)
		case errors.Contains(errorVal, errInvalidQueryParams):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Contains(errorVal, bootstrap.ErrMalformedEntity):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Contains(errorVal, bootstrap.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
		case errors.Contains(errorVal, bootstrap.ErrUnauthorizedAccess):
			w.WriteHeader(http.StatusForbidden)
		case errors.Contains(errorVal, bootstrap.ErrConflict):
			w.WriteHeader(http.StatusConflict)
		case errors.Contains(errorVal, bootstrap.ErrThings):
			w.WriteHeader(http.StatusServiceUnavailable)
		case errors.Contains(errorVal, io.EOF):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Contains(errorVal, io.ErrUnexpectedEOF):
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		if errorVal.Msg() != "" {
			if err := json.NewEncoder(w).Encode(errorRes{Err: errorVal.Msg()}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func parseUint(s string) (uint64, error) {
	if s == "" {
		return 0, nil
	}

	ret, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, errInvalidQueryParams
	}

	return ret, nil
}

func parsePagePrams(q url.Values) (uint64, uint64, error) {
	offset, err := parseUint(q.Get("offset"))
	q.Del("offset")
	if err != nil {
		return 0, 0, errors.Wrap(errInvalidOffsetParam, err)
	}

	limit, err := parseUint(q.Get("limit"))
	q.Del("limit")
	if err != nil {
		return 0, 0, errors.Wrap(errInvalidLimitParam, err)
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	if limit == 0 {
		limit = defaultLimit
	}

	return offset, limit, nil
}

func parseFilter(values url.Values) bootstrap.Filter {
	ret := bootstrap.Filter{
		FullMatch:    make(map[string]string),
		PartialMatch: make(map[string]string),
	}
	for k := range values {
		if contains(fullMatch, k) {
			ret.FullMatch[k] = values.Get(k)
		}
		if contains(partialMatch, k) {
			ret.PartialMatch[k] = strings.ToLower(values.Get(k))
		}
	}

	return ret
}

func contains(l []string, s string) bool {
	for _, v := range l {
		if v == s {
			return true
		}
	}
	return false
}
