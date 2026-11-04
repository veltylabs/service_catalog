package tests

import (
	"errors"
	"testing"

	servicecatalog "github.com/veltylabs/service_catalog"
	"webtyp.com/model"
	"webtyp.com/orm"
	"webtyp.com/storage/mem"
)

// needsPrice: el error marca el campo price y su causa es ErrActiveNeedsPrice.
func needsPrice(err error) bool {
	var fe *model.FieldError
	return errors.As(err, &fe) && fe.Field == "price" && fe.Err == servicecatalog.ErrActiveNeedsPrice
}

// Un servicio sin precio todavía (0) se puede guardar inactivo; activo necesita precio.
func TestPricePending(t *testing.T) {
	m, err := servicecatalog.New(orm.New(mem.New()), servicecatalog.Deps{TenantID: "t1", IDs: &MockIDGen{}})
	if err != nil {
		t.Fatal(err)
	}
	spec, err := m.UpsertSpecialty(servicecatalog.Specialty{TenantId: "t1", Prefix: "md", Name: "Medicina General"})
	if err != nil {
		t.Fatal(err)
	}
	base := servicecatalog.Service{TenantId: "t1", SpecialtyId: spec.Id, Sku: "mdcong", Name: "Consulta general",
		Currency: servicecatalog.CurrencyCLP, TaxTreatment: servicecatalog.TaxExempt}

	active := base
	active.IsActive = true
	if _, err := m.CreateService(active); !needsPrice(err) {
		t.Fatalf("active with price 0: err = %v, want ErrActiveNeedsPrice", err)
	}

	svc, err := m.CreateService(base)
	if err != nil {
		t.Fatalf("inactive with price 0 must be saved: %v", err)
	}

	svc.IsActive = true
	if _, err := m.UpdateService(svc); !needsPrice(err) {
		t.Fatalf("activating without price: err = %v, want ErrActiveNeedsPrice", err)
	}

	svc.Price = 30000
	if _, err := m.UpdateService(svc); err != nil {
		t.Fatalf("activating with price: %v", err)
	}
}
