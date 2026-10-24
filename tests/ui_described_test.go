//go:build browser

package tests

import (
	"testing"

	"github.com/veltylabs/service_catalog/ui"
	"webtyp.com/layout/platformd"
)

// Cada pantalla del catálogo se explica sola: el grupo y sus dos hijos
// implementan platformd.Described con un texto no vacío.
func TestBrowser_EveryScreenIsDescribed(t *testing.T) {
	screen, err := ui.Browser(&mockCaller{}, &testIDGen{}, "")
	if err != nil {
		t.Fatalf("ui.Browser: %v", err)
	}
	group, ok := screen.(platformd.GroupModule)
	if !ok {
		t.Fatalf("expected a platformd.GroupModule, got %T", screen)
	}
	all := append([]platformd.UIModule{screen}, group.Children()...)
	if len(all) != 3 {
		t.Fatalf("expected the group and 2 children, got %d screens", len(all))
	}
	for _, m := range all {
		d, ok := m.(platformd.Described)
		if !ok {
			t.Errorf("%s must implement platformd.Described", m.ModelName())
			continue
		}
		if d.Description() == "" {
			t.Errorf("%s has an empty description", m.ModelName())
		}
	}
}
