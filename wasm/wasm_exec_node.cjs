// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Trimmed from $GOROOT/lib/wasm/wasm_exec_node.js: only the Node globals setup,
// without the "run the wasm binary as a subprocess" logic (tests load it in-process).
'use strict'

globalThis.require = require
globalThis.fs = require('fs')
globalThis.TextEncoder = require('util').TextEncoder
globalThis.TextDecoder = require('util').TextDecoder

globalThis.performance ??= require('performance')

globalThis.crypto ??= require('crypto')
