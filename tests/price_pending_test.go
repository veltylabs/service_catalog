package tests

import (
	"errors"
	"testing"
	"webtyp.com/crud"
	"webtyp.com/model"
	"webtyp.com/router"
	servicecatalog "github.com/veltylabs/service_catalog"
)

// needsPrice: el error marca el campo price y su causa es ErrActiveNeedsPrice.
func needsPrice(err error) bool {
	var fe *model.FieldError
	return errors.As(err, &fe) && fe.Field == "price" && fe.Err.Error() == "servicio activo necesita precio"
}

// Un servicio sin precio todavía (0) se puede guardar inactivo; activo necesita precio.
func TestPricePending(t *testing.T) {
	_, broker, m := buildModule(t, "t1")
	c := crud.LoopbackActingAs(broker, "u1", m)

	var spec servicecatalog.Specialty
	var err error
	c.Call(router.Operation{Module: servicecatalog.ModelName, Entity: "specialty", Action: model.Create}, &crud.CreateArgs{Records: []model.Model{&servicecatalog.Specialty{Prefix: "md", Name: "Medicina General"}}}, &spec, func(e error){ err = e })
	if err != nil {
		t.Fatal(err)
	}

	base := servicecatalog.Service{SpecialtyId: spec.Id, Sku: "mdcong", Name: "Consulta general", Currency: servicecatalog.CurrencyCLP, TaxTreatment: servicecatalog.TaxExempt}

	active := base
	active.IsActive = true
	c.Call(router.Operation{Module: servicecatalog.ModelName, Entity: "service", Action: model.Create}, &crud.CreateArgs{Records: []model.Model{&active}}, &active, func(e error){ err = e })
	if err == nil || err.Error() != "loopback: 400 servicio activo necesita precio" {
		t.Fatalf("active with price 0: err = %v, want ErrActiveNeedsPrice", err)
	}

	var svc servicecatalog.Service
	c.Call(router.Operation{Module: servicecatalog.ModelName, Entity: "service", Action: model.Create}, &crud.CreateArgs{Records: []model.Model{&base}}, &svc, func(e error){ err = e })
	if err != nil {
		t.Fatalf("inactive with price 0 must be saved: %v", err)
	}

	svc.IsActive = true
	c.Call(router.Operation{Module: servicecatalog.ModelName, Entity: "service", Action: model.Update}, &crud.UpdateArgs{Record: &svc, Fields: []string{"is_active"}, IDs: []string{svc.Id}}, &svc, func(e error){ err = e })
	if err == nil || err.Error() != "loopback: 400 servicio activo necesita precio" {
		t.Fatalf("activating without price: err = %v, want ErrActiveNeedsPrice", err)
	}

	svc.Price = 30000
	c.Call(router.Operation{Module: servicecatalog.ModelName, Entity: "service", Action: model.Update}, &crud.UpdateArgs{Record: &svc, Fields: []string{"price", "is_active"}, IDs: []string{svc.Id}}, &svc, func(e error){ err = e })
	if err != nil {
		t.Fatalf("activating with price: %v", err)
	}
}
