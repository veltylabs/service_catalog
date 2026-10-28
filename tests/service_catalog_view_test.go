//go:build browser

package tests

import (
	"testing"

	servicecatalog "github.com/veltylabs/service_catalog"
	"github.com/veltylabs/service_catalog/ui"
	"webtyp.com/json"
	"webtyp.com/model"
	"webtyp.com/view"
)

func TestView_ScreensStructure(t *testing.T) {
	mock := &mockCaller{}
	screens := ui.Screens(mock)
	if len(screens) != 1 {
		t.Fatalf("expected 1 root screen, got %d", len(screens))
	}

	root := screens[0]
	if root.ID != ui.ID {
		t.Errorf("expected root ID %q, got %q", ui.ID, root.ID)
	}
	if root.Label != ui.Label {
		t.Errorf("expected root Label %q, got %q", ui.Label, root.Label)
	}

	if len(root.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(root.Children))
	}

	if root.Children[0].Label != "Servicios" {
		t.Errorf("expected first child Label 'Servicios', got %q", root.Children[0].Label)
	}
	if root.Children[1].Label != "Especialidades" {
		t.Errorf("expected second child Label 'Especialidades', got %q", root.Children[1].Label)
	}
}

func TestView_LoadList(t *testing.T) {
	mock := &mockCaller{
		onCall: func(op string, args model.Encodable, into model.Decodable) error {
			if op == servicecatalog.ModelName+"."+servicecatalog.OpListServices {
				list := servicecatalog.ServiceList{
					{Id: "1", Name: "Service One", Sku: "SKU1"},
					{Id: "2", Name: "Service Two", Sku: "SKU2"},
				}
				var out []byte
				if err := json.Encode(&list, &out); err != nil {
					return err
				}
				return json.Decode(string(out), into)
			}
			return nil
		},
	}

	screens := ui.Screens(mock)
	pres := screens[0].Children[0].Presenter

	pres.Reload(nil)

	if mock.lastOp != servicecatalog.ModelName+"."+servicecatalog.OpListServices {
		t.Errorf("expected last op %q, got %q", servicecatalog.ModelName+"."+servicecatalog.OpListServices, mock.lastOp)
	}

	items := pres.Items()
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if model.RowOf(&items[0]).Title != "Service One" {
		t.Errorf("expected item 0 label 'Service One', got %q", model.RowOf(&items[0]).Title)
	}
}

func TestRowsDeclaredOnFields(t *testing.T) {
	spec := servicecatalog.Specialty{Id: "1", Name: "Spec Name", Prefix: "SN"}
	specRow := model.RowOf(&spec)
	if specRow.Title != "Spec Name" {
		t.Errorf("expected specialty title 'Spec Name', got %q", specRow.Title)
	}
	if specRow.Detail != "SN" {
		t.Errorf("expected specialty detail 'SN', got %q", specRow.Detail)
	}

	svc := servicecatalog.Service{Id: "2", Name: "Svc Name", Sku: "SKU-1"}
	svcRow := model.RowOf(&svc)
	if svcRow.Title != "Svc Name" {
		t.Errorf("expected service title 'Svc Name', got %q", svcRow.Title)
	}
	if svcRow.Detail != "SKU-1" {
		t.Errorf("expected service detail 'SKU-1', got %q", svcRow.Detail)
	}
}

func TestView_SaveAndDelete(t *testing.T) {
	seed := &servicecatalog.Service{
		Id: "svc1", TenantId: "tenant1", Sku: "SKU1", Name: "Service One",
		Description: "General consultation service",
		Price:       100, Currency: "CLP", IsActive: true, UpdatedAt: 1,
	}

	mock := &mockCaller{
		onCall: func(op string, args model.Encodable, into model.Decodable) error {
			switch op {
			case servicecatalog.ModelName + "." + servicecatalog.OpListServices:
				list := servicecatalog.ServiceList{seed}
				var out []byte
				_ = json.Encode(&list, &out)
				return json.Decode(string(out), into)
			}
			return nil
		},
	}

	screens := ui.Screens(mock)
	pres := screens[0].Children[0].Presenter

	pres.Reload(nil)

	record := pres.Select(seed.Id)
	if record == nil {
		t.Fatal("expected Select to return the cached record")
	}

	saver, ok := pres.(view.Saver)
	if !ok {
		t.Fatal("expected Presenter to implement view.Saver")
	}
	var serr error
	saver.Save([]model.Model{record}, func(err error) { serr = err })
	if serr != nil {
		t.Fatalf("Save: %v", serr)
	}
	if mock.lastOp != servicecatalog.ModelName+"."+servicecatalog.OpUpsertService {
		t.Errorf("expected last op %q after save, got %q", servicecatalog.ModelName+"."+servicecatalog.OpUpsertService, mock.lastOp)
	}

	updater, ok := pres.(view.Updater)
	if !ok {
		t.Fatal("expected Presenter to implement view.Updater")
	}
	var uerr error
	updater.Update([]string{seed.Id}, record, []string{"name"}, func(err error) { uerr = err })
	if uerr != nil {
		t.Fatalf("Update: %v", uerr)
	}
	if mock.lastOp != servicecatalog.ModelName+"."+servicecatalog.OpUpdateService {
		t.Errorf("expected last op %q after update, got %q", servicecatalog.ModelName+"."+servicecatalog.OpUpdateService, mock.lastOp)
	}

	deleter, ok := pres.(view.Deleter)
	if !ok {
		t.Fatal("expected Presenter to implement view.Deleter")
	}
	var derr error
	deleter.Delete([]string{seed.Id}, func(err error) { derr = err })
	if derr != nil {
		t.Fatalf("Delete: %v", derr)
	}
	if mock.lastOp != servicecatalog.ModelName+"."+servicecatalog.OpDeleteService {
		t.Errorf("expected last op %q after delete, got %q", servicecatalog.ModelName+"."+servicecatalog.OpDeleteService, mock.lastOp)
	}
}