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

package localremote

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/storage/memory"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/gitrepofs/test/helpers"
)

var _ = Describe("transient local 'remote' repository", func() {

	var tmpdir string

	Context("transient repository life-cycle", Ordered, func() {

		It("creates a transient repository with correct contents", func() {
			tmpdir = CreateTransientTestRepo()
			Expect(tmpdir).NotTo(BeEmpty())
			Expect(tmpdir).To(BeADirectory())

			const tagName = "v1.1.1"
			repo := Successful(git.CloneContext(context.Background(),
				memory.NewStorage(),
				nil,
				&git.CloneOptions{
					URL:           tmpdir,
					SingleBranch:  true,
					ReferenceName: plumbing.NewTagReferenceName(tagName),
				}))
			tag := Successful(repo.Tag(tagName))
			commit := Successful(repo.CommitObject(tag.Hash()))
			tree := Successful(commit.Tree())

			canary := Successful(tree.FindEntry("folder/subfolder/canary.txt"))
			Expect(canary.Mode).To(Equal(filemode.Regular))

			schkript := Successful(tree.FindEntry("folder/subfolder/schkript.sh"))
			Expect(schkript.Mode).To(Equal(filemode.Executable))
		})

		It("has removed the transient repository", func() {
			Expect(tmpdir).NotTo(BeEmpty())
			Expect(tmpdir).NotTo(BeAnExistingFile()) // includes dirs, too.
		})

	})

})
