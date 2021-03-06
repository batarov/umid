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

package network

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
	"umid/umid"
)

const pullIntervalSec = 5

func (net *Network) puller(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(pullIntervalSec * time.Second):
			pull(ctx, net.client, net.blockchain)
		}
	}
}

func pull(ctx context.Context, client *http.Client, bc umid.IBlockchain) {
	lstBlkHeight, err := bc.LastBlockHeight()
	if err != nil {
		return
	}

	const tpl = `{"jsonrpc":"2.0","method":"listBlocks","params":{"height":%d},"id":"%d"}`
	jsn := fmt.Sprintf(tpl, lstBlkHeight+1, time.Now().UnixNano())

	req, _ := http.NewRequestWithContext(ctx, "POST", peer(), strings.NewReader(jsn))

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	if err != nil {
		return
	}

	if processResponse(body, bc) > 0 {
		pull(ctx, client, bc)
	}
}

func processResponse(body []byte, bc umid.IBlockchain) (cnt int) {
	res := new(struct {
		Result [][]byte `json:"result"`
	})

	if err := json.Unmarshal(body, res); err != nil {
		return
	}

	for _, b := range res.Result {
		if err := bc.AddBlock(b); err != nil {
			return
		}
	}

	return len(res.Result)
}
