//go:build wasm

package seed

import (
	servicecatalog "github.com/veltylabs/service_catalog"
	"github.com/veltylabs/service_catalog/ui"
	"webtyp.com/crud"
	"webtyp.com/orm"
	"webtyp.com/storage/mem"
	"webtyp.com/view"
)

func Demo() view.Demo {
	db := orm.New(mem.New())
	mod, _ := servicecatalog.New(db, servicecatalog.Deps{TenantID: "tenant-1"})
	Load(mod, "tenant-1")

	return view.Demo{
		Label:      "Catálogo de Servicios",
		Timezone:   "America/Santiago",
		ActingAs:   "u1",
		Modules:    []crud.Module{mod},
		Migrations: nil,
		Screens: func(d view.BrowserDeps) []view.Screen {
			return ui.Screens(d.Caller)
		},
	}
}
