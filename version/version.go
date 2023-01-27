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

package version

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"golang.org/x/mod/semver"
)

// VersionMatcherFn returns the semver information embedded in a refname. If the
// refname should not be taken into account, then any implementation of
// VersionMatcherFn should return an empty string instead.
type VersionMatcherFn func(refname string) (semver string)

// SemverTagMatcher is a version tag matcher that matches only on tags in semver
// version. That is, with an optional "v" prefix in the "MAJOR.MINOR.PATCH"
// format, where MINOR and PATCH are optional.
var SemverTagMatcher = NewPrefixedTagMatcher("")

// NewPrefixedTagMatcher returns a VersionMatcherFn to be used with
// [LatestReleaseTag]. The returned function only matches tags (/ref/tags/...)
// in the format <prefix><semver>. In particular, semvers must be in the format
// MAJOR, MAJOR.MINOR and MAJOR.MINOR.PATH with an optional "v" prefix, but no
// BUILD and PRERELEASE elements.
func NewPrefixedTagMatcher(prefix string) VersionMatcherFn {
	re := regexp.MustCompile(`(?m)^refs/tags/` + prefix + `((?:v)?\d(?:\.\d+(?:\.\d+)?)?)$`)
	return func(refname string) string {
		match := re.FindStringSubmatch(refname)
		if match == nil {
			return ""
		}
		semver := match[1] // 1st. group contains the semver string
		if !strings.HasPrefix(semver, "v") {
			semver = "v" + semver
		}
		return semver
	}
}

// LatestReleaseTag determines the latest release tag in the specified remote
// git repository that matches the specified pattern, especially when combined
// with [NewPrefixedTagMatcher].
func LatestReleaseTag(ctx context.Context, remoteURL string, fn VersionMatcherFn) (semanticver string, ref string, err error) {
	remote := git.NewRemote(
		memory.NewStorage(),
		&config.RemoteConfig{
			URLs: []string{remoteURL},
		})
	refs, err := remote.ListContext(ctx, &git.ListOptions{})
	if err != nil {
		return "", "", fmt.Errorf(
			"cannot list references in remote %q repository, reason: %w",
			remoteURL, err)
	}
	// an invalid semver which is automatically considered to be before any
	// valid semver.
	latest := ""
	latestref := ""
	for _, ref := range refs {
		if ref.Type() != plumbing.HashReference {
			continue
		}
		version := fn(ref.Name().String())
		if semver.Compare(version, latest) <= 0 {
			continue
		}
		latest = version
		latestref = ref.Name().String()
	}
	if latestref == "" {
		return "", "", fmt.Errorf(
			"no matching version reference in remote %q at all", remoteURL)
	}
	return latest, latestref, nil
}
