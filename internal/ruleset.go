package internal

import (
	"fmt"
	"net/http"

	"github.com/corazawaf/coraza/v3/collection"
	"github.com/corazawaf/coraza/v3/experimental/plugins/plugintypes"
	"github.com/corazawaf/coraza/v3/types"
	"github.com/corazawaf/coraza/v3/types/variables"
)

func (i *Instance) evaluateRules(r *http.Request) *types.Interruption {
	if i.engine == nil {
		return nil
	}

	tx := i.engine.NewTransaction()
	defer func() {
		if i.cfg.Verbosity >= 8 {
			txState := tx.(plugintypes.TransactionState)
			txState.Variables().All(func(_ variables.RuleVariable, v collection.Collection) bool {
				fmt.Println("--- Variable:", v.Name())
				for _, e := range v.FindAll() {
					fmt.Printf("    %s: %q\n", e.Key(), e.Value())
				}
				return true
			})

			for _, rule := range tx.MatchedRules() {
				fmt.Printf("+++ Matched rule: ID=%d Msg=%q Log=%q\n", rule.Rule().ID(), rule.Message(), rule.Data())
			}
		}

		tx.ProcessLogging()
		_ = tx.Close()
	}()

	// we cant get the real ports from http.Request
	tx.ProcessConnection(i.getRemoteIP(r), 0, "", 0)

	tx.ProcessURI(r.URL.String(), r.Method, r.Proto)

	for name, values := range r.Header {
		for _, value := range values {
			tx.AddRequestHeader(name, value)
		}
	}

	host := r.Header.Get("x-forwarded-host")
	tx.AddRequestHeader("Host", host)
	tx.SetServerName(host)

	if r.TransferEncoding != nil {
		tx.AddRequestHeader("Transfer-Encoding", r.TransferEncoding[0])
	}

	// Phase 1 done (request headers)
	if it := tx.ProcessRequestHeaders(); it != nil {
		return it
	}

	if tx.IsRequestBodyAccessible() && r.Body != nil && r.Body != http.NoBody {
		// Phase 2 (request body)
		if it, _, err := tx.ReadRequestBodyFrom(r.Body); it != nil || err != nil {
			if it != nil {
				return it
			}
			fmt.Printf("error reading request body: %v\n", err)
		}
	}

	if it, err := tx.ProcessRequestBody(); it != nil || err != nil {
		return it
	}

	// We can't process response

	return nil
}
