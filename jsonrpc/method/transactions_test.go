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
	"bytes"
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

func TestListTransaction(t *testing.T) {
	bc := &bcMock{}
	bc.FnTransactionsByAddress = func(s string) ([]*umid.Transaction, error) {
		if s == "umi1aaa" {
			return make([]*umid.Transaction, 0), nil
		}

		if s == "umi1bbb" {
			arr := make([]*umid.Transaction, 1)
			arr[0] = &umid.Transaction{Hash: strings.Repeat("00", 32)}

			return arr, nil
		}

		return nil, errors.New("invalid address")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rpc := jsonrpc.NewRPC().SetBlockchain(bc)
	go rpc.Worker(ctx, &sync.WaitGroup{})

	tests := []struct {
		request  string
		response string
	}{
		{
			`{"jsonrpc":"2.0","method":"listTransactions","id":1}`,
			`{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params"},"id":1}`,
		},
		{
			`{"jsonrpc":"2.0","method":"listTransactions","params":[],"id":2}`,
			`{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params"},"id":2}`,
		},
		{
			`{"jsonrpc":"2.0","method":"listTransactions","params":{"abc":1},"id":3}`,
			`{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params"},"id":3}`,
		},
		{
			`{"jsonrpc":"2.0","method":"listTransactions","params":{"address":1},"id":4}`,
			`{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params"},"id":4}`,
		},
		{
			`{"jsonrpc":"2.0","method":"listTransactions","params":{"address":"aaa"},"id":5}`,
			`{"jsonrpc":"2.0","error":{"code":-32603,"message":"Internal error"},"id":5}`,
		},
		{
			`{"jsonrpc":"2.0","method":"listTransactions","params":{"address":"umi1aaa"},"id":6}`,
			`{"jsonrpc":"2.0","result":[],"id":6}`,
		},
		{
			`{"jsonrpc":"2.0","method":"listTransactions","params":{"address":"umi1bbb"},"id":7}`,
			`{"jsonrpc":"2.0","result":[{"hash":"0000000000000000000000000000000000000000000000000000000000000000","block_height":0,"block_tx_idx":0,"version":0,"sender":""}],"id":7}`,
		},
		{
			`[{"jsonrpc":"2.0","method":"listTransactions","params":{"address":"umi1aaa"},"id":8}, 1]`,
			`[{"jsonrpc":"2.0","result":[],"id":8},{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"},"id":null}]`,
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

func TestSendTransaction(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bc := &bcMock{}
	bc.FnAddTransaction = func(s []byte) error {
		if bytes.Equal(make([]byte, 150), s) {
			return errors.New("invalid transaction")
		}

		return nil
	}

	rpc := jsonrpc.NewRPC().SetBlockchain(bc)
	go rpc.Worker(ctx, &sync.WaitGroup{})

	tests := []struct {
		request  string
		response string
	}{
		{
			`{"jsonrpc":"2.0","method":"sendTransaction","id":1}`,
			`{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params"},"id":1}`,
		},
		{
			`{"jsonrpc":"2.0","method":"sendTransaction","params":[],"id":2}`,
			`{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params"},"id":2}`,
		},
/*
		{
			`{"jsonrpc":"2.0","method":"sendTransaction","params":{"abc":1},"id":3}`,
			`{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params"},"id":3}`,
		},
		{
			`{"jsonrpc":"2.0","method":"sendTransaction","params":{"address":1},"id":4}`,
			`{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params"},"id":4}`,
		},
 */
		{
			`{"jsonrpc":"2.0","method":"sendTransaction","params":{"base64":"a"},"id":5}`,
			`{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params"},"id":5}`,
		},
		{
			`{"jsonrpc":"2.0","method":"sendTransaction","params":{"base64":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"},"id":6}`,
			`{"jsonrpc":"2.0","error":{"code":-32602,"message":"invalid transaction"},"id":6}`,
		},
		{
			`{"jsonrpc":"2.0","method":"sendTransaction","params":{"base64":"BAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"},"id":7}`,
			`{"jsonrpc":"2.0","result":{"hash":"cf9b9de9829bc7a8979d0969fbee6321c6bbf4ca3de72796f9aa90e5a57b7642"},"id":7}`,
		},
		{
			`{"jsonrpc":"2.0","method":"sendTransaction","params":{"base64":"!AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"},"id":8}`,
			`{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params"},"id":8}`,
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
