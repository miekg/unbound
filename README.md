# Unbound

A wrapper for Unbound in Go.

Unbound's `ub_result` has been extended with an slice of dns.RRs, this alleviates
the need to parse `ub_result.data` yourself.

Unbound's asynchronous behavior is mimicked by using goroutines, *not* by
calling `ub_resolve_async`.

The website for Unbound is https://unbound.net/, were you can find further documentation.

# TODO
copy the sort SRV and sort MX from the base library.
