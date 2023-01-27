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

var _ = Describe("file/dir information", func() {

	It("returns information about a regular file", func() {
		const filename = "README"

		tree := Successful(commit.Tree())
		readme := Successful(tree.FindEntry(filename))
		size := Successful(tree.Size(filename))

		fi := (fs.FileInfo)(NewFileInfo(*readme, size, commit.Author.When))
		Expect(fi).NotTo(BeNil())

		Expect(fi.Name()).To(Equal(filename))
		Expect(fi.Size()).To(Equal(size))
		Expect(fi.Mode()).To(Equal(Successful(filemode.Regular.ToOSFileMode())))
		Expect(fi.ModTime()).To(Equal(commit.Author.When))
		Expect(fi.IsDir()).To(BeFalse())
		Expect(fi.Sys()).To(BeNil())
	})

	It("returns information about an executable file", func() {
		const filename = "folder/subfolder/schkript.sh"

		tree := Successful(commit.Tree())
		readme := Successful(tree.FindEntry(filename))
		size := Successful(tree.Size(filename))

		fi := (fs.FileInfo)(NewFileInfo(*readme, size, commit.Author.When))
		Expect(fi).NotTo(BeNil())

		Expect(fi.Mode()).To(Equal(Successful(filemode.Executable.ToOSFileMode())))
	})

	It("returns information about a directory", func() {
		const filename = "folder"

		tree := Successful(commit.Tree())
		subfolder := Successful(tree.FindEntry(filename))
		size := Successful(tree.Size(filename))

		fi := (fs.FileInfo)(NewFileInfo(*subfolder, size, commit.Author.When))
		Expect(fi).NotTo(BeNil())

		Expect(fi.Name()).To(Equal(filename))
		Expect(fi.Size()).To(Equal(size))
		Expect(fi.Mode()).To(Equal(Successful(filemode.Dir.ToOSFileMode())))
		Expect(fi.ModTime()).To(Equal(commit.Author.When))
		Expect(fi.IsDir()).To(BeTrue())
		Expect(fi.Sys()).To(BeNil())
	})

	It("returns information about a directory tree", func() {
		tree := Successful(commit.Tree())

		fi := (fs.FileInfo)(NewFileInfoFromTree(tree, ".", commit.Author.When))
		Expect(fi).NotTo(BeNil())

		Expect(fi.Name()).To(Equal("."))
		Expect(fi.Size()).To(Equal(int64(0)))
		Expect(fi.Mode()).To(Equal(Successful(filemode.Dir.ToOSFileMode())))
		Expect(fi.ModTime()).To(Equal(commit.Author.When))
		Expect(fi.IsDir()).To(BeTrue())
		Expect(fi.Sys()).To(BeNil())
	})

})
