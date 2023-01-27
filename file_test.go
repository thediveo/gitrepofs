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
	"io"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/gitrepofs/test/helpers"
)

var _ = Describe("git files", func() {

	It("reads from a regular file", func() {
		const filename = "README"

		tree := Successful(commit.Tree())
		readme := Successful(tree.FindEntry(filename))
		size := Successful(tree.Size(filename))
		blob := Successful(repo.BlobObject(readme.Hash))

		f := NewFile(
			NewFileInfo(*readme, size, commit.Author.When),
			blob)
		Expect(f).NotTo(BeNil())
		defer func() {
			Expect(f.Close()).To(Succeed())
		}()

		Expect(f.Stat()).To(And(
			HaveField("Name()", path.Base(filename)),
			HaveField("Size()", blob.Size),
			HaveField("Mode()", Successful(readme.Mode.ToOSFileMode())),
			HaveField("IsDir()", false),
		))

		contents := Successful(io.ReadAll(f))
		Expect(contents).To(HaveLen(int(size)))
		Expect(contents).To(ContainSubstring(`"remote" git repository`))
	})

})
