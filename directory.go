// Copyright 2023 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gitrepofs

import (
	"errors"
	"io"
	"io/fs"

	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var _ fs.ReadDirFile = (*Directory)(nil)
var _ fs.File = (*Directory)(nil)

// Directory represents completely unexpectedly a git directory.
type Directory struct {
	tree     *object.Tree
	fileinfo *FileInfo
	index    int
}

// NewDirectory returns a new Directory object representing a git tree.
func NewDirectory(
	tree *object.Tree,
	fileinfo *FileInfo,
) *Directory {
	return &Directory{
		tree:     tree,
		fileinfo: fileinfo,
	}
}

// Stat returns information about this git directory.
func (d *Directory) Stat() (fs.FileInfo, error) { return d.fileinfo, nil }

func (d *Directory) ReadDir(n int) ([]fs.DirEntry, error) {
	if d.index < 0 {
		return nil, errors.New("closed directory")
	}
	if n <= 0 {
		n = len(d.tree.Entries) - d.index
		if n <= 0 {
			return nil, io.EOF
		}
	}
	count := n
	if d.index+count > len(d.tree.Entries) {
		count = len(d.tree.Entries) - d.index
	}
	if count == 0 {
		return nil, io.EOF
	}
	// nota bene: git trees are never empty.
	fileinfos := make([]fs.DirEntry, 0, count)
	for ; count > 0; count-- {
		var size int64
		entry := d.tree.Entries[d.index]
		if entry.Mode == filemode.Regular || entry.Mode == filemode.Executable {
			size, _ = d.tree.Size(entry.Name)
		}
		fileinfos = append(fileinfos,
			NewDirEntry(entry, size, d.fileinfo.mtime))
		d.index++
	}
	return fileinfos, nil
}

// Read nothing from this git directory.
func (d *Directory) Read(b []byte) (int, error) { return 0, io.EOF }

// Close this git file.
func (d *Directory) Close() error { d.index = -1; return nil }
