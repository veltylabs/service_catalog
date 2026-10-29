package tests

import (
	"testing"
	"webtyp.com/json"
	"webtyp.com/model"
	"webtyp.com/router"
	"webtyp.com/router/mock"

	servicecatalog "github.com/veltylabs/service_catalog"
	"webtyp.com/orm"
	"webtyp.com/storage/mem"
)

func TestTenantIsolation(t *testing.T) {
	db := orm.New(mem.New())
	module, err := servicecatalog.New(db, servicecatalog.Deps{
		TenantID: "tenant-1",
		IDs:      &MockIDGen{},
	})
	if err != nil {
		t.Fatal(err)
	}

	tenantA := "tenant-A"
	tenantB := "tenant-B"

	// Create specialty for tenant A
	specA, err := module.UpsertSpecialty(servicecatalog.Specialty{
		TenantId: tenantA,
		Prefix:   "md",
		Name:     "Medicina General",
	})
	if err != nil {
		t.Fatalf("failed to create specialty A: %v", err)
	}

	// Tenant B cannot get tenant A's specialty
	_, err = module.GetSpecialty(tenantB, specA.Id)
	if err != error(servicecatalog.ErrSpecialtyNotFound) {
		t.Errorf("expected ErrSpecialtyNotFound for tenant B getting tenant A's specialty, got %v", err)
	}

	// Create an item as tenant A
	itemA := servicecatalog.Service{
		TenantId:     tenantA,
		SpecialtyId:  specA.Id,
		Sku:          "md-SKU-A",
		Name:         "Tenant A Product",
		Price:        10.0,
		Currency:     "CLP",
		IsActive:     true,
		TaxTreatment: "EXEMPT",
	}

	createdA, err := module.CreateService(itemA)
	if err != nil {
		t.Fatalf("failed to create tenant A item: %v", err)
	}

	// 1. Tenant B must NOT be able to Get tenant A's item
	_, err = module.GetService(tenantB, createdA.Id)
	if err == nil || err.Error() != servicecatalog.ErrNotFound.Error() {
		t.Errorf("expected ErrNotFound for tenant B getting tenant A's item, got %v", err)
	}

	// 2. Tenant B must NOT be able to Update tenant A's item
	createdA.TenantId = tenantB
	createdA.Name = "Hijacked Name"
	_, err = module.UpdateService(createdA)
	if err == nil || err.Error() != servicecatalog.ErrNotFound.Error() {
		t.Errorf("expected ErrNotFound for tenant B updating tenant A's item, got %v", err)
	}

	// Revert tenant id on createdA helper
	createdA.TenantId = tenantA

	// 3. Tenant B must NOT be able to Deactivate tenant A's item
	err = module.DeactivateService(tenantB, createdA.Id)
	if err == nil || err.Error() != servicecatalog.ErrNotFound.Error() {
		t.Errorf("expected ErrNotFound for tenant B deactivating tenant A's item, got %v", err)
	}

	// 4. Tenant B must NOT be able to Delete tenant A's item
	err = module.DeleteService(tenantB, createdA.Id)
	if err == nil || err.Error() != servicecatalog.ErrNotFound.Error() {
		t.Errorf("expected ErrNotFound for tenant B deleting tenant A's item, got %v", err)
	}

	// Create an agreement for tenant A's item
	agA := servicecatalog.ServiceAgreement{
		TenantId:  tenantA,
		ServiceId: createdA.Id,
		Insurer:   "FONASA",
		Code:      "F-A",
		Price:     50.0,
		IsActive:  true,
	}

	createdAgA, err := module.UpsertServiceAgreement(agA)
	if err != nil {
		t.Fatalf("failed to create agreement: %v", err)
	}

	// 5. Tenant B must NOT be able to Upsert/Update tenant A's agreement
	createdAgA.TenantId = tenantB
	createdAgA.Insurer = "FONASA Hijacked"
	_, err = module.UpsertServiceAgreement(createdAgA)
	if err == nil || err.Error() != servicecatalog.ErrNotFound.Error() {
		t.Errorf("expected ErrNotFound for tenant B updating tenant A's agreement, got %v", err)
	}

	// Revert tenant ID
	createdAgA.TenantId = tenantA

	// 6. Tenant B must NOT be able to Delete tenant A's agreement
	err = module.DeleteServiceAgreement(tenantB, createdAgA.Id)
	if err == nil || err.Error() != servicecatalog.ErrNotFound.Error() {
		t.Errorf("expected ErrNotFound for tenant B deleting tenant A's agreement, got %v", err)
	}
}
func TestTenantFromServer(t *testing.T) {
	db := orm.New(mem.New())

	modA, _ := servicecatalog.New(db, servicecatalog.Deps{TenantID: "tenant-A", IDs: &MockIDGen{}})
	modB, _ := servicecatalog.New(db, servicecatalog.Deps{TenantID: "tenant-B", IDs: &MockIDGen{}})

	rA := &mock.Router{}
	rA.Configure(mock.Config{Authn: func(next router.HandlerFunc) router.HandlerFunc {
		return func(ctx router.Context) {
			ctx.SetUserID("u1")
			next(ctx)
		}
	}, Authorize: func(string, model.Resource, model.Action) bool { return true }})
	modA.MountOperations(rA)

	rB := &mock.Router{}
	rB.Configure(mock.Config{Authn: func(next router.HandlerFunc) router.HandlerFunc {
		return func(ctx router.Context) {
			ctx.SetUserID("u1")
			next(ctx)
		}
	}, Authorize: func(string, model.Resource, model.Action) bool { return true }})
	modB.MountOperations(rB)

	// Upsert via module A
	specArg := servicecatalog.UpsertSpecialtyArgs{
		Records: []servicecatalog.Specialty{{
			Prefix: "ma",
			Name:   "Med A",
		}},
	}
	var inA []byte
	_ = json.Encode(&specArg, &inA)
	ctxA := &mock.Context{InBody: inA}
	rA.Invoke("OP", "/"+servicecatalog.OpUpsertSpecialty, ctxA)
	if ctxA.Status != 0 && ctxA.Status != 200 {
		t.Fatalf("modA upsert failed: %d", ctxA.Status)
	}

	// List via module B
	listArg := servicecatalog.ListSpecialtiesArgs{}
	var inListB []byte
	_ = json.Encode(&listArg, &inListB)
	ctxListB := &mock.Context{InBody: inListB}
	rB.Invoke("OP", "/"+servicecatalog.OpListSpecialties, ctxListB)

	var listResB servicecatalog.SpecialtyList
	_ = json.Decode(ctxListB.ResponseBody(), &listResB)
	if listResB.Len() != 0 {
		t.Fatalf("expected module B to see 0 specialties, got %d", listResB.Len())
	}

	// Upsert via module A with tenant-B should be stored with tenant-A
	specArg2 := servicecatalog.UpsertSpecialtyArgs{
		Records: []servicecatalog.Specialty{{
			TenantId: "tenant-B", // hijacked
			Prefix:   "mb",
			Name:     "Med B",
		}},
	}
	var inA2 []byte
	_ = json.Encode(&specArg2, &inA2)
	ctxA2 := &mock.Context{InBody: inA2}
	rA.Invoke("OP", "/"+servicecatalog.OpUpsertSpecialty, ctxA2)

	var created servicecatalog.Specialty
	_ = json.Decode(ctxA2.ResponseBody(), &created)

	if created.TenantId != "tenant-A" {
		t.Fatalf("expected tenant-A, got %s", created.TenantId)
	}
}

func TestNew_TenantRequired(t *testing.T) {
	db := orm.New(mem.New())
	_, err := servicecatalog.New(db, servicecatalog.Deps{
		IDs: &MockIDGen{},
	})
	if err == nil || err.Error() != "service_catalog: Deps.TenantID is required" {
		t.Fatalf("expected Deps.TenantID is required error, got %v", err)
	}
}
