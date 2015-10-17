# Golang SLOC counting library

Forked from [bytebox](https://github.com/bytbox/sloc) and refactored into a library


## Commandline usages

```shell

Usage of sloc: sloc [options] [...files]
  -V  display version info and exit
  -cpuprofile string
      write cpu profile to file
  -json
      JSON-format output

```

## API Usage

```go

sloc := &SlocCounter{}

sloc.Add('/path/to/directory/containing/source-code')
sloc.Add('/path/to/file.go')

sloc.Sloc()


```