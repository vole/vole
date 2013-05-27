Building
--------

    git clone https://github.com/vole/vole.git
    cd vole
    export GOPATH=`pwd`
    export PATH=$PATH:$GOPATH/bin
    make

In your browser, navigate to http://localhost:6789



Testing
-------

To run tests, execute `go test -v lib/store` from the root of the Vole project. This will create a ~/VoleTest directory with test data. This can be safely deleted once the tests have run.

Go cheat sheet
--------------

| Command | Description |
| ------- | ----------- |
fmt.Println(name) | Print a string
fmt.Printf("%+v", user) | Print a struct
