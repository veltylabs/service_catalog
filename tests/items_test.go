package tests

import (
	"testing"

	servicecatalog "github.com/veltylabs/service_catalog"
)

func TestItemsDeclaredOnFields(t *testing.T) {
	// Construye un servicio con valores conocidos
	srv := servicecatalog.Service{
		Id: "svc123",
		Name: "Consultation",
		Sku: "SRV-01",
	}

	// Verifica Item para Service
	srvItem := srv.Item()
	if srvItem.ID != "svc123" {
		t.Errorf("expected Service item ID 'svc123', got %q", srvItem.ID)
	}
	if srvItem.Label != "Consultation" {
		t.Errorf("expected Service item Title 'Consultation', got %q", srvItem.Label)
	}
	if srvItem.Description != "SRV-01" {
		t.Errorf("expected Service item Detail 'SRV-01', got %q", srvItem.Description)
	}

	// Construye una especialidad con valores conocidos
	spec := servicecatalog.Specialty{
		Id: "spc456",
		Name: "Cardiology",
		Prefix: "CAR",
	}

	// Verifica Item para Specialty
	specItem := spec.Item()
	if specItem.ID != "spc456" {
		t.Errorf("expected Specialty item ID 'spc456', got %q", specItem.ID)
	}
	if specItem.Label != "Cardiology" {
		t.Errorf("expected Specialty item Title 'Cardiology', got %q", specItem.Label)
	}
	if specItem.Description != "CAR" {
		t.Errorf("expected Specialty item Detail 'CAR', got %q", specItem.Description)
	}
}
