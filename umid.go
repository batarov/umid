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

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"umid/blockchain"
	"umid/jsonrpc"
	"umid/network"
	"umid/storage"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(os.Stdout)

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	db := storage.NewStorage()
	bc := blockchain.NewBlockchain().SetStorage(db)
	rpc := jsonrpc.NewRPC().SetBlockchain(bc)
	net := network.NewNetwork().SetBlockchain(bc)
	srv := network.NewServer()

	http.HandleFunc("/json-rpc", jsonrpc.CORS(jsonrpc.Filter(rpc.HTTP)))
	http.HandleFunc("/json-rpc-ws", rpc.WebSocket)

	go db.Worker(ctx, wg)
	go bc.Worker(ctx, wg)
	go rpc.Worker(ctx, wg)
	go net.Worker(ctx, wg)
	go srv.Serve()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	srv.DrainConnections()
	cancel()
	srv.Shutdown()

	wg.Wait()
}
