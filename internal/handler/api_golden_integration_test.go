//go:build integration

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/pangaunn/harness-engineering-workshop/internal/jobs"
	"github.com/pangaunn/harness-engineering-workshop/internal/model"
)

// TestAPI_Golden drives the API the way a client does and compares every
// response with testdata/golden. The golden files are the public API contract:
// clients depend on these fields and values. Numbers compare as exact decimals,
// so 3103 and 3103.00 are equal.
func TestAPI_Golden(t *testing.T) {
	s := newTestServer(t)
	job := &jobs.RevenueReport{Payments: s.payments, Reports: s.reports}
	book := func(roomType, checkIn, checkOut string) map[string]string {
		return map[string]string{"room_type_id": roomType, "guest_name": "Somchai", "check_in": checkIn, "check_out": checkOut}
	}
	report := func(date model.Date) {
		t.Helper()
		if _, err := job.Run(context.Background(), date); err != nil {
			t.Fatal(err)
		}
	}

	s.clock.t = time.Date(2026, 10, 14, 9, 30, 0, 0, model.HotelTime)
	checkGolden(t, s, "rooms", "GET", "/rooms", nil, http.StatusOK)
	checkGolden(t, s, "quote_deluxe", "POST", "/quotes", book("deluxe", "2026-10-15", "2026-10-18"), http.StatusOK)
	checkGolden(t, s, "book_standard", "POST", "/bookings", book("standard", "2026-10-19", "2026-10-21"), http.StatusCreated) // id 1
	checkGolden(t, s, "book_deluxe_1_night", "POST", "/bookings", book("deluxe", "2026-10-19", "2026-10-20"), http.StatusCreated)
	checkGolden(t, s, "book_deluxe_weekend", "POST", "/bookings", book("deluxe", "2026-10-15", "2026-10-18"), http.StatusCreated) // id 3
	checkGolden(t, s, "cancel_deluxe_weekend", "POST", "/bookings/3/cancel", nil, http.StatusOK)
	report("2026-10-14")
	checkGolden(t, s, "report_2026-10-14", "GET", "/reports/revenue?date=2026-10-14", nil, http.StatusOK)

	// Next day, refunds exceed charges.
	s.clock.t = time.Date(2026, 10, 15, 11, 0, 0, 0, model.HotelTime)
	checkGolden(t, s, "book_standard_weekend", "POST", "/bookings", book("standard", "2026-10-23", "2026-10-24"), http.StatusCreated)
	checkGolden(t, s, "cancel_standard", "POST", "/bookings/1/cancel", nil, http.StatusOK)
	report("2026-10-15")
	checkGolden(t, s, "report_2026-10-15", "GET", "/reports/revenue?date=2026-10-15", nil, http.StatusOK)
}

func checkGolden(t *testing.T, s *testServer, name, method, path string, body any, wantStatus int) {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatal(err)
		}
	}
	req, err := http.NewRequest(method, s.URL+path, &buf)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != wantStatus {
		t.Fatalf("%s: %s %s: status %d, want %d: %s", name, method, path, resp.StatusCode, wantStatus, got)
	}
	want, err := os.ReadFile(filepath.Join("testdata", "golden", name+".json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range diffJSON("", decodeJSON(t, want), decodeJSON(t, got)) {
		t.Errorf("%s: %s %s: %s", name, method, path, d)
	}
}

func decodeJSON(t *testing.T, data []byte) any {
	t.Helper()
	var v any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		t.Fatalf("decode %s: %v", data, err)
	}
	return v
}

// diffJSON lists every place where got differs from want.
func diffJSON(path string, want, got any) []string {
	switch w := want.(type) {
	case map[string]any:
		g, ok := got.(map[string]any)
		if !ok {
			return []string{fmt.Sprintf("%s: got %v, want an object", path, got)}
		}
		var diffs []string
		for _, k := range sortedKeys(w) {
			if _, ok := g[k]; !ok {
				diffs = append(diffs, fmt.Sprintf("%s.%s: missing", path, k))
				continue
			}
			diffs = append(diffs, diffJSON(path+"."+k, w[k], g[k])...)
		}
		for _, k := range sortedKeys(g) {
			if _, ok := w[k]; !ok {
				diffs = append(diffs, fmt.Sprintf("%s.%s: unexpected field", path, k))
			}
		}
		return diffs
	case []any:
		g, ok := got.([]any)
		if !ok || len(g) != len(w) {
			return []string{fmt.Sprintf("%s: got %v, want %v", path, got, want)}
		}
		var diffs []string
		for i := range w {
			diffs = append(diffs, diffJSON(fmt.Sprintf("%s[%d]", path, i), w[i], g[i])...)
		}
		return diffs
	case json.Number:
		g, ok := got.(json.Number)
		if !ok || !sameNumber(g, string(w)) {
			return []string{fmt.Sprintf("%s: got %v, want %s", path, got, w)}
		}
		return nil
	default:
		if got != want {
			return []string{fmt.Sprintf("%s: got %v, want %v", path, got, want)}
		}
		return nil
	}
}

func sameNumber(got json.Number, want string) bool {
	g, ok1 := new(big.Rat).SetString(string(got))
	w, ok2 := new(big.Rat).SetString(want)
	return ok1 && ok2 && g.Cmp(w) == 0
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
