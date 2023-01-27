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
	"io/fs"
	"time"

	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var _ fs.FileInfo = (*FileInfo)(nil)

// FileInfo implements fs.FileInfo for both a git (regular/executable) file blob
// as well as a directory.
//
// According to the “fs.FS Zoo” FileInfo objects are
// returned from:
//   - [fs.File.Stat]
//   - directory entries: [fs.DirEntry.Info]
type FileInfo struct {
	entry object.TreeEntry
	size  int64
	mtime time.Time
}

// NewFileInfo returns a new FileInfo object, given a git tree entry, file size
// and modification time stamp.
//
// The size can (and must) be determined beforehand given a tree object and then
// calling [object.Tree.Size] with the correct path.
func NewFileInfo(
	entry object.TreeEntry,
	size int64,
	mtime time.Time,
) *FileInfo {
	return &FileInfo{
		entry: entry,
		size:  size,
		mtime: mtime,
	}
}

// NewFileInfoFromTree returns a new FileInfo object that represents the
// specified tree as a directory itself.
func NewFileInfoFromTree(
	tree *object.Tree,
	name string,
	mtime time.Time,
) *FileInfo {
	return &FileInfo{
		entry: object.TreeEntry{
			Name: name,
			Mode: filemode.Dir,
			Hash: tree.Hash,
		},
		mtime: mtime,
	}
}

// Name returns the name of a file (that can actually also happened to be a
// directory).
func (f *FileInfo) Name() string { return f.entry.Name }

// Size returns the size of the file. Always returns 0 for a directory (this is
// system-dependent anyway).
func (f *FileInfo) Size() int64 { return f.size }

// Mode returns the file's mode and permission bits.
func (f *FileInfo) Mode() fs.FileMode {
	osmode, err := f.entry.Mode.ToOSFileMode()
	if err != nil {
		return fs.ModeIrregular
	}
	return osmode
}

// ModTime returns the file's modification time.
func (f *FileInfo) ModTime() time.Time { return f.mtime }

// IsDir returns true if the file actually is a directory.
func (f *FileInfo) IsDir() bool { return f.entry.Mode == filemode.Dir }

// Sys always return nil, as there is no-system specific file information
// provided by our git file system.
func (f *FileInfo) Sys() any { return nil }
