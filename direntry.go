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

	"github.com/go-git/go-git/v5/plumbing/object"
)

var _ fs.DirEntry = (*DirEntry)(nil)

// DirEntry represents an entry read from a directory. They are returned by
// [fs.ReadDir] and [fs.ReadDirFile.ReadDir]. DirEntry objects also contain
// [fs.FileInfo] objects.
type DirEntry struct {
	fileinfo *FileInfo
}

// NewDirEntry returns a new DirEntry object for a file or directory entry.
func NewDirEntry(
	entry object.TreeEntry,
	size int64,
	mtime time.Time,
) *DirEntry {
	return &DirEntry{
		fileinfo: NewFileInfo(entry, size, mtime),
	}
}

// Name returns the name of the file (or subdirectory) described by the entry.
// This name is only the final element of the path (the base name), not the
// entire path. For example, Name would return "hello.go" not
// "home/gopher/hello.go".
func (e *DirEntry) Name() string { return e.fileinfo.Name() }

// IsDir reports whether the entry describes a directory.
func (e *DirEntry) IsDir() bool { return e.fileinfo.IsDir() }

// Type returns the type bits for the entry. The type bits are a subset of the
// usual [fs.FileMode] bits, those returned by the [fs.FileMode.Type] method.
func (e *DirEntry) Type() fs.FileMode { return e.fileinfo.Mode() }

// Info returns the FileInfo for the file or subdirectory described by the
// entry. The returned FileInfo may be from the time of the original directory
// read or from the time of the call to Info. If the file has been removed or
// renamed since the directory read, Info may return an error satisfying
// [errors.Is](err, ErrNotExist). If the entry denotes a symbolic link, Info
// reports the information about the link itself, not the link's target.
func (e *DirEntry) Info() (fs.FileInfo, error) { return e.fileinfo, nil }
