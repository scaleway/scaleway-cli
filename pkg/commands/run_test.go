// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

func ExampleRun() {
	ctx := ExampleCommandContext()
	args := RunArgs{
		Image: "ubuntu-trusuty",
	}
	Run(ctx, args)
}

func ExampleRun_complex() {
	ctx := ExampleCommandContext()
	args := RunArgs{
		Attach:     false,
		Bootscript: "rescue",
		Command:    []string{"ls", "-la"},
		Detach:     false,
		Gateway:    "my-gateway",
		Image:      "ubuntu-trusty",
		Name:       "my-test-server",
		Tags:       []string{"testing", "fake"},
		Volumes:    []string{"50G", "1G"},
	}
	Run(ctx, args)
}
