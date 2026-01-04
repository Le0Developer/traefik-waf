package internal

import (
	"fmt"
	"net/http"

	"github.com/corazawaf/coraza/v3/types"
)

func (i *Instance) evaluateRules(r http.Request) *types.Interruption {
	if i.engine == nil {
		return nil
	}

	tx := i.engine.NewTransaction()
	//nolint:errcheck
	defer tx.Close()
	defer tx.ProcessLogging()

	// We use sample data because getting the real IP is difficult
	// due to proxy nesting
	// TODO: ^ fix above
	tx.ProcessConnection("127.0.0.1", 1337, "127.0.0.1", 80)

	tx.ProcessURI(r.URL.String(), r.Method, r.Proto)

	for name, values := range r.Header {
		for _, value := range values {
			tx.AddRequestHeader(name, value)
		}
	}

	// Phase 1 done (request headers)
	if it := tx.ProcessRequestHeaders(); it != nil {
		return it
	}

	// Phase 2 (request body)
	if it, _, err := tx.ReadRequestBodyFrom(r.Body); it != nil || err != nil {
		if it != nil {
			return it
		}
		fmt.Printf("error reading request body: %v\n", err)
	}

	// We can't process response

	return nil
}
