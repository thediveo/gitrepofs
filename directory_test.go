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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/gitrepofs/test/helpers"
)

var _ = Describe("git directories", Ordered, func() {

	var dir *Directory

	BeforeEach(func() {
		tree := Successful(commit.Tree())
		dir = NewDirectory(
			tree,
			NewFileInfoFromTree(tree, ".", commit.Author.When))
	})

	It("won't read a directory after closing it", func() {
		Expect(dir.Close()).To(Succeed())
		entries, err := dir.ReadDir(-1)
		Expect(err).To(HaveOccurred())
		Expect(entries).To(BeNil())
	})

	It("never reads anything 'dir'ectly", func() {
		b := make([]byte, 256)
		n, err := dir.Read(b)
		Expect(err).To(Equal(io.EOF))
		Expect(n).To(BeZero())
	})

	It("reads a directory en bloc", func() {
		entries := Successful(dir.ReadDir(-1))
		Expect(entries).To(HaveLen(3))
		Expect(entries).To(ConsistOf(
			HaveField("Name()", "fodder"),
			HaveField("Name()", "folder"),
			HaveField("Name()", "README"),
		))

		entries, err := dir.ReadDir(-1)
		Expect(err).To(Equal(io.EOF))
		Expect(entries).To(BeNil())

		de := Successful(dir.Stat())
		Expect(de.Name()).To(Equal("."))
		Expect(de.IsDir()).To(BeTrue())
	})

	It("reads a directory piece-wise", func() {
		expecteds := []string{
			"fodder", "folder", "README",
		}
		previous := ""
		for i := 1; i <= len(expecteds); i++ {
			entries := Successful(dir.ReadDir(1))
			Expect(entries).To(HaveLen(1))
			Expect(entries[0]).To(
				HaveField("Name()", BeElementOf(expecteds)))
			Expect(entries[0]).To(HaveField("Name()", Not(Equal(previous))))
			previous = entries[0].Name()
		}
		entries, err := dir.ReadDir(1)
		Expect(err).To(Equal(io.EOF))
		Expect(entries).To(BeEmpty())
	})

})
