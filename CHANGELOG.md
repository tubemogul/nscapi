# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/) 
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]
### Added
- Server that receive nsca calls and put them in a non-locking queue
- Worker that reads from the non-locking queue and put the data in a cache
- HTTP server that will read from the cache and display the check results
  in a json format in response to a call on `/api/reports` based on a template.
  The rest of the calls will display a usage summary based on another template
- Support for custom fields
- Support configuration via both command-line and environment variables
