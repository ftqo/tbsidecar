package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	tbt "github.com/tigerbeetle/tigerbeetle-go/pkg/types"

	"github.com/ftqo/tbsidecar"
)

type sidecarServer struct {
	tbClient tb.Client
}

var (
	clusterID      = flag.Int("clusterID", 0, "cluster ID")
	concurrencyMax = flag.Int("concurrencyMax", 32, "max concurrency")
	port           = flag.Int("port", 8081, "port")
)

func main() {
	flag.Parse()

	c, err := tb.NewClient(uint32(*clusterID), flag.Args(), uint(*concurrencyMax))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create new tigerbeetle client: %v", err)
		return
	}

	s := sidecarServer{tbClient: c}
	r := chi.NewRouter()

	r.Post("/accounts", s.createAccounts)
	r.Get("/accounts/{accountID}", s.lookupAccount)
	r.Post("/transfers", s.createTransfers)
	r.Get("/transfers/{transferID}", s.lookupTransfer)

	http.ListenAndServe(fmt.Sprintf(":%d", *port), r)
}

func (s sidecarServer) createAccounts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var accounts []tbsidecar.Account
	err := json.NewDecoder(r.Body).Decode(&accounts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode create_accounts json body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tbAccounts := make([]tbt.Account, 0, len(accounts))
	for _, a := range accounts {
		tbAccounts = append(tbAccounts, tbt.Account{
			ID:             a.ID,
			UserData:       a.UserData,
			Reserved:       a.Reserved,
			Ledger:         a.Ledger,
			Code:           a.Code,
			Flags:          a.Flags,
			DebitsPending:  a.DebitsPending,
			DebitsPosted:   a.DebitsPosted,
			CreditsPending: a.CreditsPending,
			CreditsPosted:  a.CreditsPosted,
			Timestamp:      a.Timestamp,
		})
	}

	tbResults, err := s.tbClient.CreateAccounts(tbAccounts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to make create_accounts request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	results := make([]tbsidecar.AccountEventResult, 0, len(tbResults))
	for _, tbResult := range tbResults {
		results = append(results, tbsidecar.AccountEventResult{
			Index:  tbResult.Index,
			Result: tbsidecar.CreateAccountResult(tbResult.Result),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(results)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode create_accounts results: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(buf.Bytes())
}

func (s sidecarServer) createTransfers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var transfers []tbsidecar.Transfer
	err := json.NewDecoder(r.Body).Decode(&transfers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode create_transfers json body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tbTransfers := make([]tbt.Transfer, 0, len(transfers))
	for _, t := range transfers {
		tbTransfers = append(tbTransfers, tbt.Transfer{
			ID:              t.ID,
			DebitAccountID:  t.DebitAccountID,
			CreditAccountID: t.CreditAccountID,
			UserData:        t.UserData,
			Reserved:        t.Reserved,
			PendingID:       t.PendingID,
			Timeout:         t.Timeout,
			Ledger:          t.Ledger,
			Code:            t.Code,
			Flags:           t.Flags,
			Amount:          t.Amount,
			Timestamp:       t.Timestamp,
		})
	}

	tbResults, err := s.tbClient.CreateTransfers(tbTransfers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to make create_transfers request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	results := make([]tbsidecar.TransferEventResult, 0, len(tbResults))
	for _, tbResult := range tbResults {
		results = append(results, tbsidecar.TransferEventResult{
			Index:  tbResult.Index,
			Result: tbsidecar.CreateTransferResult(tbResult.Result),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(results)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode create_transfers results: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(buf.Bytes())
}

func (s sidecarServer) lookupAccount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	accountIDstr := chi.URLParam(r, "accountID")
	accountID := tbsidecar.StringToBytes(accountIDstr)
	tbAccountID := tbt.BytesToUint128(accountID)

	results, err := s.tbClient.LookupAccounts([]tbt.Uint128{tbAccountID})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to make lookup_accounts request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(results[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode lookup_accounts results: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(buf.Bytes())
}

func (s sidecarServer) lookupTransfer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	transferIDstr := chi.URLParam(r, "accountID")
	transferID := tbsidecar.StringToBytes(transferIDstr)
	tbTransferID := tbt.BytesToUint128(transferID)

	results, err := s.tbClient.LookupTransfers([]tbt.Uint128{tbTransferID})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to make lookup_account request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(results[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode lookup_transfer results: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(buf.Bytes())
}
