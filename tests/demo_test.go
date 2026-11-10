//go:build wasm

package tests

import (
	"testing"
	"github.com/veltylabs/service_catalog/seed"
	"webtyp.com/view"
)

func TestDemo_Valid(t *testing.T) {
	d := seed.Demo()
	if err := d.Validate(); err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if len(d.Screens(view.BrowserDeps{})) == 0 {
		t.Errorf("expected screens, got 0")
	}
}
