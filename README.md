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
