package tests

import (
	"testing"
	"webtyp.com/crud"
	"webtyp.com/router/mock"
)

func TestMount_MatchesDefinition(t *testing.T) {
	_, _, mod := buildModule(t, "tenant-1")
	reg := &mock.Router{}
	err := crud.Mount(reg, nil, nil, mod)
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
}
