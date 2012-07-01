# Unbound

A wrapper for Unbound in Go.

As `cgo` does not support function callbacks (calling a Go function from within
the C library) I'm still pondering how to implement the `*_async` function
defined in `libunbound`.

The website for Unbound is https://unbound.net/, were you can find further documentation.
