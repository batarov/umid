// Copyright (c) 2020 UMI
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package method_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"umid/jsonrpc"
	"umid/umid"
)

func TestGetBalance(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bc := &bcMock{}
	bc.FnBalance = func(a string) (*umid.Balance, error) {
		switch a {
		case "umi1aaa":
			return &umid.Balance{Confirmed: 1, Interest: 2, Unconfirmed: 3, Composite: nil}, nil
		case "umi1bbb":
			return &umid.Balance{Confirmed: 10, Interest: 20, Unconfirmed: 30, Composite: new(uint64)}, nil
		}

		return nil, errors.New("invalid address")
	}

	rpc := jsonrpc.NewRPC().SetBlockchain(bc)
	go rpc.Worker(ctx, &sync.WaitGroup{})

	tests := []struct {
		request  string
		response string
	}{
		{
			`{"jsonrpc":"2.0","method":"getBalance","id":1}`,
			`{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params"},"id":1}`,
		},
		{
			`{"jsonrpc":"2.0","method":"getBalance","params":[],"id":2}`,
			`{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params"},"id":2}`,
		},
		{
			`{"jsonrpc":"2.0","method":"getBalance","params":{"abc":1},"id":3}`,
			`{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params"},"id":3}`,
		},
		{
			`{"jsonrpc":"2.0","method":"getBalance","params":{"address":1},"id":4}`,
			`{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params"},"id":4}`,
		},
		{
			`{"jsonrpc":"2.0","method":"getBalance","params":{"address":"aaa"},"id":5}`,
			`{"jsonrpc":"2.0","error":{"code":-32603,"message":"Internal error"},"id":5}`,
		},
		{
			`{"jsonrpc":"2.0","method":"getBalance","params":{"address":"umi1aaa"},"id":6}`,
			`{"jsonrpc":"2.0","result":{"confirmed":1,"interest":2,"unconfirmed":3,"type":""},"id":6}`,
		},
		{
			`{"jsonrpc":"2.0","method":"getBalance","params":{"address":"umi1bbb"},"id":7}`,
			`{"jsonrpc":"2.0","result":{"confirmed":10,"interest":20,"unconfirmed":30,"composite":0,"type":""},"id":7}`,
		},
		{
			`[{"jsonrpc":"2.0","method":"getBalance","params":{"address":"umi1bbb"},"id":8}, 1]`,
			`[{"jsonrpc":"2.0","result":{"confirmed":10,"interest":20,"unconfirmed":30,"composite":0,"type":""},"id":8},{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"},"id":null}]`,
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequestWithContext(ctx, "POST", "/json-rpc", strings.NewReader(test.request))
		req.Header.Set("Content-Type", "application/json")

		res := httptest.NewRecorder()
		handler := http.HandlerFunc(rpc.HTTP)
		handler.ServeHTTP(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("wrong http code: got %v want %v", res.Code, http.StatusOK)
		}

		if res.Header().Get("Content-Type") != "application/json" {
			t.Errorf("wrong Content-Type header: got %v want %v", res.Header().Get("Content-Type"), "application/json")
		}

		if res.Body.String() != test.response {
			t.Errorf("unexpected body: got %v want %v", res.Body.String(), test.response)
		}
	}
}
