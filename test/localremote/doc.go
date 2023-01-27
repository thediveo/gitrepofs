/*
Package localremote aids (go-)git-related unit tests by managing temporary git
repositories in the file system. Now, go-git can perfectly access remote
repositories which are actually local (using the "file:" protocol), hence the
package name "localremote".

This package leverages [Ginko] testing framework and [Gomega] matcher library
for simplifying the code and at the same time making it expressive, without the
overly load noise of Go's (barely existing) assertion checking and error
handling.

[Ginko]: https://github.com/onsi/ginkgo
[Gomega]: https://github.com/onsi/gomega
*/
package localremote
