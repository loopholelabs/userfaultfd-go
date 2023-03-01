# userfaultfd-go

Go bindings to userfaultfd.

![Go Version](https://img.shields.io/badge/go%20version-%3E=1.19-61CFDD.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/loopholelabs/userfaultfd-go.svg)](https://pkg.go.dev/github.com/loopholelabs/userfaultfd-go)

## Overview

`userfaultfd-go` provides Go bindings to [userfaultfd](https://man7.org/linux/man-pages/man2/userfaultfd.2.html), which allows for handling page faults in a memory region in userspace.

It enables you to...

- **Expose any `io.ReaderAt` as a Go slice**: You can use this feature to work with any `io.ReaderAt` interface as a Go slice, which can be useful when working with external libraries that can't work with Go buffers or readers.
- **Access remote files as a slice without fetching**: With `userfaultfd-go`, you can access a remote file stored in S3 or a file as a slice without needing to fetch its complete contents locally, which can be beneficial when working with large files.
- **Track changes to `mmap`ed regions**: You can track changes made to a `mmap`ed region by one process from another process.

## Installation

You can add `userfaultfd-go` to your Go project by running the following:

```shell
$ go get github.com/loopholelabs/userfaultfd/...@latest
```

## Reference

To make getting started with `userfaultfd-go` easier, take a look at the following examples:

- [I/O benchmark](./cmd/userfaultfd-go-benchmark/main.go)
- [Exposing a pattern](./cmd/userfaultfd-go-example-abc/main.go)
- [Exposing a file](./cmd/userfaultfd-go-example-file/main.go)
- [Exposing an object in S3](./cmd/userfaultfd-go-example-s3/main.go)

## Acknowledgements

- [bytecodealliance/userfaultfd-rs](https://github.com/bytecodealliance/userfaultfd-rs) inspired the API design.

## Contributing

Bug reports and pull requests are welcome on GitHub at [https://github.com/loopholelabs/userfaultfd-go][gitrepo]. For more contribution information check out [the contribution guide](./CONTRIBUTING.md).

## License

The `userfaultfd-go` project is available as open source under the terms of the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).

## Code of Conduct

Everyone interacting in the `userfaultfd-go` project's codebases, issue trackers, chat rooms and mailing lists is expected to follow the [CNCF Code of Conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md).

## Project Managed By:

[![https://loopholelabs.io][loopholelabs]](https://loopholelabs.io)

[gitrepo]: https://github.com/loopholelabs/userfaultfd-go
[loopholelabs]: https://cdn.loopholelabs.io/loopholelabs/LoopholeLabsLogo.svg
