package tests
import (
	"testing"
	"webtyp.com/crud"
	"webtyp.com/model"
	"webtyp.com/router"
	servicecatalog "github.com/veltylabs/service_catalog"
)

func TestReadModel_Currencies(t *testing.T) {
	_, broker, mod := buildModule(t, "tenant-1")
	c := crud.LoopbackActingAs(broker, "u1", mod)

	var err error
	c.Call(router.Operation{Module: servicecatalog.ModelName, Entity: "currency_option", Action: model.Read}, &crud.ReadArgs{}, nil, func(e error) { err = e })
	if err != nil {
		t.Fatal(err)
	}


}

func TestCurrencyValidation(t *testing.T) {
	_, broker, mod := buildModule(t, "tenant-1")
	c := crud.LoopbackActingAs(broker, "u1", mod)

	var svc servicecatalog.Service
	var err error
	c.Call(router.Operation{Module: servicecatalog.ModelName, Entity: "service", Action: model.Create}, &crud.CreateArgs{Records: []model.Model{&servicecatalog.Service{
		SpecialtyId: "spec-1", Sku: "SK1", Name: "Svc 1", Currency: "USD", TaxTreatment: "EXEMPT", Price: 10, IsActive: true,
	}}}, &svc, func(e error) { err = e })

	if err == nil {
		t.Fatal("expected ErrUnsupportedCurrency creating with USD, got nil")
	}
}
