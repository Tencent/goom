// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build plan9 || windows
// +build plan9 windows

package hack

import (
	"os"
)

// signalsToIgnore ignore the quit signal
var SignalsToIgnore = []os.Signal{os.Interrupt}
