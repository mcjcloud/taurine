# Taurine

This is a hobby programming language, fueled by caffeine o.O

To view a list of language features, check out [the spec](https://github.com/mcjcloud/taurine/blob/master/docs/spec.md).

## Demo

Follow these steps to clone and build taurine (requires Go to be installed).

1. Run `git clone https://github.com/mcjcloud/taurine.git`
2. From the src directory, run `go build`
3. Run `./taurine ../example/num_guesser.tc`. See other example programs in the example directory.
4. Use the `--ast` flag before the filename to print the Abstract Syntax Tree in JSON format.
5. Use the `--print-tokens` flag to print the source files' tokens and their indecies.

## Install taurine

You can also run `go install` to install taurine to your `GOBIN`

## Tests

Tests are located in the `test` directory. Run `go run test/test.go` to execute tests.

Each directory in `test` contains a `src.tc`, `input.txt`, `output.txt`, and `ast.json` file. `test.go` works by running `src.tc` 
and comparing the AST and output with those specified in the respective files. Any program input needed should be included in `input.txt`.

