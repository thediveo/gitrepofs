// Copyright 2023 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package version

import (
	"context"

	"github.com/thediveo/gitrepofs/test/localremote"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("finding the latest version in a remote repository", func() {

	Context("tag matcher for prefixed versions", func() {

		It("ignores non-matching and invalid references", func() {
			ptm := NewPrefixedTagMatcher("libfoo-")
			Expect(ptm("refs/heads/ohlala")).To(BeEmpty())
			Expect(ptm("refs/heads/libfoo1.2")).To(BeEmpty())
			Expect(ptm("refs/tags/libfoo1.2")).To(BeEmpty())
		})

		It("matches version references", func() {
			ptm := NewPrefixedTagMatcher("libfoo-")
			Expect(ptm("refs/tags/libfoo-1.2")).To(Equal("v1.2"))
			Expect(ptm("refs/tags/libfoo-v1.2")).To(Equal("v1.2"))
			Expect(ptm("refs/tags/libfoo-v1.2.3")).To(Equal("v1.2.3"))
		})
	})

	Context("with a remote repository", Ordered, func() {

		var tmprepdir string

		BeforeAll(func() {
			tmprepdir = localremote.CreateTransientTestRepo()
		})

		It("reports an error when tag matcher doesn't match", func(ctx context.Context) {
			Expect(LatestReleaseTag(ctx, tmprepdir, func(refname string) (semver string) {
				return ""
			})).Error().To(HaveOccurred())
		})

		It("finds latest and greatest version", func(ctx context.Context) {
			semver, ref, err := LatestReleaseTag(ctx, tmprepdir, SemverTagMatcher)
			Expect(err).NotTo(HaveOccurred())
			Expect(semver).To(Equal("v1.1.1"))
			Expect(ref).To(Equal("refs/tags/v1.1.1"))
		})

	})

})
