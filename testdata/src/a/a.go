// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package a

import (
	"context"
	"time"
)

func MyContextIsPassedAsContextContext() {
	mc := myContext{}
	ContextContextIsUsedInParams(mc) // want "do not use a.myContext as context.Context in function call"
}

func ContextContextIsUsedInParams(ctx context.Context) {
	ctx.Done()
}

type myContext struct{}

func (mc myContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}
func (mc myContext) Done() <-chan struct{} {
	return nil
}
func (mc myContext) Err() error {
	return nil
}
func (mc myContext) Value(key any) any {
	return nil
}
