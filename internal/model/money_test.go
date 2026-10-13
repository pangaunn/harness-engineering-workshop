package model

import (
	"encoding/json"
	"testing"
)

func TestMoney_JSON(t *testing.T) {
	cases := map[Money]string{0: "0.00", 5: "0.05", 123450: "1234.50", -2550: "-25.50"}
	for m, want := range cases {
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != want {
			t.Errorf("Money(%d): got %s, want %s", m, b, want)
		}
	}
}
