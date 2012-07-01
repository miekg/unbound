// Package unbound implements a wrapper for libunbound(3).
// Unbound is a DNSSEC aware resolver, see https://unbound.net/
// for more information.
//
// The method's documentation can be found in libunbound(3).
// The names of the methods are in sync with the
// names used in unbound, but the underscores are removed and they
// are in camel-case, e.g. ub_ctx_resolv_conf becomes u.ResolvConf.
// Except for ub_ctx_create() and ub_ctx_delete(),
// which beome: New() and Destroy().
//
// Basic use pattern:
//	u := unbound.New()
//	defer u.Destroy()
//	err := u.ResolvConf("/etc/resolv.conf")
package unbound

/*
#cgo LDFLAGS: -lunbound
#include <stdlib.h>
#include <unbound.h>

typedef struct ub_ctx ctx;

int array_len(int *l)			{ return sizeof(l)/sizeof(int); }
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
	"os"
	"unsafe"
)

type Unbound struct {
	ctx *C.ctx
}

// Results is Unbound's ub_result adapted for Go.
type Result struct {
	Qname        string   // Text string, original question
	Qtype        uint16   // Type code asked for
	Qclass       uint16   // Class code asked for 
	Data         [][]byte // Slice of rdata items,
	CanonName    string   // Canonical name of result
	Rcode        int      // Additional error code in case of no data
	AnswerPacket *dns.Msg // Full network format answer packet
	HaveData     bool     // True if there is data
	NxDomain     bool     // True if nodata because name does not exist
	Secure       bool     // True if result is secure
	Bogus        bool     // True if a security failure happened
	WhyBogus     string   // String with error if bogus
}

// UnboundError is an error returned from Unbound, it wraps both the
// return code and the error string as returned by ub_strerror.
type UnboundError struct {
	Err  string
	code int
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
	e.code = i
	return e
}

func errorString(i int) string {
	return C.GoString(C.ub_strerror(C.int(i)))
}

// New wraps Unbound's ub_ctx_create.
func New() *Unbound {
	u := new(Unbound)
	u.ctx = C.ub_ctx_create()
	return u
}

// Destroy wraps Unbound's ub_ctx_delete.
func (u *Unbound) Destroy() {
	C.ub_ctx_delete(u.ctx)
}

// ResolvConf wraps Unbound's ub_ctx_resolvconf.
func (u *Unbound) ResolvConf(fname string) error {
	cfname := C.CString(fname)
	defer C.free(unsafe.Pointer(cfname))
	i := C.ub_ctx_resolvconf(u.ctx, cfname)
	return newError(int(i))
}

// SetOption wraps Unbound's ub_ctx_set_option.
func (u *Unbound) SetOption(opt, val string) error {
	copt := C.CString(opt)
	defer C.free(unsafe.Pointer(copt))
	cval := C.CString(val)
	defer C.free(unsafe.Pointer(cval))
	i := C.ub_ctx_set_option(u.ctx, copt, cval)
	return newError(int(i))
}

// GetOption wraps Unbound's ub_ctx_get_option.
func (u *Unbound) GetOption(opt string) (string, error) {
	val := ""
	cval := C.CString(val)
	defer C.free(unsafe.Pointer(cval))
	i := C.ub_ctx_get_option(u.ctx, C.CString(opt), &cval)
	// Not sure if this works...?
	return val, newError(int(i))
}

// Config wraps Unbound's ub_ctx_config.
func (u *Unbound) Config(fname string) error {
	cfname := C.CString(fname)
	defer C.free(unsafe.Pointer(cfname))
	i := C.ub_ctx_config(u.ctx, cfname)
	return newError(int(i))
}

// SetFwd wraps Unbound's ub_ctx_set_fwd.
func (u *Unbound) SetFwd(addr string) error {
	caddr := C.CString(addr)
	defer C.free(unsafe.Pointer(caddr))
	i := C.ub_ctx_set_fwd(u.ctx, caddr)
	return newError(int(i))
}

// Hosts wraps Unbound's ub_ctx_hosts.
func (u *Unbound) Hosts(fname string) error {
	cfname := C.CString(fname)
	defer C.free(unsafe.Pointer(cfname))
	i := C.ub_ctx_hosts(u.ctx, cfname)
	return newError(int(i))
}

// Resolve wraps Unbound's ub_resolve.
func (u *Unbound) Resolve(name string, rrtype, rrclass uint16) (*Result, error) {
	res := C.new_ub_result()
	r := new(Result)
	defer C.ub_resolve_free(res)
	i := C.ub_resolve(u.ctx, C.CString(name), C.int(rrtype), C.int(rrclass), &res)
	err := newError(int(i))
	if err != nil {
		return nil, err
	}

	r.Qname = C.GoString(res.qname)
	r.Qtype = uint16(res.qtype)
	r.Qclass = uint16(res.qclass)
	r.Data = make([][]byte, 0)
	j := 0
	b := C.GoBytes(unsafe.Pointer(C.array_elem_char(res.data, C.int(j))), C.array_elem_int(res.len, C.int(j)))
	for len(b) != 0 {
		r.Data = append(r.Data, b)
		j++
		b = C.GoBytes(unsafe.Pointer(C.array_elem_char(res.data, C.int(j))), C.array_elem_int(res.len, C.int(j)))
	}
	// Try to create an RR

	r.CanonName = C.GoString(res.canonname)
	r.Rcode = int(res.rcode)
	r.AnswerPacket = new(dns.Msg)
	r.AnswerPacket.Unpack(C.GoBytes(res.answer_packet, res.answer_len)) // Should always work
	r.HaveData = res.havedata == 1
	r.NxDomain = res.nxdomain == 1
	r.Secure = res.secure == 1
	r.Bogus = res.bogus == 1
	r.WhyBogus = C.GoString(res.why_bogus)
	return r, err
}

// ResolveAsync does *not* wrap the Unbound function, instead
// it utilizes Go's goroutines to mimic the async behavoir Unbound
// implements. As a result the function signature is different.
// The function f is called after the resolution is finished.
// Also the ub_cancel, ub_wait_, ub_fd? are not implemented.
func (u *Unbound) ResolveAsync(name string, rrtype, rrclass uint16, m interface{}, f func(interface{}, error, *Result)) error {
	go func() {
		r, e := u.Resolve(name, rrtype, rrclass)
		f(m, e, r)
	}()
	return newError(0)
}

// AddTa wraps Unbound's ub_ctx_add_ta.
func (u *Unbound) AddTa(ta string) error {
	cta := C.CString(ta)
	i := C.ub_ctx_add_ta(u.ctx, cta)
	return newError(int(i))
}

// AddTaFile wraps Unbound's ub_ctx_add_ta_file.
func (u *Unbound) AddTaFile(fname string) error {
	cfname := C.CString(fname)
	defer C.free(unsafe.Pointer(cfname))
	i := C.ub_ctx_add_ta_file(u.ctx, cfname)
	return newError(int(i))
}

// AddTaFile wraps Unbound's ub_ctx_trustedkeys.
func (u *Unbound) TrustedKeys(fname string) error {
	cfname := C.CString(fname)
	defer C.free(unsafe.Pointer(cfname))
	i := C.ub_ctx_trustedkeys(u.ctx, cfname)
	return newError(int(i))
}

// ZoneAdd wraps Unbound's ub_ctx_zone_add.
func (u *Unbound) ZoneAdd(zone_name, zone_type string) error {
	czone_name := C.CString(zone_name)
	defer C.free(unsafe.Pointer(czone_name))
	czone_type := C.CString(zone_type)
	defer C.free(unsafe.Pointer(czone_type))
	i := C.ub_ctx_zone_add(u.ctx, czone_name, czone_type)
	return newError(int(i))
}

// ZoneRemove wraps Unbound's ub_ctx_zone_remove.
func (u *Unbound) ZoneRemove(zone_name string) error {
	czone_name := C.CString(zone_name)
	defer C.free(unsafe.Pointer(czone_name))
	i := C.ub_ctx_zone_remove(u.ctx, czone_name)
	return newError(int(i))
}

// DataAdd wraps Unbound's ub_ctx_data_add.
func (u *Unbound) DataAdd(data string) error {
	cdata := C.CString(data)
	defer C.free(unsafe.Pointer(cdata))
	i := C.ub_ctx_data_add(u.ctx, cdata)
	return newError(int(i))
}

// DataRemove wraps Unbound's ub_ctx_data_remove.
func (u *Unbound) DataRemove(data string) error {
	cdata := C.CString(data)
	defer C.free(unsafe.Pointer(cdata))
	i := C.ub_ctx_data_remove(u.ctx, cdata)
	return newError(int(i))
}

// DebugOut wraps Unbound's ub_ctx_debugout
func (u *Unbound) DebugOut(out *os.File) error {
	// TODO(mg): How to converto os.File to *FILE?
	i := 0
	return newError(int(i))
}

// DataRemove wraps Unbound's ub_ctx_data_level
func (u *Unbound) DebugLevel(d int) error {
	i := C.ub_ctx_debuglevel(u.ctx, C.int(d))
	return newError(int(i))
}
