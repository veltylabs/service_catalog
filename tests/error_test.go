package tests

import (
	"testing"

	"webtyp.com/model"
	servicecatalog "github.com/veltylabs/service_catalog"
)

func TestSentinelErrorTexts(t *testing.T) {
	cases := []struct {
		err  error
		text string
	}{
		{servicecatalog.ErrNotFound, "service not found"},
		{servicecatalog.ErrAlreadyExists, "service already exists"},
		{servicecatalog.ErrSpecialtyNotFound, "specialty not found"},
		{servicecatalog.ErrSpecialtyInUse, "specialty in use"},
		{servicecatalog.ErrSpecialtyPrefixExists, "specialty prefix already exists"},
	}

	for _, c := range cases {
		if c.err.Error() != c.text {
			t.Errorf("expected error text %q, got %q", c.text, c.err.Error())
		}
	}
}

func TestSentinelErrorKinds(t *testing.T) {
	cases := []struct {
		err  error
		kind model.ErrorKind
	}{
		{servicecatalog.ErrNotFound, model.ErrorNotFound},
		{servicecatalog.ErrSpecialtyNotFound, model.ErrorNotFound},
		{servicecatalog.ErrSupplyLineNotFound, model.ErrorNotFound},
		{servicecatalog.ErrAlreadyExists, model.ErrorConflict},
		{servicecatalog.ErrSpecialtyPrefixExists, model.ErrorConflict},
		{servicecatalog.ErrSupplyLineExists, model.ErrorConflict},
		{servicecatalog.ErrSpecialtyInUse, model.ErrorInvalid},
		{servicecatalog.ErrNoRecords, model.ErrorInvalid},
		{servicecatalog.ErrNoIds, model.ErrorInvalid},
		{servicecatalog.ErrQuantity, model.ErrorInvalid},
		{servicecatalog.ErrUnsupportedCurrency, model.ErrorInvalid},
		{servicecatalog.ErrCannotRead, model.ErrorInvalid},
	}

	for _, c := range cases {
		if model.KindOf(c.err) != c.kind {
			t.Errorf("expected error kind %q for %q, got %q", c.kind, c.err.Error(), model.KindOf(c.err))
		}
	}
}
