package tests

import (
	"testing"

	servicecatalog "github.com/veltylabs/service_catalog"
	"webtyp.com/model"
)

func TestRowsDeclaredOnFields(t *testing.T) {
	// Construye un servicio con valores conocidos
	srv := servicecatalog.Service{
		Id: "svc123",
		Name: "Consultation",
		Sku: "SRV-01",
	}

	// Verifica RowOf para Service
	srvRow := model.RowOf(&srv)
	if srvRow.ID != "svc123" {
		t.Errorf("expected Service row ID 'svc123', got %q", srvRow.ID)
	}
	if srvRow.Title != "Consultation" {
		t.Errorf("expected Service row Title 'Consultation', got %q", srvRow.Title)
	}
	if srvRow.Detail != "SRV-01" {
		t.Errorf("expected Service row Detail 'SRV-01', got %q", srvRow.Detail)
	}

	// Construye una especialidad con valores conocidos
	spec := servicecatalog.Specialty{
		Id: "spc456",
		Name: "Cardiology",
		Prefix: "CAR",
	}

	// Verifica RowOf para Specialty
	specRow := model.RowOf(&spec)
	if specRow.ID != "spc456" {
		t.Errorf("expected Specialty row ID 'spc456', got %q", specRow.ID)
	}
	if specRow.Title != "Cardiology" {
		t.Errorf("expected Specialty row Title 'Cardiology', got %q", specRow.Title)
	}
	if specRow.Detail != "CAR" {
		t.Errorf("expected Specialty row Detail 'CAR', got %q", specRow.Detail)
	}
}
