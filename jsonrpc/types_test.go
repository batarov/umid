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

package jsonrpc_test

import "umid/umid"

type bcMock struct {
	FnBalance               func(string) (*umid.Balance, error)
	FnAddTransaction        func([]byte) error
	FnStructureByPrefix     func(string) (*umid.Structure, error)
	FnStructures            func() ([]umid.Structure, error)
	FnTransactionsByAddress func(string) ([]umid.Transaction, error)
	FnAddBlock              func([]byte) error
	FnLastBlockHeight       func() (uint32, error)
}

func (m *bcMock) Balance(s string) (*umid.Balance, error) {
	return m.FnBalance(s)
}

func (m *bcMock) AddTransaction(b []byte) error {
	return m.FnAddTransaction(b)
}

func (m *bcMock) StructureByPrefix(s string) (*umid.Structure, error) {
	return m.FnStructureByPrefix(s)
}

func (m *bcMock) Structures() ([]umid.Structure, error) {
	return m.FnStructures()
}

func (m *bcMock) TransactionsByAddress(s string) ([]umid.Transaction, error) {
	return m.FnTransactionsByAddress(s)
}

func (m *bcMock) LastBlockHeight() (uint32, error) {
	return m.FnLastBlockHeight()
}

func (m *bcMock) AddBlock(b []byte) error {
	return m.FnAddBlock(b)
}