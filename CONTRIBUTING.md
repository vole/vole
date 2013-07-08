General
-------

We welcome contributions. For new features or changes to the API, please write up your proposal and make a Github issue *before* coding it. There will always be discussion around implementation and we want to make sure your time isn't wasted. You can also search the [mailing list](https://groups.google.com/forum/#!forum/vole-dev).

For obvious bugs, just send a pull request.

Building
--------

Compiling and building Vole requires Go version 1.1+. Download it from [http://golang.org/](http://golang.org/).

To grab the latest development version of Vole and run it:

    git clone https://github.com/vole/vole.git
    cd vole
    export GOPATH=`pwd`
    export PATH=$PATH:$GOPATH/bin
    go run src/vole/vole.go

In your browser, navigate to http://localhost:6789.

Windows
-------

To set the Go path on Windows, use:

    set GOPATH="c:\Users..."
    set PATH=%GOPATH%\bin;%PATH%

Testing
-------

To run tests, execute `go test -v lib/store` from the root of the Vole project. This will create a ~/VoleTest folder with test data. This can be safely deleted once the tests have run.
