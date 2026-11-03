//go:build !browser

package module

import (
	"webtyp.com/appkit"
	"webtyp.com/model"

	servicecatalog "github.com/veltylabs/service_catalog"
	"github.com/veltylabs/service_catalog/migrate"
)

// server es lo que el módulo aporta al servidor: sus tablas y cómo se construye.
var server = appkit.ModuleServer{Migrate: migrate.Migrate, Build: build}

// En el servidor el módulo no aporta pantallas.
var browser appkit.ModuleBrowser

func build(d appkit.Deps) (model.ModuleNaming, error) {
	return servicecatalog.New(d.DB, servicecatalog.Deps{IDs: d.IDs, Publisher: d.Events, TenantID: d.TenantID})
}
