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
	"context"
	"fmt"
	"io/fs"
	"path"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

var _ fs.FS = (*FS)(nil)

// FS provides a view into a specific git tree.
type FS struct {
	repo  *git.Repository
	tree  *object.Tree
	mtime time.Time
}

// NewForRevision returns a [fs.FS] git repository file system object that
// provides read access to the files in the remote repository at the specified
// tag. The repository gets downloaded into memory only, so there is no need for
// handling and cleaning up any temporary directories on the file system.
//
// revision can be (as supported by
// [github.com/go-git/go-git/v5/Repository.ResolveRevision]):
//   - HEAD
//   - branch
//   - tag
//   - ...
func NewForRevision(ctx context.Context, remoteURL string, revision string) (fs.FS, error) {
	repo, err := git.Clone(
		memory.NewStorage(),
		nil,
		&git.CloneOptions{
			URL: remoteURL,
		})
	if err != nil {
		return nil, fmt.Errorf(
			"cannot clone remote repository %q", remoteURL)
	}
	commitHash, err := repo.ResolveRevision(plumbing.Revision(revision))
	if err != nil {
		return nil, fmt.Errorf(
			"no such revision %q in remote repository %q",
			revision, remoteURL)
	}
	commit, err := repo.CommitObject(*commitHash)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid commit hash for reference %q in remote repository %q",
			revision, remoteURL)
	}
	tree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf(
			"invalid tree hash for reference %q  in remote repository %q",
			revision, remoteURL)
	}
	return New(repo, tree, commit.Author.When), nil
}

// New returns a [fs.FS] for the specified tree of the git repository object,
// and using the specified modification time.
func New(repo *git.Repository, tree *object.Tree, mtime time.Time) fs.FS {
	return &FS{
		repo:  repo,
		tree:  tree,
		mtime: mtime,
	}
}

// Open opens the named file or directory. The name must conform to the rules
// implemented in [fs.ValidPath]:
//   - unrooted, slash-separated path elements, like “x/y/z”, but not “/x/y/z”.
//     Double slashes as separators are invlid.
//   - the root (top-level) directory on its own is named “.”.
//   - otherwise, neither “.” nor “..' are allowed.
//   - finally, the empty name “” isn't allowed either.
//
// Please note that [fs.ReadDir] uses [fs.ReadDirFS.ReadDir] when available, but
// otherwise falls back to [fs.FS.Open].
//
// When Open returns an error, it is of type [*fs.PathError] with the Op field
// set to "open", the Path field set to name, and the Err field describing the
// problem.
//
// Open rejects attempts to open names that do not satisfy [fs.ValidPath](name),
// returning a [*fs.PathError with Err set to [fs.ErrInvalid] or
// [fs.ErrNotExist].
func (gfs *FS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{
			Op:   "open",
			Path: name, // report original name/path
			Err:  fs.ErrInvalid,
		}
	}
	var entry object.TreeEntry
	if name == "." {
		entry.Mode = filemode.Dir
	} else {
		e, err := gfs.tree.FindEntry(name)
		if err != nil { // reports object.ErrDirectoryNotFound, ErrEntryNotFound
			return nil, &fs.PathError{
				Op:   "open",
				Path: name,
				Err:  fs.ErrNotExist,
			}
		}
		entry = *e
	}
	switch entry.Mode {
	case filemode.Regular, filemode.Executable:
		return gfs.openFile(name, entry)
	case filemode.Dir:
		return gfs.openDir(name, entry.Hash)
	}
	return nil, &fs.PathError{
		Op:   "open",
		Path: name,
		Err:  fs.ErrInvalid,
	}
}

// openFile returns a File object for the specified file name+path. The
// name+path must have been validated before using [fs.ValidPath].
func (gfs *FS) openFile(name string, entry object.TreeEntry) (fs.File, error) {
	blob, err := gfs.repo.BlobObject(entry.Hash)
	if err != nil {
		return nil, &fs.PathError{
			Op:   "open",
			Path: name,
			Err:  fs.ErrNotExist, // now that is embarrassing
		}
	}
	f := NewFile(NewFileInfo(entry, blob.Size, gfs.mtime), blob)
	if f == nil {
		return nil, &fs.PathError{
			Op:   "open",
			Path: name,
			Err:  fs.ErrInvalid,
		}
	}
	return f, nil
}

// openDir returns a Directory object for the specified directory path. The
// name+path must have been validated before using [fs.ValidPath]. In case of a
// non-root (non-toplevel) directory, h specifies the tree's hash, as commonly
// found in TreeEntry objects.
func (gfs *FS) openDir(name string, h plumbing.Hash) (f fs.File, err error) {
	var tree *object.Tree
	if name == "." {
		tree = gfs.tree
	} else {
		tree, err = gfs.repo.TreeObject(h)
		if err != nil {
			return nil, &fs.PathError{
				Op:   "open",
				Path: name,
				Err:  fs.ErrNotExist, // now that is embarrassing
			}
		}
	}
	entry := object.TreeEntry{
		Name: path.Base(name),
		Mode: filemode.Dir,
		Hash: tree.ID(), // ... albeit we actually don't need it.
	}
	return NewDirectory(tree, NewFileInfo(entry, 0, gfs.mtime)), nil
}
