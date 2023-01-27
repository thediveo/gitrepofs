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

package localremote

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"

	. "github.com/onsi/ginkgo/v2"                  // we don't care about dot-imports
	. "github.com/onsi/gomega"                     // we don't care about dot-imports
	. "github.com/thediveo/gitrepofs/test/helpers" // guess what we don't care about?
)

var (
	fileMode, _ = filemode.Regular.ToOSFileMode()
	exeMode, _  = filemode.Executable.ToOSFileMode()
	dirMode, _  = filemode.Dir.ToOSFileMode()
)

const gitDir = ".git"

var commitOptions = &git.CommitOptions{
	Author: &object.Signature{
		Name:  "Brian",
		Email: "brian@palace.herodes",
		When:  time.Now(),
	},
}

//go:embed files
var contentfs embed.FS

// CreateTransientTestRepo initializes and populates a fresh git repository in a
// new temporary directory and then returns the path to this newly created
// directory.
func CreateTransientTestRepo() (repopath string) {
	By("creating a temporary directory to initialize a new git repository in")
	tmpdir := Successful(os.MkdirTemp("", "localremote-*"))
	DeferCleanup(func() {
		Expect(os.RemoveAll(tmpdir)).To(Succeed())
	})

	By("initializing git repository")
	repo := Successful(git.PlainInit(tmpdir, false))
	Expect(path.Join(tmpdir, gitDir)).To(BeADirectory())

	By("checking in and tagging stuff")
	worktree := Successful(repo.Worktree())

	Expect(copyFile("README", path.Join(tmpdir, "README"), fileMode)).To(Succeed())
	Expect(worktree.Add("README")).Error().NotTo(HaveOccurred())
	commit := Successful(worktree.Commit("initial check-in", commitOptions))
	Expect(repo.CreateTag("v1.0", commit, nil)).Error().NotTo(HaveOccurred())

	Expect(os.Mkdir(path.Join(tmpdir, "fodder"), dirMode)).Error().NotTo(HaveOccurred())
	Expect(copyFile("fodder/empty", path.Join(tmpdir, "fodder/empty"), fileMode)).To(Succeed())
	Expect(worktree.Add("fodder")).Error().NotTo(HaveOccurred())
	Expect(os.Mkdir(path.Join(tmpdir, "folder"), dirMode)).Error().NotTo(HaveOccurred())
	Expect(os.Mkdir(path.Join(tmpdir, "folder/subfolder"), dirMode)).Error().NotTo(HaveOccurred())
	Expect(copyFile("folder/subfolder/canary.txt", path.Join(tmpdir, "folder/subfolder/canary.txt"), fileMode)).To(Succeed())
	Expect(copyFile("folder/subfolder/schkript.sh", path.Join(tmpdir, "folder/subfolder/schkript.sh"), exeMode)).To(Succeed())
	Expect(worktree.Add("folder")).Error().NotTo(HaveOccurred())
	commit = Successful(worktree.Commit("adds canary", commitOptions))
	Expect(repo.CreateTag("v1.1.1", commit, nil)).Error().NotTo(HaveOccurred())

	return tmpdir
}

func copyFile(from, to string, mode fs.FileMode) error {
	contents, err := contentfs.ReadFile(path.Join("files", from))
	if err != nil {
		return fmt.Errorf("cannot read test file from embedded file system, reason: %w", err)
	}
	if err := os.WriteFile(to, contents, mode); err != nil {
		return fmt.Errorf("cannot write test file, reason: %w", err)
	}
	return nil
}
