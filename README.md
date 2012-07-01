# Unbound

A wrapper for Unbound in Go.

Unbound's asynchronous behavior is mimicked by using goroutines, *not* by
calling `ub_resolve_async`.

The website for Unbound is https://unbound.net/, were you can find further documentation.
