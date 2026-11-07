package tests

import (
	"testing"
	servicecatalog "github.com/veltylabs/service_catalog"
	"webtyp.com/json"
	"webtyp.com/router"
	"webtyp.com/router/mock"
	"webtyp.com/view"
	"webtyp.com/model"
	"webtyp.com/orm"
	"webtyp.com/storage/mem"
)

func TestSpecialty_ListAcceptsListArgs(t *testing.T) {
	db := orm.New(mem.New())
	tenantID := "tenant-1"
	module, err := servicecatalog.New(db, servicecatalog.Deps{
		TenantID:  tenantID,
		IDs:       &MockIDGen{},
	})
	if err != nil {
		t.Fatal(err)
	}
		r := &mock.Router{}
	r.Configure(mock.Config{
		Authn: func(next router.HandlerFunc) router.HandlerFunc {
			return func(ctx router.Context) {
				ctx.SetUserID("user-123")
				next(ctx)
			}
		},
		Authorize: func(userID string, res model.Resource, act model.Action) bool {
			return true
		},
	})
	module.MountOperations(r)

	// Create seed data
	for i := 0; i < 3; i++ {
		_, err := module.UpsertSpecialty(servicecatalog.Specialty{
			TenantId: tenantID,
			Prefix:   "p" + string(rune('a'+i)),
			Name:     "Test Specialty " + string(rune('A'+i)),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// 1. Limit 1 -> exactly 1 row
	listArgs := view.ListArgs{Limit: 1}
	var listBytes []byte
	_ = json.Encode(&listArgs, &listBytes)

	ctxList := &mock.Context{InBody: listBytes}
	r.Invoke("OP", "/"+servicecatalog.OpListSpecialties, ctxList)
	if ctxList.Status != 0 && ctxList.Status != 200 {
		t.Fatalf("expected list status success, got %d", ctxList.Status)
	}

	var listResult servicecatalog.SpecialtyList
	if err := json.Decode(ctxList.ResponseBody(), &listResult); err != nil {
		t.Fatal(err)
	}
	if listResult.Len() != 1 {
		t.Fatalf("expected 1 row, got %d", listResult.Len())
	}

	firstID := listResult[0].Id

	// 2. Search exact row -> that row
	searchArgs := view.ListArgs{Search: "Test Specialty B", Limit: 10}
	var searchBytes []byte
	_ = json.Encode(&searchArgs, &searchBytes)

	ctxSearch := &mock.Context{InBody: searchBytes}
	r.Invoke("OP", "/"+servicecatalog.OpListSpecialties, ctxSearch)

	var searchResult servicecatalog.SpecialtyList
	json.Decode(ctxSearch.ResponseBody(), &searchResult)

	if searchResult.Len() != 1 || searchResult[0].Name != "Test Specialty B" {
		t.Fatalf("expected to find 'Test Specialty B', got %d results", searchResult.Len())
	}

	// 3. After the first row -> first row is not in result
	afterArgs := view.ListArgs{After: firstID, Limit: 10}
	var afterBytes []byte
	_ = json.Encode(&afterArgs, &afterBytes)

	ctxAfter := &mock.Context{InBody: afterBytes}
	r.Invoke("OP", "/"+servicecatalog.OpListSpecialties, ctxAfter)

	var afterResult servicecatalog.SpecialtyList
	json.Decode(ctxAfter.ResponseBody(), &afterResult)

	if afterResult.Len() != 2 {
		t.Fatalf("expected 2 rows after the first, got %d", afterResult.Len())
	}
	if afterResult[0].Id == firstID {
		t.Fatalf("first row should not be in the results when querying with After")
	}
}

func TestService_ListAcceptsListArgs(t *testing.T) {
	db := orm.New(mem.New())
	tenantID := "tenant-1"
	module, err := servicecatalog.New(db, servicecatalog.Deps{
		TenantID:  tenantID,
		IDs:       &MockIDGen{},
	})
	if err != nil {
		t.Fatal(err)
	}
		r := &mock.Router{}
	r.Configure(mock.Config{
		Authn: func(next router.HandlerFunc) router.HandlerFunc {
			return func(ctx router.Context) {
				ctx.SetUserID("user-123")
				next(ctx)
			}
		},
		Authorize: func(userID string, res model.Resource, act model.Action) bool {
			return true
		},
	})
	module.MountOperations(r)

	spec, _ := module.UpsertSpecialty(servicecatalog.Specialty{
		TenantId: tenantID,
		Prefix:   "xx",
		Name:     "Spec",
	})

	// Create seed data
	for i := 0; i < 3; i++ {
		_, err := module.CreateService(servicecatalog.Service{
			TenantId:     tenantID,
			SpecialtyId:  spec.Id,
			Sku:          "SKU" + string(rune('0'+i)),
			Name:         "Test Service " + string(rune('A'+i)),
			TaxTreatment: servicecatalog.TaxExempt,
			Currency:     servicecatalog.CurrencyCLP,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// 1. Limit 1 -> exactly 1 row
	listArgs := view.ListArgs{Limit: 1}
	var listBytes []byte
	_ = json.Encode(&listArgs, &listBytes)

	ctxList := &mock.Context{InBody: listBytes}
	r.Invoke("OP", "/"+servicecatalog.OpListServices, ctxList)
	if ctxList.Status != 0 && ctxList.Status != 200 {
		t.Fatalf("expected list status success, got %d", ctxList.Status)
	}

	var listResult servicecatalog.ServiceList
	if err := json.Decode(ctxList.ResponseBody(), &listResult); err != nil {
		t.Fatal(err)
	}
	if listResult.Len() != 1 {
		t.Fatalf("expected 1 row, got %d", listResult.Len())
	}

	firstID := listResult[0].Id

	// 2. Search exact row -> that row
	searchArgs := view.ListArgs{Search: "Test Service B", Limit: 10}
	var searchBytes []byte
	_ = json.Encode(&searchArgs, &searchBytes)

	ctxSearch := &mock.Context{InBody: searchBytes}
	r.Invoke("OP", "/"+servicecatalog.OpListServices, ctxSearch)

	var searchResult servicecatalog.ServiceList
	json.Decode(ctxSearch.ResponseBody(), &searchResult)

	if searchResult.Len() != 1 || searchResult[0].Name != "Test Service B" {
		t.Fatalf("expected to find 'Test Service B', got %d results", searchResult.Len())
	}

	// 3. After the first row -> first row is not in result
	afterArgs := view.ListArgs{After: firstID, Limit: 10}
	var afterBytes []byte
	_ = json.Encode(&afterArgs, &afterBytes)

	ctxAfter := &mock.Context{InBody: afterBytes}
	r.Invoke("OP", "/"+servicecatalog.OpListServices, ctxAfter)

	var afterResult servicecatalog.ServiceList
	json.Decode(ctxAfter.ResponseBody(), &afterResult)

	if afterResult.Len() != 2 {
		t.Fatalf("expected 2 rows after the first, got %d", afterResult.Len())
	}
	if afterResult[0].Id == firstID {
		t.Fatalf("first row should not be in the results when querying with After")
	}
}