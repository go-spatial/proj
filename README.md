[![Build Status](https://travis-ci.org/go-spatial/proj.svg?branch=master)](https://travis-ci.org/go-spatial/proj)
[![Report Card](https://goreportcard.com/badge/github.com/go-spatial/proj)](https://goreportcard.com/badge/github.com/go-spatial/proj)
[![Coverage Status](https://coveralls.io/repos/github/go-spatial/proj/badge.svg?branch=master)](https://coveralls.io/github/go-spatial/proj?branch=master)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/go-spatial/proj)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/go-spatial/proj/blob/master/LICENSE.md)

# proj: PROJ4, for Go!

This project is **UNDER ACTIVE DEVELOPMENT** and is, therefore, **NOT STABLE**, which subsequently means it **SHOULD NOT BE USED FOR PRODUCTION PURPOSES**.

Contact `mpg@flaxen.com` if you'd like to help out on the project.


# Guiding Principles and Goals and Plans

As much as we all love PROJ4, we are not going to attempt to blindly (or even blithely) port every jot and tittle of the project.

In no particular order, these are the conditions we're imposing on ourselves:

* We are going to use the PROJ 5.0.1 release as our starting point.
* We will look to the `proj4js` project for suggestions as to what PROJ4 code does and does not need to be ported, and how.
* We will be "DONE" when the targetted test cases from PROJ are passing, including both direct invocations of the coordinate operations and proj-string tests.
* The coordinate operations are going to be ported pretty directly, but the function signatures and "catalog" will be idiomatic Go.
* The `proj` command-line app will be ported and will be idiomatic.
* All code will pass [the Go `metalinter`](https://github.com/alecthomas/gometalinter) cleanly all the time.
* Unit tests will be implemented for pretty much everything, using the "side-by-side" `_test` package style.
* Tests will be extracted from the various PROJ tests into idiomatic Go test cases. We will not port the new PROJ test harness.
* Go-style source code documentation will be provided.
* A set of small, clean usage examples will be provided.


# Project Layout

Packages (directories):
* `proj`: the `proj` command-line tool
* `operations`: the actual coordinate operations
* `support`: misc stuff, including the core API
* `examples`: simple yet instructive self-validating demos

The API:
* ...
