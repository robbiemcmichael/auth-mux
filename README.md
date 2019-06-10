# auth-mux

`auth-mux` is an authentication multiplexer, its job is to accept many
different authentication schemes and return a validation result which can
optionally be used to issue new credentials using a different authentication
scheme. It relies on being able to take an authentication request and return
either an error message or a standard set of
[identity claims](./internal/types/types.go).

## Use cases

- Transforming credentials from `m` issuers to a common authentication scheme
- Transforming credentials from some authentication scheme to the `n` different
  schemes supported by the systems that you run
- Supporting any of the `m x n` combinations above
- Securerly reissuing credentials while ensuring that the original `m` issuers
  cannot issue conflicting credentials
- Forking this repo and implementing a single interface to allow a non-standard
  authentication scheme to place nicely with the other systems that you run

## Design

Authentication schemes are supported by implementing the `Input` or `Output`
interface. An `Input` must implement a handler that takes an HTTP request and
returns the `Validation` type which represents a validation result. An `Output`
must implement a handler that takes this `Validation` type and returns an HTTP
response.

These standard interfaces allow any `Input` and `Output` combination to be
composed. The HTTP server automatically creates a handler for each combination
with path `/:input/:output`. This allows the caller to use the HTTP path to
select the authentication schemes, forming a multiplexer for the `Input` and a
demultiplexer for the `Output`.

<p align="center"><img src="docs/mux-demux.svg"></p>
