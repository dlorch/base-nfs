// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xdr

func newStructTagState() *structTagState {
	return new(structTagState)
}

func (sts *structTagState) switchStatement(u uint32) {
	sts.isSwitch = true
	sts.switchValue = u
	sts.isCase = false
	sts.matched = false
}

func (sts *structTagState) caseStatement(u uint32) {
	if sts.switchValue == uint32(u) {
		sts.isCase = true
		sts.currentCase = uint32(u)
		sts.matched = true
	} else {
		sts.isCase = false
	}
}

func (sts *structTagState) defaultStatement() {
	// if a previous case matched, the default statement will not be executed
	// if no previous case matched, the default statement will be executed
	sts.isCase = !sts.matched
}

func (sts *structTagState) caseMatch() bool {
	return !sts.isSwitch || (sts.isSwitch && sts.isCase)
}
