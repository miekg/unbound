package unbound

import (
	"dns"
)

// AddTaRR calls AddTa, but allows to directly use an dns.RR.
// This method is not found in Unbound.`
func (u *Unbound) AddTaRR(ta dns.RR) error { return u.AddTa(ta.String()) }

// DataAddRR calls DataAdd, but allows to directly use an dns.RR.
// This method is not found in Unbound.
func (u *Unbound) DataAddRR(data dns.RR) error { return u.DataAdd(data.String()) }

// DataRemoveRR calls DataRemove, but allows to directly use an dns.RR.
// This method is not found in Unbound.
func (u *Unbound) DataRemoveRR(data dns.RR) error { return u.DataRemove(data.String()) }
