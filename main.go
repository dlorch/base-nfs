package main

import (
	"fmt"
	"os"

	"github.com/dlorch/nfsv3/mountv3"
	"github.com/dlorch/nfsv3/portmapv2"
	"github.com/dlorch/nfsv3/rpcv2"
)

/*
	NFSv3 Server

	BSD 2-Clause License

	Copyright (c) 2019, Daniel Lorch
	All rights reserved.

	Redistribution and use in source and binary forms, with or without
	modification, are permitted provided that the following conditions are met:

	1. Redistributions of source code must retain the above copyright notice, this
	   list of conditions and the following disclaimer.

	2. Redistributions in binary form must reproduce the above copyright notice,
       this list of conditions and the following disclaimer in the documentation
       and/or other materials provided with the distribution.

	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
	AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
	IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
	DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
	FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
	DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
	SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
	CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
	OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

func main() {
	portmapService := rpcv2.NewRPCService("portmap", portmapv2.Program, portmapv2.Version)

	err := portmapService.Listen("udp", ":111")

	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	// TODO investigate contexts to run services in the background: https://blog.golang.org/context
	go portmapService.HandleClients()

	mountService := rpcv2.NewRPCService("mount", mountv3.Program, mountv3.Version)

	err = mountService.Listen("tcp", ":892")

	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	mountService.HandleClients()
}
