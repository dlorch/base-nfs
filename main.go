// Copyright 2019 Daniel Lorch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/dlorch/base-nfs/mountv3"
	"github.com/dlorch/base-nfs/nfsv3"
	"github.com/dlorch/base-nfs/portmapv2"
)

func main() {
	portmapService := portmapv2.NewPortmapService()

	err := portmapService.AddListener("udp", ":111")

	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	err = portmapService.AddListener("tcp", ":111")

	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	go portmapService.HandleClients()

	mountService := mountv3.NewMountService()

	err = mountService.AddListener("tcp", ":892")

	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	go mountService.HandleClients()

	nfsv3Service := nfsv3.NewNFSv3Service()

	err = nfsv3Service.AddListener("tcp", ":2049")

	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	go nfsv3Service.HandleClients()

	portmapService.WaitUntilDone()
	mountService.WaitUntilDone()
	nfsv3Service.WaitUntilDone()
}
