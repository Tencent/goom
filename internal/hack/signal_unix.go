// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build aix || darwin || dragonfly || freebsd || js || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd js linux netbsd openbsd solaris

package hack

import (
	"os"
	"syscall"
)

// SignalsToIgnore ignore the quit signal
var SignalsToIgnore = []os.Signal{os.Interrupt, syscall.SIGQUIT}
