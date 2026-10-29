//go:build browser

package tests

import (
	"testing"

	"github.com/veltylabs/service_catalog/ui"
	"webtyp.com/view"
)

// Cada pantalla del catálogo se explica sola: el grupo y sus dos hijos
// validan contra view.Screen y tienen descripciones no vacías.
func TestBrowser_EveryScreenIsDescribed(t *testing.T) {
	screens := ui.Screens(view.NewRecords().Caller(&mockCaller{}))
	if len(screens) != 1 {
		t.Fatalf("expected 1 root group screen, got %d", len(screens))
	}
	root := screens[0]
	if err := root.Validate(); err != nil {
		t.Fatalf("root screen validation failed: %v", err)
	}
	if root.ID != ui.ID {
		t.Errorf("root ID = %q, want %q", root.ID, ui.ID)
	}
	if root.Help == "" {
		t.Errorf("root screen has empty help")
	}
	if len(root.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(root.Children))
	}
	for i, c := range root.Children {
		if err := c.Validate(); err != nil {
			t.Fatalf("child %d validation failed: %v", i, err)
		}
		if c.Permission == "" {
			t.Errorf("child %d permission is empty", i)
		}
		if c.Presenter.Description() == "" {
			t.Errorf("child %d presenter must have a non-empty description", i)
		}
	}
}
