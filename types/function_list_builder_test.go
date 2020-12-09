// Copyright (c) OpenFaaS Project 2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package types

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/openfaas/faas-provider/auth"
	"github.com/openfaas/faas-provider/types"
)

func TestBuildSingleMatchingFunction(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/system/namespaces" {
			namespaces := []string{"openfaas-fn"}
			bytesOut, _ := json.Marshal(namespaces)
			_, _ = w.Write(bytesOut)
		} else {
			var functions []types.FunctionStatus
			annotationMap := make(map[string]string)
			annotationMap["topic"] = "topic1"

			functions = append(functions, types.FunctionStatus{
				Name:        "echo",
				Annotations: &annotationMap,
				Namespace:   "openfaas-fn",
			})
			bytesOut, _ := json.Marshal(functions)
			_, _ = w.Write(bytesOut)
		}
	}))

	client := srv.Client()
	builder := FunctionLookupBuilder{
		Client:         client,
		GatewayURL:     srv.URL,
		TopicDelimiter: ",",
	}

	lookup, err := builder.Build()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(lookup) != 1 {
		t.Errorf("Lookup - want: %d items, got: %d", 1, len(lookup))
	}
}
func Test_Build_SingleFunctionNoDelimiter(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/system/namespaces" {
			namespaces := []string{"openfaas-fn"}
			bytesOut, _ := json.Marshal(namespaces)
			_, _ = w.Write(bytesOut)
		} else {
			var functions []types.FunctionStatus
			annotationMap := make(map[string]string)
			annotationMap["topic"] = "topic1"

			functions = append(functions, types.FunctionStatus{
				Name:        "echo",
				Annotations: &annotationMap,
			})
			bytesOut, _ := json.Marshal(functions)
			_, _ = w.Write(bytesOut)
		}
	}))

	client := srv.Client()
	builder := FunctionLookupBuilder{
		Client:     client,
		GatewayURL: srv.URL,
	}

	lookup, err := builder.Build()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(lookup) != 1 {
		t.Errorf("Lookup - want: %d items, got: %d", 1, len(lookup))
	}
}

func TestBuildMultiMatchingFunction(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/system/namespaces" {
			namespaces := []string{"openfaas-fn"}
			bytesOut, _ := json.Marshal(namespaces)
			_, _ = w.Write(bytesOut)
		} else {
			var functions []types.FunctionStatus
			annotationMap := make(map[string]string)
			annotationMap["topic"] = "topic1,topic2,topic3"

			functions = append(functions, types.FunctionStatus{
				Name:        "echo",
				Annotations: &annotationMap,
			})
			bytesOut, _ := json.Marshal(functions)
			_, _ = w.Write(bytesOut)
		}
	}))

	client := srv.Client()
	builder := FunctionLookupBuilder{
		Client:         client,
		GatewayURL:     srv.URL,
		TopicDelimiter: ",",
	}

	lookup, err := builder.Build()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(lookup) != 3 {
		t.Errorf("Lookup - want: %d items, got: %d", 3, len(lookup))
	}
}

func TestBuildNoFunctions(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var functions []types.FunctionStatus
		bytesOut, _ := json.Marshal(functions)
		_, _ = w.Write(bytesOut)
	}))

	client := srv.Client()
	builder := FunctionLookupBuilder{
		Client:         client,
		GatewayURL:     srv.URL,
		TopicDelimiter: ",",
	}

	lookup, err := builder.Build()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(lookup) != 0 {
		t.Errorf("Lookup - want: %d items, got: %d", 0, len(lookup))
	}
}

func Test_Build_JustDelim(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/system/namespaces" {
			namespaces := []string{"openfaas-fn"}
			bytesOut, _ := json.Marshal(namespaces)
			_, _ = w.Write(bytesOut)
		} else {

			var functions []types.FunctionStatus
			annotationMap := make(map[string]string)
			annotationMap["topic"] = ","

			functions = append(functions, types.FunctionStatus{
				Name:        "echo",
				Annotations: &annotationMap,
			})
			bytesOut, _ := json.Marshal(functions)
			_, _ = w.Write(bytesOut)
		}
	}))

	client := srv.Client()
	builder := FunctionLookupBuilder{
		Client:         client,
		GatewayURL:     srv.URL,
		TopicDelimiter: ",",
	}

	lookup, err := builder.Build()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(lookup) != 0 {
		t.Errorf("Lookup - want: %d items, got: %d", 0, len(lookup))
	}
}

func Test_Build_MultiMatchingFunctionBespokeDelim(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/system/namespaces" {
			namespaces := []string{"openfaas-fn"}
			bytesOut, _ := json.Marshal(namespaces)
			_, _ = w.Write(bytesOut)
		} else {
			var functions []types.FunctionStatus
			annotationMap := make(map[string]string)
			annotationMap["topic"] = "topic1|topic2|topic3,withcomma"

			functions = append(functions, types.FunctionStatus{
				Name:        "echo",
				Annotations: &annotationMap,
			})
			bytesOut, _ := json.Marshal(functions)
			_, _ = w.Write(bytesOut)
		}
	}))

	client := srv.Client()
	builder := FunctionLookupBuilder{
		Client:         client,
		GatewayURL:     srv.URL,
		TopicDelimiter: "|",
	}

	lookup, err := builder.Build()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(lookup) != 3 {
		t.Errorf("Lookup - want: %d items, got: %d", 3, len(lookup))
	}
}

func Test_appendServiceMap(t *testing.T) {
	var TestCases = []struct {
		Name               string
		Key                string
		Function           string
		Namespace          string
		InputServiceMap    map[string][]string
		ExpectedServiceMap map[string][]string
	}{
		{
			Name:            "Empty starting map - key with length",
			Key:             "newKey",
			Function:        "fnName",
			Namespace:       "openfaas-fn",
			InputServiceMap: map[string][]string{},
			ExpectedServiceMap: map[string][]string{
				"newKey": {"fnName.openfaas-fn"},
			},
		},
		{
			Name:               "Empty starting map - zero key length",
			Key:                "",
			Function:           "fnName",
			Namespace:          "openfaas-fn",
			InputServiceMap:    map[string][]string{},
			ExpectedServiceMap: map[string][]string{},
		},
		{
			Name:            "Populated starting map - key with length",
			Key:             "theKey",
			Function:        "newName",
			Namespace:       "fn",
			InputServiceMap: map[string][]string{"theKey": {"fnName"}},
			ExpectedServiceMap: map[string][]string{
				"theKey": {"fnName", "newName.fn"},
			},
		},
		{
			Name:            "Populated starting map - zero key length",
			Key:             "",
			Function:        "newName",
			Namespace:       "fn",
			InputServiceMap: map[string][]string{"theKey": {"fnName"}},
			ExpectedServiceMap: map[string][]string{
				"theKey": {"fnName"},
			},
		},
		{
			Name:            "Populated starting map - new key with length",
			Key:             "newKey",
			Function:        "newName",
			Namespace:       "fn",
			InputServiceMap: map[string][]string{"theKey": {"fnName"}},
			ExpectedServiceMap: map[string][]string{
				"theKey": {"fnName"},
				"newKey": {"newName.fn"},
			},
		},
		{
			Name:            "Populated starting map - existing key new function",
			Key:             "newKey",
			Function:        "secondName",
			Namespace:       "openfaas-fn",
			InputServiceMap: map[string][]string{"theKey": {"fnName"}, "newKey": {"newName"}},
			ExpectedServiceMap: map[string][]string{
				"theKey": {"fnName"},
				"newKey": {"newName", "secondName.openfaas-fn"},
			},
		},
	}

	for _, test := range TestCases {

		serviceMap := appendServiceMap(test.Key, test.Function, test.Namespace, test.InputServiceMap)

		if len(serviceMap) != len(test.ExpectedServiceMap) {
			t.Errorf("Testcase %s failed on serviceMap size. want - %d, got - %d", test.Name, len(test.ExpectedServiceMap), len(serviceMap))
		}

		for key := range serviceMap {

			if _, exists := test.ExpectedServiceMap[key]; !exists {
				t.Errorf("Testcase %s failed on serviceMap keys. found value - %s doesnt exist in expected", test.Name, key)
			}

			if len(serviceMap[key]) != len(test.ExpectedServiceMap[key]) {
				t.Errorf("Testcase %s failed on key slice size. want - %d, got - %d", test.Name, len(test.ExpectedServiceMap[key]), len(serviceMap[key]))
			}

			lookupMap := make(map[string]bool)
			for _, fn := range serviceMap[key] {
				lookupMap[fn] = true
			}

			for _, v := range test.ExpectedServiceMap[key] {
				if _, found := lookupMap[v]; !found {
					t.Errorf("Testcase %s failed on key slice values. found value - %s doesnt exist in expected", test.Name, v)
				}
			}
		}
	}
}

func Test_BuildMultipleNamespaceFunction(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/system/namespaces" {
			namespaces := []string{"openfaas-fn", "fn"}
			bytesOut, _ := json.Marshal(namespaces)
			_, _ = w.Write(bytesOut)
		} else {
			var functions []types.FunctionStatus
			annotationMap := make(map[string]string)
			annotationMap["topic"] = "topic1"

			functions = append(functions, types.FunctionStatus{
				Name:        "echo",
				Annotations: &annotationMap,
				Namespace:   "openfaas-fn",
			})
			bytesOut, _ := json.Marshal(functions)
			_, _ = w.Write(bytesOut)
		}
	}))

	client := srv.Client()
	builder := FunctionLookupBuilder{
		Client:         client,
		GatewayURL:     srv.URL,
		TopicDelimiter: ",",
	}

	lookup, err := builder.Build()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(lookup) != 1 {
		t.Errorf("Lookup - want: %d items, got: %d", 1, len(lookup))
	}
	functions, ok := lookup["topic1"]
	if !ok {
		t.Errorf("Topic %s does not exists", "topic1")
	}
	if len(functions) != 2 {
		t.Errorf("Topic %s - want: %d functions, got: %d", "topic1", 2, len(functions))
	}
}

func Test_GetFunctions(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/system/namespaces" {
			namespaces := []string{"openfaas-fn"}
			bytesOut, _ := json.Marshal(namespaces)
			_, _ = w.Write(bytesOut)
		} else {
			var functions []types.FunctionStatus
			if r.URL.Query().Get("namespace") == "openfaas-fn" {
				annotationMap := make(map[string]string)
				annotationMap["topic"] = "topic1"

				functions = append(functions, types.FunctionStatus{
					Name:        "echo",
					Annotations: &annotationMap,
					Namespace:   "openfaas-fn",
				})
			}
			bytesOut, _ := json.Marshal(functions)
			_, _ = w.Write(bytesOut)
		}
	}))

	client := srv.Client()
	builder := FunctionLookupBuilder{
		Client:         client,
		GatewayURL:     srv.URL,
		TopicDelimiter: ",",
	}

	functions, err := builder.getFunctions("openfaas-fn")
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(functions) != 1 {
		t.Errorf("Functions - want: %d items, got: %d", 1, len(functions))
	}
}

func Test_GetEmptyFunctions(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/system/namespaces" {
			namespaces := []string{"openfaas-fn"}
			bytesOut, _ := json.Marshal(namespaces)
			_, _ = w.Write(bytesOut)
		} else {
			var functions []types.FunctionStatus
			if r.URL.Query().Get("namespace") == "openfaas-fn" {
				annotationMap := make(map[string]string)
				annotationMap["topic"] = "topic1"

				functions = append(functions, types.FunctionStatus{
					Name:        "echo",
					Annotations: &annotationMap,
					Namespace:   "openfaas-fn",
				})

			}
			bytesOut, _ := json.Marshal(functions)
			_, _ = w.Write(bytesOut)
		}
	}))

	client := srv.Client()
	builder := FunctionLookupBuilder{
		Client:         client,
		GatewayURL:     srv.URL,
		TopicDelimiter: ",",
	}

	functions, err := builder.getFunctions("fn")
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(functions) != 0 {
		t.Errorf("Functions - want: %d items, got: %d", 0, len(functions))
	}
}

func TestFunctionLookupBuilder_getNamespaces(t *testing.T) {
	var srv *httptest.Server
	defer func() {
		if srv != nil {
			srv.Close()
		}
	}()

	type fields struct {
		handler     http.Handler
		credentials *auth.BasicAuthCredentials
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr error
	}{
		{
			name: "401 invalid credentials (empty user and password)",
			fields: fields{
				handler:     authHandler(t),
				credentials: &auth.BasicAuthCredentials{},
			},
			want:    nil,
			wantErr: fmt.Errorf("authentication failure against gateway: %s", http.StatusText(401)),
		},
		{
			name: "401 no credentials with gateway enforcing basic auth",
			fields: fields{
				handler:     authHandler(t),
				credentials: nil,
			},
			want:    nil,
			wantErr: fmt.Errorf("authentication failure against gateway: %s", http.StatusText(401)),
		},
		{
			name: "500 internal server error",
			fields: fields{
				handler:     internalSrvErrHandler(t),
				credentials: nil,
			},
			want:    nil,
			wantErr: fmt.Errorf("get namespaces unexpected HTTP response: %s", http.StatusText(500)),
		},
		{
			name: "no namespaces returned (empty body)",
			fields: fields{
				handler:     emptyBodyHandler(t),
				credentials: nil,
			},
			want:    nil,
			wantErr: fmt.Errorf("empty response body"),
		},
		{
			name: "namespaces returned",
			fields: fields{
				handler:     namespaceHandler(t),
				credentials: nil,
			},
			want:    []string{"openfaas-fn"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		srv = httptest.NewServer(tt.fields.handler)

		t.Run(tt.name, func(t *testing.T) {
			s := &FunctionLookupBuilder{
				GatewayURL:  srv.URL,
				Client:      srv.Client(),
				Credentials: tt.fields.credentials,
			}
			got, err := s.getNamespaces()
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("getNamespaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNamespaces() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func authHandler(t *testing.T) http.Handler {
	t.Helper()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u, p, ok := r.BasicAuth(); ok {
			if u != "" && p != "" {
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	})
}

func internalSrvErrHandler(t *testing.T) http.Handler {
	t.Helper()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	})
}

func emptyBodyHandler(t *testing.T) http.Handler {
	t.Helper()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func namespaceHandler(t *testing.T) http.Handler {
	t.Helper()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ns []string
		if r.URL.Path != "/system/namespaces" {
			http.NotFoundHandler()
		}

		ns = []string{"openfaas-fn"}
		b, err := json.Marshal(ns)
		if err != nil {
			t.Errorf("marshal JSON: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(b)
		if err != nil {
			t.Errorf("write response: %v", err)
			return
		}
	})
}
