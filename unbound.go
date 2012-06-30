package unbound

/*
#cgo LDFLAGS: -lunbound
#include <unbound.h>

typedef struct ub_ctx ub;

*/
import "C"

type Unbound struct {
	ctx *C.ub
}

func New() *Unbound {
	u := new(Unbound)
	u.ctx = C.ub_ctx_create()
	return u
}
