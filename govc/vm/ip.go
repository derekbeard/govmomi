/*
Copyright (c) 2014 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vm

import (
	"flag"
	"fmt"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/govc/cli"
	"github.com/vmware/govmomi/govc/flags"
	"github.com/vmware/govmomi/govc/host/esxcli"
)

type ip struct {
	*flags.OutputFlag
	*flags.SearchFlag

	esx bool
}

func init() {
	cli.Register("vm.ip", &ip{})
}

func (cmd *ip) Register(f *flag.FlagSet) {
	cmd.SearchFlag = flags.NewSearchFlag(flags.SearchVirtualMachines)
	f.BoolVar(&cmd.esx, "esxcli", false, "Use esxcli instead of guest tools")
}

func (cmd *ip) Process() error { return nil }

func (cmd *ip) Run(f *flag.FlagSet) error {
	c, err := cmd.Client()
	if err != nil {
		return err
	}

	vms, err := cmd.VirtualMachines(f.Args())
	if err != nil {
		return err
	}

	var get func(*govmomi.VirtualMachine) (string, error)

	if cmd.esx {
		get = esxcli.NewGuestInfo(c).IpAddress
	} else {
		get = func(vm *govmomi.VirtualMachine) (string, error) {
			return vm.WaitForIP(c)
		}
	}

	for _, vm := range vms {
		ip, err := get(vm)
		if err != nil {
			return err
		}

		// TODO(PN): Display inventory path to VM
		fmt.Fprintf(cmd, "%s\n", ip)
	}

	return nil
}