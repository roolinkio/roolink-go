# Contributing

This repository is a copy of a package that Roolink develops in a private repository. Every commit here is exported from there, so nothing is pushed here directly.

- **Issues:** open them here.
- **Pull requests:** welcome. Open one here as usual, and CI checks it. A maintainer imports it into the private repository, where it's reviewed and merged, and the next export brings it back here with you as the commit's author. The pull request here is then closed.
- **Dependency updates:** these come from the private repository, so please don't open pull requests that only bump versions.

Build and test it with Go's own tools:

```sh
go test ./...
```
