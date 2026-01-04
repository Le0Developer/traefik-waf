package internal

import (
	"fmt"
	"net/http"
)

func (i *Instance) reference(r *http.Request) string {
	if i.cfg.RefHeader != "" {
		return r.Header.Get(i.cfg.RefHeader)
	}

	ref := ""
	for _, v := range []string{
		"X-Request-ID",
		"CF-Ray",
		"CDN-RequestID",
	} {
		if ref_ := r.Header.Get(v); ref_ != "" {
			if ref != "" {
				if i.cfg.Verbosity >= 3 {
					fmt.Printf("multiple reference values found: %s and %s. falling back", ref, ref_)
				}

				return ""
			}
			ref = ref_
		}
	}

	return ref
}
