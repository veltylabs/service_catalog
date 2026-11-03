//go:build !wasm

// Package server es service_catalog tal como lo enchufa una app.
package server

import (
	"webtyp.com/appkit"
	"webtyp.com/model"

	servicecatalog "github.com/veltylabs/service_catalog"
	"github.com/veltylabs/service_catalog/migrate"
)

// Module es service_catalog para appkit: va en App.Modules de la app.
//
// Ejemplo: Modules: []appkit.Module{catalog.Module}
var Module = appkit.Module{Name: "service_catalog", Migrate: migrate.Migrate, Build: build}

func build(d appkit.Deps) (model.ModuleNaming, error) {
	return servicecatalog.New(d.DB, servicecatalog.Deps{IDs: d.IDs, Publisher: d.Events, TenantID: d.Tenant})
}
