// Copyright 2020 Go Noto Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Noto is a trademark of Google Inc. Noto fonts are open source.
// All Noto fonts are published under the SIL Open Font License, Version 1.1.

// package notomono provides the "Noto Mono" font collection. It is a fixed-width, serif font.
// This font collection provides broad unicode coverage.
// Special software is required to use OpenType font collections.
//
// See https://github.com/gonoto/gonoto for details.
package notomono

import (
	"compress/gzip"
	"io"
	"sync"
)

type chunkDecoder struct{}

func (d chunkDecoder) Read(p []byte) (n int, err error) {
	for len(p) >= 8 {
		if len(chunks) < 1 {
			return n, io.EOF
		}
		if len(chunks[0]) < 1 {
			chunks = chunks[1:]
			continue
		}
		u := chunks[0][0]
		chunks[0] = chunks[0][1:]
		p[0] = byte(u & 0xff)
		p[1] = byte(u & 0xff00 >> 8)
		p[2] = byte(u & 0xff0000 >> 16)
		p[3] = byte(u & 0xff000000 >> 24)
		p[4] = byte(u & 0xff00000000 >> 32)
		p[5] = byte(u & 0xff0000000000 >> 40)
		p[6] = byte(u & 0xff000000000000 >> 48)
		p[7] = byte(u & 0xff00000000000000 >> 56)
		p = p[8:]
		n += 8
	}
	return n, nil
}

var initOnce sync.Once
var otcData []byte

// OTC returns the font data as an OpenType collection.
func OTC() []byte {
	initOnce.Do(func() {
		var cr chunkDecoder
		otcData = make([]byte, decompressedSize)
		r, _ := gzip.NewReader(cr)
		_, _ = io.ReadFull(r, otcData)
		chunks = nil
	})
	return otcData
}
