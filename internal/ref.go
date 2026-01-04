package internal

import (
	"fmt"
	"net/http"
)

var fallbackReference = "unknown"

func (i *Instance) reference(r *http.Request) string {
	if i.cfg.RefHeader != "" {
		if ref := r.Header.Get(i.cfg.RefHeader); ref != "" {
			return ref
		}
		return fallbackReference
	}

	ref := ""
	for _, v := range []string{
		"X-Request-ID",
		"CF-Ray",
		"CDN-Uid",
	} {
		if ref_ := r.Header.Get(v); ref_ != "" {
			if ref != "" {
				fmt.Printf("multiple reference values found: %s and %s. falling back", ref, ref_)
				return fallbackReference
			}
			ref = ref_
		}
	}

	if ref == "" {
		return fallbackReference
	}
	return ref
}
