// Copyright (c) 2015-2021 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package prompt

import (
	"bufio"
	"fmt"
)

func ProvideSeed() ([]byte, error) {
	return nil, fmt.Errorf("prompt not supported in WebAssembly")
}

func ProvidePrivPassphrase() ([]byte, error) {
	return nil, fmt.Errorf("prompt not supported in WebAssembly")
}

func PrivatePass(_ *bufio.Reader) ([]byte, error) {
	return nil, fmt.Errorf("prompt not supported in WebAssembly")
}

func PublicPass(_ *bufio.Reader, _, _, _ []byte) ([]byte, error) {
	return nil, fmt.Errorf("prompt not supported in WebAssembly")
}

func Seed(_ *bufio.Reader) ([]byte, error) {
	return nil, fmt.Errorf("prompt not supported in WebAssembly")
}
