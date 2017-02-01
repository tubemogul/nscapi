# NSCAPI: a server that ingests nsca calls and exposes the result via a json API

[![TravisBuild](https://travis-ci.org/tubemogul/nscapi.svg?branch=master)](https://travis-ci.org/tubemogul/nscapi)

## Introduction

Nscapi is a go application meant to help moving from a nsca model to a pull
model by providing a JSON endpoint to summarize the state of the nsca calls it
received. 

The tool is split into 3 main parts:

* Server that receive nsca calls and put them in a non-locking queue
* Worker that reads from the non-locking queue and put the data in a cache
* HTTP server that will read from the cache and display the check results 
  in a json format in response to a call on `/api/reports` based on a template.
  The rest of the calls will display a usage summary based on another template

## Custom fields

A custom field is a key-value couple that is not contained in the nsca check
itself but will be returned by the API based on the values of the check itself.

The custom fields are defined in a hiera-style hierarchy.

Provided that:
* `%{hostgroup}` is the hostname received in the nsca check result minus any trailing number
* `%{check name}` is the name of the check as received in the nsca check result

The level of importance is defined as follows:
* `common.yaml` is where the default values are defined.
* `service/%{hostgroup}/all.yaml` is where you define the default values of your
  custom fields per host group. It has the priority over the custom fields
  already defined in `common.yaml` and can also add new fields of its own.
  If a host don't match any of the defined hostgroups, only the custom fields
  defined in `common.yaml` will be applied.
* `service/%{hostgroup}/%{check name}.yaml` has the highest priority. Any field value in
  this file will overwrite the values of the previously defined custom fields.

## Using the Makefile

The Makefile is used as a helper for building and testing the project. Current
capabilities:

 * `make lint`: runs golint on your files (requires `github.com/golang/lint/golint` installed)
 * `make fmt`: checks that the files are compliant with the gofmt format
 * `make vet`: runs `go tool vet` on your files to ensure there's no problems
 * `make test`: runs `make lint`, `make fmt`, `make vet` before running all the
   test, printing also the percentage of code coverage
 * `make race`: runs the tests with the `-race` option to detect race conditions
 * `make bench`: runs the benchmarks
 * `make gocov`: runs a gocov report (requires `github.com/axw/gocov/gocov`)
 * `make dep`: gets (with update option `-u`) the go dependencies required by
   the application
 * `make build`: runs `make dep` and `make test` before running a clean and
   build the binary
 * `make install`: runs `make dep` and `make test` before running a clean and
   install the application
 * `make` / `make all`: run `make test` and `make race`

## Contributions

Contributions to this project are welcome, though please
[file an issue](https://github.com/tubemogul/nscapi/issues/new).
before starting work on anything major as someone else could already be working
on it.

Contributions that do not provide the corresponding tests will most likely not
be accepted (we can help building the tests if needed).

Contributions that do not pass the basic gofmt, vet and other basic checks
provided in the Makefile will not be accepted. It's just a question of trying to
keep a basic code standard. Thanks for your help! :)

