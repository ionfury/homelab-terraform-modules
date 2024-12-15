# Automated Tests

This folder contains automated tests for the modules in the `modules` folder.  The tests are written in `go` and use a test library called [terratest](https://terratest.gruntwork.io/).

## Running the Tests

> **WARNING**: These tests are deploying real resources.

### Pre-requisites

```sh
> brew bundle
```

### Run all tests

```sh
go test -v -timeout 10m
```

### Run a specific test

```sh
go test -v -timeout 10m -run '<TEST_NAME>'
```
