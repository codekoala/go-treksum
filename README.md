# Treksum

[![Travis CI Status](https://travis-ci.org/codekoala/go-treksum.svg?branch=master)](https://travis-ci.org/codekoala/go-treksum)
[![License BSD3](https://img.shields.io/badge/license-BSD3-blue.svg)](https://raw.githubusercontent.com/codekoala/go-treksum/master/LICENSE)
[![Downloads](https://img.shields.io/github/downloads/codekoala/go-treksum/total.svg)](https://github.com/codekoala/go-treksum/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/codekoala/treksum.svg?label=docker+pulls)](https://hub.docker.com/r/codekoala/treksum/)

Treksum provides access to the transcripts for all episodes of the following TV
series:

* Star Trek
* Star Trek: The Next Generation
* Star Trek: Deep Space Nine
* Star Trek: Voyager
* Star Trek: Enterprise

## What's Included

There are currently two utilities contained in this repository:

* treksum-scraper: scrapes Star Trek transcripts and shoves everything into a
  PostgreSQL database.
* treksum-api: a simple API to serve up random quotes from the populated
  PostgreSQL database. This utility is useless without the PostgreSQL database.

## Usage

The following environment variables are used to configure both utilities:

* ``TREKSUM_DBHOST``: IP address for PostgreSQL service. Default: ``localhost``
* ``TREKSUM_DBPORT``: Port on which PostgreSQL is listening. Default: ``5432``
* ``TREKSUM_DBNAME``: Name of PostgreSQL database where transcripts reside.
  Default: ``treksum``
* ``TREKSUM_DBUSER``: PostgreSQL username. Default: ``treksum``
* ``TREKSUM_DBPASSWORD``: PostgreSQL password. No default.

The following environment variable is used only for ``treksum-api``:

* ``TREKSUM_APIADDR``: interface and port to bind when serving requests.
  Default: ``:1323``

## Credits

The transcripts used by treksum were scraped from http://chakoteya.net/StarTrek/

## Legal

Star Trek and related marks are trademarks of CBS Studios Inc. This project is
purely for educational and entertainment purposes only. All other copyrights
remain property of their respective owners.
