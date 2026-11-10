package tests

import (
	"testing"
	"webtyp.com/crud"
	"webtyp.com/model"
	"webtyp.com/router"
	servicecatalog "github.com/veltylabs/service_catalog"
)

func TestServiceRecipe_LinesInOrder(t *testing.T) {
	_, broker, m := buildModule(t, "t1")
	c := crud.LoopbackActingAs(broker, "u1", m)

	var spec servicecatalog.Specialty
	var err error
	c.Call(router.Operation{Module: servicecatalog.ModelName, Entity: "specialty", Action: model.Create}, &crud.CreateArgs{Records: []model.Model{&servicecatalog.Specialty{Prefix: "md", Name: "Medicina General"}}}, &spec, func(e error){ err = e })
	if err != nil { t.Fatal(err) }

	var svc servicecatalog.Service
	c.Call(router.Operation{Module: servicecatalog.ModelName, Entity: "service", Action: model.Create}, &crud.CreateArgs{Records: []model.Model{&servicecatalog.Service{SpecialtyId: spec.Id, Sku: "mdcong", Name: "Consulta general", Currency: servicecatalog.CurrencyCLP, TaxTreatment: servicecatalog.TaxExempt, Price: 10, IsActive: true}}}, &svc, func(e error){ err = e })
	if err != nil { t.Fatal(err) }

	var sup servicecatalog.ServiceSupply
	c.Call(router.Operation{Module: servicecatalog.ModelName, Entity: "service_supply", Action: model.Create}, &crud.CreateArgs{Records: []model.Model{&servicecatalog.ServiceSupply{ServiceId: svc.Id, ItemId: "item-1", Quantity: 2}}}, &sup, func(e error){ err = e })
	c.Call(router.Operation{Module: servicecatalog.ModelName, Entity: "service_supply", Action: model.Create}, &crud.CreateArgs{Records: []model.Model{&servicecatalog.ServiceSupply{ServiceId: svc.Id, ItemId: "item-2", Quantity: 5}}}, &sup, func(e error){ err = e })

	itemIds, qty, err := m.ServiceRecipe("t1", svc.Id)
	if err != nil {
		t.Fatal(err)
	}
	if len(itemIds) != 2 || itemIds[0] != "item-1" || itemIds[1] != "item-2" || qty[0] != 2 || qty[1] != 5 {
		t.Fatalf("expected item-1 (2), item-2 (5), got %v %v", itemIds, qty)
	}
}

func TestServiceRecipe_Rules(t *testing.T) {
	_, broker, m := buildModule(t, "t1")
	c := crud.LoopbackActingAs(broker, "u1", m)

	var spec servicecatalog.Specialty
	var err error
	c.Call(router.Operation{Module: servicecatalog.ModelName, Entity: "specialty", Action: model.Create}, &crud.CreateArgs{Records: []model.Model{&servicecatalog.Specialty{Prefix: "md", Name: "Medicina General"}}}, &spec, func(e error){ err = e })

	var svc servicecatalog.Service
	c.Call(router.Operation{Module: servicecatalog.ModelName, Entity: "service", Action: model.Create}, &crud.CreateArgs{Records: []model.Model{&servicecatalog.Service{SpecialtyId: spec.Id, Sku: "mdcong", Name: "Consulta general", Currency: servicecatalog.CurrencyCLP, TaxTreatment: servicecatalog.TaxExempt, Price: 10, IsActive: true}}}, &svc, func(e error){ err = e })

	var sup servicecatalog.ServiceSupply
	c.Call(router.Operation{Module: servicecatalog.ModelName, Entity: "service_supply", Action: model.Create}, &crud.CreateArgs{Records: []model.Model{&servicecatalog.ServiceSupply{ServiceId: svc.Id, ItemId: "item-1", Quantity: 2}}}, &sup, func(e error){ err = e })

	// same supply twice
	c.Call(router.Operation{Module: servicecatalog.ModelName, Entity: "service_supply", Action: model.Create}, &crud.CreateArgs{Records: []model.Model{&servicecatalog.ServiceSupply{ServiceId: svc.Id, ItemId: "item-1", Quantity: 5}}}, &sup, func(e error){ err = e })
	if err == nil || err.Error() != "loopback: 409 receta_duplicada" {
		t.Errorf("same supply twice: got %v", err)
	}

	// negative quantity
	c.Call(router.Operation{Module: servicecatalog.ModelName, Entity: "service_supply", Action: model.Create}, &crud.CreateArgs{Records: []model.Model{&servicecatalog.ServiceSupply{ServiceId: svc.Id, ItemId: "item-2", Quantity: -1}}}, &sup, func(e error){ err = e })
	if err == nil || err.Error() != "loopback: 400 cantidad debe ser mayor a cero" {
		t.Errorf("negative quantity: got %v", err)
	}
}
