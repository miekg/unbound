# unbound

## Name

*unbound* - resolve names using libunbound.

## Description

With *unbound* you can (recursively) resolve names using the Unbound resolver from [NLnet
Labs](https://nlnetlabs.nl).

Unbound's `ub_result` has been extended with an slice of dns.RRs, this alleviates the need to parse
`ub_result.data` yourself.

## Syntax

Just enable the plugin with:

~~~ corefile
unbound
~~~

A wrapper for Unbound in Go.

## Notes

Compilation of this plugin requires CGO, which means the executables will use shared libraries
(OpenSSL, ldns and libunbound).

## See Also

The website for Unbound is https://unbound.net/, where you can find further documentation.
Tested/compiled to work for versions: 1.4.22 and 1.6.0-3+deb9u1 (Debian Stretch).

The tutorials found here are the originals ones adapted to Go.
