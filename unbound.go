package unbound

/*
#cgo LDFLAGS: -lunbound
#include <stdlib.h>
#include <unbound.h>

typedef struct ub_ctx ctx;

int array_len(int *l) { return sizeof(l)/sizeof(int); }
int array_elem_int(int *l, int i)       { return l[i]; }
char * array_elem_char(char **l, int i) { return l[i]; }

struct ub_result *new_ub_result() {
	struct ub_result *r;
	r = calloc(sizeof(struct ub_result), 1);
	return r;
}
*/
import "C"

import (
	"github.com/miekg/dns"
	"unsafe"
)

type Unbound struct {
	ctx *C.ctx
}

// ub_result adapted to Go.
type Result struct {
	Qname        string   // text string, original question
	Qtype        uint16   // type code asked for
	Qclass       uint16   // class code asked for 
	Data         [][]byte // slice of rdata items,
	CanonName    string   // canonical name of result
	Rcode        int      // additional error code in case of no data
	AnswerPacket *dns.Msg // full network format answer packet
	HaveData     bool     // true if there is data
	NxDomain     bool     // true if nodata because name does not exist
	Secure       bool     // true if result is secure
	Bogus        bool     // true if a security failure happened
	WhyBogus     string   // string with error if bogus
}

type UnboundError struct {
	Err  string
	Code int // Internal unbound error code
}

func (e *UnboundError) Error() string {
	return e.Err
}

func newError(i int) error {
	if i == 0 {
		return nil
	}
	e := new(UnboundError)
	e.Err = errorString(i)
	e.Code = i
	return e
}

func errorString(i int) string {
	return C.GoString(C.ub_strerror(C.int(i)))
}

func New() *Unbound {
	u := new(Unbound)
	u.ctx = C.ub_ctx_create()
	return u
}

func (u *Unbound) Delete() {
	C.ub_ctx_delete(u.ctx)
}

func (u *Unbound) ResolvConf(fname string) error {
	i := C.ub_ctx_resolvconf(u.ctx, C.CString(fname))
	return newError(int(i))
}

func (u *Unbound) Hosts(fname string) error {
	i := C.ub_ctx_hosts(u.ctx, C.CString(fname))
	return newError(int(i))
}

func (u *Unbound) Resolve(name string, rrtype, rrclass uint16) (*Result, error) {
	res := C.new_ub_result()
	r := new(Result)
	i := C.ub_resolve(u.ctx, C.CString(name), C.int(rrtype), C.int(rrclass), &res)
	err := newError(int(i))
	if err != nil {
		return nil, err
	}

	r.Qname = C.GoString(res.qname)
	r.Qtype = uint16(res.qtype)
	r.Qclass = uint16(res.qclass)
	r.Data = make([][]byte, 0)
	for i := 0; i < int(C.array_len(res.len))-1; i++ {
		r.Data = append(r.Data,
			C.GoBytes(
				unsafe.Pointer(C.array_elem_char(res.data, C.int(i))),
				C.array_elem_int(res.len, C.int(i))))
	}
	r.CanonName = C.GoString(res.canonname)
	r.Rcode = int(res.rcode)
	r.AnswerPacket = new(dns.Msg)
	r.AnswerPacket.Unpack(C.GoBytes(res.answer_packet, res.answer_len)) // TODO(mg): return code
	r.HaveData = res.havedata == 1
	r.NxDomain = res.nxdomain == 1
	r.Secure = res.secure == 1
	r.Bogus = res.bogus == 1
	r.WhyBogus = C.GoString(res.why_bogus)

	C.ub_resolve_free(res)
	return r, err
}

// AddTa wraps Unbound's ub_ctx_add_ta.
func (u *Unbound) AddTa(ta string) error {
	i := C.ub_ctx_add_ta(u.ctx, C.CString(ta))
	return newError(int(i))
}

// AddTaFile wraps Unbound's ub_ctx_add_ta_file.
func (u *Unbound) AddTaFile(fname string) error {
	i := C.ub_ctx_add_ta_file(u.ctx, C.CString(fname))
	return newError(int(i))
}

// AddTaFile wraps Unbound's ub_ctx_trustedkeys.
func (u *Unbound) TrustedKeys(fname string) error {
	i := C.ub_ctx_trustedkeys(u.ctx, C.CString(fname))
	return newError(int(i))
}

// ZoneAdd wraps Unbound's ub_ctx_zone_add.
func (u *Unbound) ZoneAdd(zone_name, zone_type string) error {
	i := C.ub_ctx_zone_add(u.ctx, C.CString(zone_name), C.CString(zone_type))
	return newError(int(i))
}

// ZoneRemove wraps Unbound's ub_ctx_zone_remove.
func (u *Unbound) ZoneRemove(zone_name string) error {
	i := C.ub_ctx_zone_remove(u.ctx, C.CString(zone_name))
	return newError(int(i))
}

// DataAdd wraps Unbound's ub_ctx_data_add.
func (u *Unbound) DataAdd(data string) error {
	i := C.ub_ctx_data_add(u.ctx, C.CString(data))
	return newError(int(i))
}

// DataRemove wraps Unbound's ub_ctx_data_remove.
func (u *Unbound) DataRemove(data string) error {
	i := C.ub_ctx_data_remove(u.ctx, C.CString(data))
	return newError(int(i))
}
