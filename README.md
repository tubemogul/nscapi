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

