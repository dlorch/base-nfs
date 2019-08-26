// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xdr_test

import (
	"fmt"

	"github.com/dlorch/base-nfs/xdr"
)

// /*
//  * Example from section 6 in RFC 1014
//  */
//
// const MAXUSERNAME = 32;     /* max length of a user name */
// const MAXFILELEN = 65535;   /* max length of a file      */
// const MAXNAMELEN = 255;     /* max length of a file name */
//
// /*
//  * Types of files:
//  */
// enum filekind {
//   TEXT = 0,       /* ascii data */
//   DATA = 1,       /* raw data   */
//   EXEC = 2        /* executable */
// };
//
// /*
// * File information, per kind of file:
// */
// union filetype switch (filekind kind) {
// case TEXT:
//    void;                           /* no extra information */
// case DATA:
//    string creator<MAXNAMELEN>;     /* data creator         */
// case EXEC:
//    string interpretor<MAXNAMELEN>; /* program interpretor  */
// };
//
// /*
//  * A complete file:
//  */
// struct file {
//    string filename<MAXNAMELEN>; /* name of file    */
//    filetype type;               /* info about file */
//    string owner<MAXUSERNAME>;   /* owner of file   */
//    opaque data<MAXFILELEN>;     /* file data       */
// };

const (
	MAXUSERNAME uint32 = 32    /* max length of a user name */
	MAXFILENAME uint32 = 65535 /* max length of a file      */
	MAXNAMELEN  uint32 = 255   /* max length of a file name */
)

const (
	TEXT uint32 = iota /* ascii data */
	DATA               /* raw data   */
	EXEC               /* executable */
)

// type and struct fields must start with capital letter to be marshallable by xdr
type Filetype struct {
	Filekind uint32 `xdr:"switch"`
	// note that "void" for "case TEXT" was omitted
	Creator     string `xdr:"case=1"` // note that size limit indicators need to be verified by the application
	Interpretor string `xdr:"case=2"`
}

type File struct {
	Filename string
	Type     Filetype
	Owner    string
	Data     []byte
}

func Example_sillyprog() {
	/*
		Suppose now that there is a user named "john" who wants to store his
		lisp program "sillyprog" that contains just the data "(quit)".  His
		file would be encoded as follows:

			OFFSET  HEX BYTES       ASCII    COMMENTS
			------  ---------       -----    --------
			0       00 00 00 09     ....     -- length of filename = 9
			4       73 69 6c 6c     sill     -- filename characters
			8       79 70 72 6f     ypro     -- ... and more characters ...
			12      67 00 00 00     g...     -- ... and 3 zero-bytes of fill
			16      00 00 00 02     ....     -- filekind is EXEC = 2
			20      00 00 00 04     ....     -- length of interpretor = 4
			24      6c 69 73 70     lisp     -- interpretor characters
			28      00 00 00 04     ....     -- length of owner = 4
			32      6a 6f 68 6e     john     -- owner characters
		    36      00 00 00 06     ....     -- length of file data = 6
			40      28 71 75 69     (qui     -- file data bytes ...
			44      74 29 00 00     t)..     -- ... and 2 zero-bytes of fill
	*/
	sillyprog := &File{
		Filename: "sillyprog",
		Type: Filetype{
			Filekind:    EXEC,
			Interpretor: "lisp",
		},
		Owner: "john",
		Data:  []byte("(quit)"),
	}
	b, err := xdr.Marshal(sillyprog)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("% x", b)
	// Output:
	// 00 00 00 09 73 69 6c 6c 79 70 72 6f 67 00 00 00 00 00 00 02 00 00 00 04 6c 69 73 70 00 00 00 04 6a 6f 68 6e 00 00 00 06 28 71 75 69 74 29 00 00
}
