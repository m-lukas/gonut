// Copyright © 2019 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cf_test

import (
	"encoding/json"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/homeport/gonut/internal/gonut/cf"
)

var _ = Describe("Cloud Foundry JSON structs and contracts", func() {
	Context("Cloud Foundry API result JSON", func() {
		It("should parse Cloud Foundry API app details", func() {
			data, err := ioutil.ReadFile("../../../assets/test/cf-curl/v2/apps/nodejs-app.json")
			Expect(err).ToNot(HaveOccurred())

			var app AppDetails
			Expect(json.Unmarshal(data, &app)).ToNot(HaveOccurred())
			Expect(app.Entity.DetectedBuildpackGUID).To(BeEquivalentTo("6b70e2d7-1c63-4af9-b06d-37ae841ca8ae"))
		})

		It("should parse Cloud Foundry API buildpacks details", func() {
			data, err := ioutil.ReadFile("../../../assets/test/cf-curl/v2/buildpacks/nodejs-buildpack.json")
			Expect(err).ToNot(HaveOccurred())

			var buildpack BuildpackDetails
			Expect(json.Unmarshal(data, &buildpack)).ToNot(HaveOccurred())
			Expect(buildpack.Entity.Name).To(BeEquivalentTo("nodejs_buildpack"))
		})

		It("should parse Cloud Foundry API stacks details", func() {
			data, err := ioutil.ReadFile("../../../assets/test/cf-curl/v2/stacks/cflinuxfs3.json")
			Expect(err).ToNot(HaveOccurred())

			var stack StackDetails
			Expect(json.Unmarshal(data, &stack)).ToNot(HaveOccurred())
			Expect(stack.Entity.Name).To(BeEquivalentTo("cflinuxfs3"))
			Expect(stack.Entity.Description).To(BeEquivalentTo("Cloud Foundry Linux-based filesystem (Ubuntu 18.04)"))
		})
	})
})
