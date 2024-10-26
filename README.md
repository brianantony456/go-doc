# Go docs

## About Go
- Statically & Strong Typed
- Compiled & Fast Compile Time
- Built in Concurrency
- Simplicity & Garbase collection

# Structure
- Modules - Collection of related packages
- Packages - Collection of go files

## Create a go module:
Create a module & run code

```bash
go mod init github.com/brianantony456/go-doc`       # Create module `<location>/<username>/<name_of_module>``
go build <file>.go                                  # Compile go code
go run <file>.go                                    # Compiles & Runs the code

go mod tidy                                         # Updates all the imported packages & creates hashes

go vendor                                           # Self contained local copy of dependencies
export GO11MODULE=on                                # Set to use the vendor
```


## Setup web app with gin framework
[Go Blueprint](https://go-blueprint.dev/)
