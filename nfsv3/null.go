// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nfsv3

// VoidReply is an empty reply
type VoidReply struct{}

func nfsProcedure3Null(procedureArguments []byte) (interface{}, error) {
	return &VoidReply{}, nil
}
