// params_test.go
//
// Copyright (C) 2015 Selection Pressure LLC
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package ionconnect

import (
	"flag"
	"github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/codegangsta/cli"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Params", func() {
	var (
		context *cli.Context
		config  Command
		set     *flag.FlagSet
	)

	BeforeEach(func() {
		Debug = true
		set = flag.NewFlagSet("set", 0)
		command := cli.Command{Name: "scan-git"}
		context = cli.NewContext(nil, set, nil)
		context.Command = command

		config, _ = GetConfig().FindSubCommandConfig("scanner", "scan-git")
	})

	Context("When generating Post Params", func() {
		It("should populate the fields from the context flags", func() {
			args := []string{"ernie"}
			params := PostParams{}.Generate(args, config.Args)
			Expect(params).To(Equal(PostParams{Project: "ernie"}))
			Expect(params.String()).To(Equal("Project=ernie, Url=, Type=, Checksum=, Id=, Text=, Version=, File="))
		})

		It("should populate the fields from json data", func() {
			config, _ = GetConfig().FindSubCommandConfig("test", "test-json")
			args := []string{"{\"key\":\"value\"}", "[\"key\",\"value\"]"}
			params := PostParams{}.Generate(args, config.Args)
			m := make(map[string]interface{})
			m["key"] = "value"

			a := make([]interface{}, 2)
			a[0] = "key"
			a[1] = "value"

			Expect(params).To(Equal(PostParams{Results: m, Rules: a}))
		})

		It("should populate the fields with file contents if type is file ", func() {
			config, _ = GetConfig().FindSubCommandConfig("scanner", "scan-dependencies")
			args := []string{"../test/Gemfile"}
			params := PostParams{}.Generate(args, config.Args)
			Expect(params).To(Equal(PostParams{File: "source 'https://rubygems.org'\nruby '2.2.4'\n\ngem 'sinatra', '1.4.6', require: 'sinatra/base'\ngem 'sinatra-contrib', '1.4.6'\ngem 'rack-standards', '0.0.1'\ngem 'i18n', '0.7.0'\ngem 'rake', '10.4.2'\ngem 'json', '1.8.3'\ngem 'httpclient', '2.7.0.1'\ngem 'patron', '0.5.0'\ngem 'foreman', '0.78.0'\ngem 'thin', '1.6.4'\n\n# geocoding\ngem 'geocoder', '1.2.14'\n# languge\ngem 'cld2', '1.0.3', require: 'cld'\n# sentiment\ngem 'sentimental', '1.0.4'\n# dependencies\ngem 'lowendinsight', github: 'ion-channel/lowendinsight', :branch => 'master'\n\n# aws\ngem 'aws-sdk', '2.2.9'\n\ngem 'bundler', '>1.10.6'\n\ngroup :test do\n  gem 'minitest'\nend\n"}))
			Expect(params.String()).To(Equal(`Project=, Url=, Type=, Checksum=, Id=, Text=, Version=, File=source 'https://rubygems.org'
ruby '2.2.4'

gem 'sinatra', '1.4.6', require: 'sinatra/base'
gem 'sinatra-contrib', '1.4.6'
gem 'rack-standards', '0.0.1'
gem 'i18n', '0.7.0'
gem 'rake', '10.4.2'
gem 'json', '1.8.3'
gem 'httpclient', '2.7.0.1'
gem 'patron', '0.5.0'
gem 'foreman', '0.78.0'
gem 'thin', '1.6.4'

# geocoding
gem 'geocoder', '1.2.14'
# languge
gem 'cld2', '1.0.3', require: 'cld'
# sentiment
gem 'sentimental', '1.0.4'
# dependencies
gem 'lowendinsight', github: 'ion-channel/lowendinsight', :branch => 'master'

# aws
gem 'aws-sdk', '2.2.9'

gem 'bundler', '>1.10.6'

group :test do
  gem 'minitest'
end
`))
		})
	})

	Context("When generating Get Params", func() {
		It("should populate the fields from the context flags", func() {
			args := []string{"ernie"}
			options := map[string]string{
				"limit":  "22",
				"offset": "105",
			}
			params := GetParams{}.Generate(args, config.Args).UpdateFromMap(options)
			Expect(params).To(Equal(GetParams{Project: "ernie", Limit: "22", Offset: "105"}))
			Expect(params.String()).To(Equal("Project=ernie, Url=, Type=, Checksum=, Id=, Text=, Version=, Limit=22, Offset=105"))
		})
	})

	AfterEach(func() {
		Debug = false
	})
})
