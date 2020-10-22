# go-oif
> Opinionated Imports Formatter

**go-oif** is an opinionated imports formatter which sorts imports into three
categories:

1. Built-in imports (`os`, `io`, `net/http`, etc)
2. Third-party imports (all your dependencies)
3. Project dependencies

## Install

Run the following to automatically download and install the latest version
to `$(go env GOPATH)/bin/go-oif`

```
curl -sSfL https://raw.githubusercontent.com/heyvito/go-oif/main/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

## Usage
This tool attempts to guess the project name (in order to detect local
dependencies) using data from `go.mod`. In case your project does not use
it, please suppply the project base name through the `--project-name` (`-n`)
flag. For instance:

### Using go.mod
Consider a project with the following go.mod file:
```go.mod
module github.com/heyvito/go-foo

go 1.15

// ...
```

And the following imports in an arbitrary file:

```go
package main

import "io"

import (
    "github.com/heyvito/go-foo/importc"
    "net"
    "github.com/heyvito/go-foo/importb"
    "github.com/heyvito/go-foo/importj"
    "os"
    ig "github.com/heyvito/go-foo/importg"
    ni "github.com/heyvito/go-foo/namedimport"
    _ "github.com/lib/side_effects_import"
    "github.com/foo/barA"
    "fmt"
    "github.com/heyvito/go-foo/importi"

    "github.com/heyvito/go-foo/importa"
    "github.com/foo/barC"
    "github.com/heyvito/go-foo/importf"
    "github.com/heyvito/go-foo/importh"
    "github.com/heyvito/go-foo/importk"
    "github.com/foo/barB"
    "github.com/heyvito/go-foo/importd"
    _ "github.com/heyvito/go-foo/side-effects-import"
    "context"

    annotatedImport "github.com/foo/barB"
    "github.com/foo/barD"
    "log"
    "github.com/heyvito/go-foo/importe"
)
```

Running `go-oif ./...` would rewrite the file import's as the following:

```go
package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "net"
    "os"

    "github.com/foo/barA"
    "github.com/foo/barB"
    "github.com/foo/barC"
    "github.com/foo/barD"
    "github.com/lib/side_effects_import"

    "github.com/heyvito/go-foo/importa"
    "github.com/heyvito/go-foo/importb"
    "github.com/heyvito/go-foo/importc"
    "github.com/heyvito/go-foo/importd"
    "github.com/heyvito/go-foo/importe"
    "github.com/heyvito/go-foo/importf"
    "github.com/heyvito/go-foo/importg"
    "github.com/heyvito/go-foo/importh"
    "github.com/heyvito/go-foo/importi"
    "github.com/heyvito/go-foo/importj"
    "github.com/heyvito/go-foo/importk"
    "github.com/heyvito/go-foo/namedimport"
    "github.com/heyvito/go-foo/side-effects-import"
)
```


### Without go.mod
The same effect can be achieved without a go.mod file:

```
$ go-oif -n github.com/heyvito/go-foo ./...
```

## License

```
MIT License

Copyright (c) 2020 Vito Sartori

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

```
