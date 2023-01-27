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
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/thediveo/gitrepofs/test/localremote"
	"github.com/thediveo/gitrepofs/version"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestComposerDecorator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gitloaderfs package")
}

var tmprepdir string
var repo *git.Repository
var commit *object.Commit

var _ = BeforeSuite(func(ctx context.Context) {
	By("creating a temporary test repository")
	tmprepdir = localremote.CreateTransientTestRepo()
	semver, ref, err := version.LatestReleaseTag(
		ctx, tmprepdir, version.SemverTagMatcher)
	Expect(err).NotTo(HaveOccurred())

	By(fmt.Sprintf("using version %s in %q", semver, tmprepdir))
	repo, err = git.CloneContext(ctx,
		memory.NewStorage(),
		nil,
		&git.CloneOptions{
			URL:           tmprepdir,
			SingleBranch:  true,
			ReferenceName: plumbing.ReferenceName(ref),
		})
	Expect(err).NotTo(HaveOccurred(), "cloning into memory failed")
	tag, err := repo.Tag(strings.Split(ref, "/")[2])
	Expect(err).NotTo(HaveOccurred(), "reference %s gone", ref)
	commit, err = repo.CommitObject(tag.Hash())
	Expect(err).NotTo(HaveOccurred())
})
