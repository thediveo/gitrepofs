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
	"io"
	"io/fs"

	"github.com/go-git/go-git/v5/plumbing/object"
)

var _ fs.File = (*File)(nil)

// File represents a git regular or executable file. It never represents a
// directory, that is served by [Directory] instead.
type File struct {
	fileinfo *FileInfo
	r        io.ReadCloser
}

// NewFile returns a new File object, given a file information object and the
// file's contents blob.
func NewFile(fileinfo *FileInfo, blob *object.Blob) *File {
	r, err := blob.Reader()
	if err != nil {
		return nil
	}
	return &File{
		fileinfo: fileinfo,
		r:        r,
	}
}

// Stat returns information about this git file.
func (f *File) Stat() (fs.FileInfo, error) { return f.fileinfo, nil }

// Read some amount of contents from this git file.
func (f *File) Read(b []byte) (int, error) { return f.r.Read(b) }

// Close this git file.
func (f *File) Close() error { return f.r.Close() }
