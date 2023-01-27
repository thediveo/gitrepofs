/*
package gitrepofs provides a fs.FS interface onto a remote git repository at a
specific tag.

The main envisioned use case is for “go generate”-based updates from upstream
repositories to fetch the latest C definitions without the need to integrate
upstream C libraries.

		remoteURL := "https://gohub.org/froozle/baduzle"
		latest, latestref, err := version.LatestReleaseTag(context.Background(),
	      remoteURL, version.SemverMatcher)
		gfs, err := NewForRevision(context.Background(), remoteURL, latestref)
		contents, err := fs.ReadFile(gfs, "some/useful/file.h")

# The fs.FS Zoo

The number of interfaces and their relationships in [fs.FS] look like a (small)
zoo and it is easy to get lost in what returns what, or extends that, or
whatever. So here's our own little fs.FS Zoo map. For whatever reason, there's
no cafeteria and no rest rooms (unless counting in [runtime.GC]).

[fs.FS] provides access to a hierarchical file system. Additional optional
interfaces (discussed next) then offer more functionality.

  - [fs.FS.Open] opens not only files, but also directories, returning an
    [fs.File]. However, in order to successfully open directories, either the
    file system itself must additionally implement the interface [fs.ReadDirFS],
    or the [fs.File] returned must also implement [fs.ReadDirFile].

[fs.File] provides access to a single file or directory; for directories, the
additional interface [fs.ReadDirFiles] should also be implemented (Golang
soundbite). We implement regular and executable file access in the [File] type
and directory access in the [Directory] type.

  - [fs.File.Stat] returns an [fs.FileInfo].
  - [fs.File.Read] reads the file contents; it doesn't return anything for
    directories.
  - [fs.File.Close] closes the [fs.File].

[fs.FileInfo] describes a file or directory and is returned by [fs.File.Stat],
but can also be returned from [fs.DirEntry.Info]. We implement the fs.FileInfo
interface in our aptly named [FileInfo] type.

  - Name
  - Size
  - Mode: file mode bits.
  - ModTime
  - IsDir
  - Sys: underlying data source of nil.

Next on to [fs.ReadDirFile]: it provides [fs.File] operations and on top of it
reading a directory. We implemented this interface in the [Directory] type.

  - [fs.ReadDirFile.ReadDir] returns a bunch of [fs.DirEntry] objects. To
    complicate things, callers are allowed to read only piecemeal wise.

[fs.DirEntry] is an entry from a directory. These entries also have
[fs.FileInfo] objects attached to them. We implement the interface in the
[DirEntry] type.

  - Name
  - IsDir
  - Type: file mode bits.
  - [fs.DirEntry.Info]: returns an [fs.FileInfo] about this directory entry.
*/
package gitrepofs
