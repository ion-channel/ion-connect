// params_test.go
//
// Copyright (C) 2015 Selection Pressure LLC
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.


package ionconnect

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/codegangsta/cli"
  "flag"
)

var _ = Describe("Params", func() {
  var (
		context *cli.Context
		set *flag.FlagSet
	)

  BeforeEach(func() {
    Debug = true
    set = flag.NewFlagSet("set", 0)
		set.String("Name", "p", "the name")
    set.Set("Name", "ernie")
    command := cli.Command{Name: "scan-git"}
    command.Flags = []cli.Flag{cli.StringFlag{Name:"Name"}}
		context = cli.NewContext(nil, set, nil)
		context.Command = command


  })

  Context("When generating Post Params", func() {
    It("should populate the fields from the context flags", func() {
      Expect(PostParams{}.Generate(context)).To(Equal(PostParams{Name:"ernie"}))
    })
  })

  Context("When generating Get Params", func() {
    It("should populate the fields from the context flags", func() {
        Expect(GetParams{}.Generate(context)).To(Equal(GetParams{Name:"ernie"}))
    })
  })

  AfterEach(func() {
    Debug = false
  })
})
