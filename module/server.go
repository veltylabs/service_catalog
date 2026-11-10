//go:build !browser

package module

import (
	"webtyp.com/appkit"
	"webtyp.com/model"
	"webtyp.com/view"

	servicecatalog "github.com/veltylabs/service_catalog"
)

// server es lo que el módulo aporta al servidor: cómo se construye.
var server = appkit.ModuleServer{Build: build}

// En el servidor el módulo no aporta pantallas.
var browser view.ModuleBrowser

func build(d appkit.Deps) (model.ModuleNaming, error) {
	return servicecatalog.New(d.DB, servicecatalog.Deps{IDs: d.IDs, Publisher: d.Events, TenantID: d.TenantID})
}
