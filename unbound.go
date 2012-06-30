package unbound

/*
#cgo LDFLAGS: -lunbound
#include <stdlib.h>
#include <unbound.h>

typedef struct ub_ctx ctx;

struct ub_result *new_ub_result() {
	struct ub_result *r;
	r = calloc(sizeof(struct ub_result), 1);
	return r;
}

*/
import "C"

import "dns"

type Unbound struct {
	ctx *C.ctx
}

type Result struct {
	Qname  string   // text string, original question
	Qtype  uint16      // type code asked for
	Qclass uint16      // class code asked for 
	Data   []string // array of rdata items, NULL terminated
	// Not needed                 int* len;    /* array with lengths of rdata items
	CanonName    string   // canonical name of result
	Rcode        int      // additional error code in case of no data
	AnswerPacket *dns.Msg // full network format answer packet
	// Not needed                 int answer_len; // length of packet in octets
	HaveData bool   // true if there is data
	NxDomain bool   // true if nodata because name does not exist
	Secure   bool   // true if result is secure
	Bogus    bool   // true if a security failure happened
	WhyBogus string // string with error if bogus
}

type UnboundError struct {
	Err  string
	Code int // Internal unbound error code
}

func (e *UnboundError) Error() string {
	// Also Code here?
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

func (u *Unbound) ResolvConf(fname string) error {
	i := C.ub_ctx_resolvconf(u.ctx, C.CString(fname))
	return newError(int(i))
}

func (u *Unbound) Resolve(name string, rrtype, rrclass uint16) (*Result, error) {
	res := C.new_ub_result()
	r := new(Result)
	i := C.ub_resolve(u.ctx, C.CString(name), C.int(rrtype), C.int(rrclass), &res)
	err := newError(int(i))

	// Copy the data from res to Result for easy dismembering in Go
	r.Qname = C.GoString(res.qname)
	r.Qtype = uint16(res.qtype)		// Create full blown RR?
	r.Qclass = uint16(res.qclass)
	// res.data to ...			// And incorperate these...?
	r.CanonName = C.GoString(res.canonname)
	r.Rcode = int(res.rcode)
	r.AnswerPacket = new(dns.Msg)
	r.AnswerPacket.Unpack(C.GoBytes(res.answer_packet, res.answer_len))	// return code??
	r.HaveData = res.havedata == 1
	r.NxDomain = res.nxdomain == 1
	r.Secure = res.secure == 1
	r.Bogus = res.bogus == 1
	r.WhyBogus = C.GoString(res.why_bogus)

	C.ub_resolve_free(res)

	return r, err
}
