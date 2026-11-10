package tests

import (
	"testing"
	"webtyp.com/crud"
	"webtyp.com/model"
	"webtyp.com/router"
	servicecatalog "github.com/veltylabs/service_catalog"
	"webtyp.com/orm"
	"webtyp.com/storage/mem"
	"webtyp.com/events/mock"
)


func buildModule(t *testing.T, tenantId string) (*orm.DB, *mock.Broker, *servicecatalog.Module) {
	db := orm.New(mem.New())
	broker := &mock.Broker{}
	m, err := servicecatalog.New(db, servicecatalog.Deps{
		TenantID:   tenantId,
		Publisher:  broker,
		IDs:        &testIDGen{},
		Currencies: []string{"CLP"},
	})
	if err != nil {
		t.Fatal(err)
	}
	return db, broker, m
}

func opCreate(entity string) router.Operation {
	return router.Operation{Module: servicecatalog.ModelName, Entity: entity, Action: model.Create}
}
func opRead(entity string) router.Operation {
	return router.Operation{Module: servicecatalog.ModelName, Entity: entity, Action: model.Read}
}
func opDelete(entity string) router.Operation {
	return router.Operation{Module: servicecatalog.ModelName, Entity: entity, Action: model.Delete}
}

func TestTenant_ReadOtherTenantIsNotFound(t *testing.T) {
	db, broker, mod := buildModule(t, "tenant-a")
	c := crud.LoopbackActingAs(broker, "u1", mod)

	var svc servicecatalog.Service
	var err error
	c.Call(opCreate("service"), &crud.CreateArgs{Records: []model.Model{&servicecatalog.Service{
		SpecialtyId: "spec-1", Sku: "SK1", Name: "Svc 1", Currency: "CLP", TaxTreatment: "EXEMPT", Price: 10, IsActive: true,
	}}}, &svc, func(e error) { err = e })
	if err != nil { t.Fatal(err) }

	modB, _ := servicecatalog.New(db, servicecatalog.Deps{TenantID: "tenant-b", Currencies: []string{"CLP"}, IDs: &testIDGen{}})
	cB := crud.LoopbackActingAs(nil, "u1", modB)

	var readSvc servicecatalog.Service
	cB.Call(opRead("service"), &crud.ReadArgs{ID: svc.Id}, &readSvc, func(e error) { err = e })
	if err == nil || model.KindOf(err) != model.ErrorNotFound {
		t.Fatalf("expected ErrorNotFound, got %v", err)
	}
}

func TestService_DuplicateSkuRejected(t *testing.T) {
	_, broker, mod := buildModule(t, "tenant-1")
	c := crud.LoopbackActingAs(broker, "u1", mod)

	var err error
	var svc servicecatalog.Service
	c.Call(opCreate("service"), &crud.CreateArgs{Records: []model.Model{&servicecatalog.Service{
		SpecialtyId: "spec-1", Sku: "X1", Name: "Svc 1", Currency: "CLP", TaxTreatment: "EXEMPT", Price: 10, IsActive: true,
	}}}, &svc, func(e error) { err = e })
	if err != nil { t.Fatal(err) }

	c.Call(opCreate("service"), &crud.CreateArgs{Records: []model.Model{&servicecatalog.Service{
		SpecialtyId: "spec-2", Sku: "X1", Name: "Svc 2", Currency: "CLP", TaxTreatment: "EXEMPT", Price: 10, IsActive: true,
	}}}, &svc, func(e error) { err = e })

	if err == nil {
		t.Fatal("expected error for duplicate sku")
	}
	if err == nil { t.Fatal() }
	if err.Error() != "loopback: 409 servicio_duplicado" {
		t.Fatalf("expected error service already exists, got %v", err)
	}


}

func TestService_FindBySkuFilter(t *testing.T) {
	_, broker, mod := buildModule(t, "tenant-1")
	c := crud.LoopbackActingAs(broker, "u1", mod)

	var err error
	var svc servicecatalog.Service
	c.Call(opCreate("service"), &crud.CreateArgs{Records: []model.Model{&servicecatalog.Service{
		SpecialtyId: "spec-1", Sku: "X1", Name: "Svc 1", Currency: "CLP", TaxTreatment: "EXEMPT", Price: 10, IsActive: true,
	}}}, &svc, func(e error) { err = e })
	if err != nil { t.Fatal(err) }

	c.Call(opRead("service"), &crud.ReadArgs{Filter: &servicecatalog.Service{Sku: "X1"}}, nil, func(e error) { err = e })
	if err != nil { t.Fatal(err) }


}

func TestSpecialty_DeleteInUseRejected(t *testing.T) {
	_, broker, mod := buildModule(t, "tenant-1")
	c := crud.LoopbackActingAs(broker, "u1", mod)

	var err error
	var spec servicecatalog.Specialty
	c.Call(opCreate("specialty"), &crud.CreateArgs{Records: []model.Model{&servicecatalog.Specialty{
		Prefix: "AA", Name: "Spec AA",
	}}}, &spec, func(e error) { err = e })
	if err != nil { t.Fatal(err) }

	var svc servicecatalog.Service
	c.Call(opCreate("service"), &crud.CreateArgs{Records: []model.Model{&servicecatalog.Service{
		SpecialtyId: spec.Id, Sku: "AA01", Name: "Svc", Currency: "CLP", TaxTreatment: "EXEMPT", Price: 10, IsActive: true,
	}}}, &svc, func(e error) { err = e })
	if err != nil { t.Fatal(err) }

	c.Call(opDelete("specialty"), &crud.DeleteArgs{IDs: []string{spec.Id}}, nil, func(e error) { err = e })
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "loopback: 400 especialidad en uso" {
		t.Fatalf("expected %v, got %v", "loopback: 400 especialidad en uso", err)
	}
}
