/*
Package version implements finding the latest and greatest tagged version in a
remote git repository.

# Usage

In its simplest form, the following snippet finds the latest version for tags in
semver format (MAJOR.MINOR.PATH, with MINOR and PATH being optional) with or
without a "v" prefix. The latest semver found is returned, as well as its
corresponding tag reference for further use.

	remoteURL := "..." // ...could even be just a local path.
	semver, ref, err := LatestReleaseTag(
	    context.Background(),
	    remoteURL,
	    SemverTagMatcher(""))

In case the version tags feature a prefix, such as in "libfoo-1.6.66", use
[NewPrefixedTagMatcher]:

	semver, ref, err := LatestReleaseTag(
		context.Background(),
		remoteURL,
		NewPrefixedTagMatcher("libfoo-"))
*/
package version
