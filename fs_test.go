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

package gitrepofs

import (
	"context"
	"io/fs"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/gitrepofs/test/helpers"
)

var _ = Describe("git directories", func() {

	var gfs fs.FS

	BeforeEach(func() {
		tree := Successful(commit.Tree())
		gfs = New(repo, tree, commit.Author.When)
	})

	When("opening files and directories", func() {

		When("successful", func() {

			It("opens the root directory", func() {
				names := []string{
					"README", "fodder", "folder",
				}

				direntries := Successful(fs.ReadDir(gfs, "."))
				Expect(direntries).To(HaveLen(len(names)))
				Expect(direntries).To(ContainElement(HaveField("Name()", BeElementOf(names))))
			})

			It("opens a directory", func() {
				direntries := Successful(fs.ReadDir(gfs, "folder/subfolder"))
				Expect(direntries).To(ContainElement(And(
					HaveField("Name()", "canary.txt"),
					HaveField("IsDir()", false))))
			})

			It("opens a file", func() {
				contents := Successful(fs.ReadFile(gfs, "folder/subfolder/canary.txt"))
				Expect(contents).To(ContainSubstring("chirp!"))
			})

		})

	})

	Context("failure", func() {

		DescribeTable("rejects invalid paths",
			func(name string, experr error) {
				_, err := gfs.Open(name)
				var perr *fs.PathError
				Expect(err).To(BeAssignableToTypeOf(perr))
				perr = err.(*fs.PathError)
				Expect(perr.Op).To(Equal("open"))
				Expect(perr.Path).To(Equal(name))
				Expect(perr.Err).To(Equal(experr))
			},
			Entry("empty name", "", fs.ErrInvalid),
			Entry("invalid name", "/a/b", fs.ErrInvalid),
			Entry("missing directory", "folder/folder/canary.txt", fs.ErrNotExist),
			Entry("missing file", "missing.txt", fs.ErrNotExist),
		)

		It("returns failures from helpers", func() {
			fs := gfs.(*FS)
			Expect(fs.openFile("missing.txt", object.TreeEntry{})).Error().To(HaveOccurred())
			Expect(fs.openDir("missing.txt", plumbing.Hash{})).Error().To(HaveOccurred())
		})

	})

	It("reports an error for a non-existing remote repository", func() {
		Expect(NewForRevision(context.Background(), "/", "invalidref")).Error().
			To(HaveOccurred())
	})

	It("reports an error for a non-existing reference", func() {
		Expect(NewForRevision(context.Background(), tmprepdir, "invalidref")).Error().
			To(HaveOccurred())
	})

	DescribeTable("returns an fs.FS for a repository and reference",
		func(ref string, hascanary bool) {
			gfs := Successful(NewForRevision(context.Background(), tmprepdir, ref))
			contents := Successful(fs.ReadFile(gfs, "README"))
			Expect(contents).To(ContainSubstring(`"remote" git repository`))
			if hascanary {
				Expect(fs.ReadFile(gfs, "folder/subfolder/canary.txt")).Error().NotTo(HaveOccurred())
			} else {
				Expect(fs.ReadFile(gfs, "folder/subfolder/canary.txt")).Error().To(HaveOccurred())
			}
		},
		Entry("v1.0", "v1.0", false),
		Entry("master", "master", true),
	)

})
