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
	"io/fs"

	"github.com/go-git/go-git/v5/plumbing/filemode"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/gitrepofs/test/helpers"
)

var _ = Describe("directory entry information", func() {

	It("returns a file dir entry", func() {
		const filename = "README"

		tree := Successful(commit.Tree())
		readme := Successful(tree.FindEntry(filename))
		size := Successful(tree.Size(filename))

		de := (fs.DirEntry)(NewDirEntry(*readme, size, commit.Author.When))
		Expect(de).NotTo(BeNil())

		Expect(de.Name()).To(Equal(filename))
		Expect(de.IsDir()).To(BeFalse())
		Expect(de.Type()).To(Equal(Successful(filemode.Regular.ToOSFileMode())))

		fi := Successful(de.Info())
		Expect(fi.Name()).To(Equal(filename))
		Expect(fi.Size()).To(Equal(size))
		Expect(fi.Mode()).To(Equal(Successful(filemode.Regular.ToOSFileMode())))
		Expect(fi.ModTime()).To(Equal(commit.Author.When))
		Expect(fi.IsDir()).To(BeFalse())
		Expect(fi.Sys()).To(BeNil())
	})

	It("returns a file dir entry", func() {
		const filename = "folder"

		tree := Successful(commit.Tree())
		readme := Successful(tree.FindEntry(filename))
		size := Successful(tree.Size(filename))

		de := (fs.DirEntry)(NewDirEntry(*readme, size, commit.Author.When))
		Expect(de).NotTo(BeNil())

		Expect(de.Name()).To(Equal(filename))
		Expect(de.IsDir()).To(BeTrue())
		Expect(de.Type()).To(Equal(Successful(filemode.Dir.ToOSFileMode())))
	})

})
