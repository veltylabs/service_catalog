package migrate

import (
	"webtyp.com/ddl"

	servicecatalog "github.com/veltylabs/service_catalog"
)

// Migrate reconcilia el esquema de base de datos que pertenece a service_catalog: Specialty,
// Service y ServiceAgreement, en orden de dependencia FK (ServiceAgreement hace referencia
// a Service, Service hace referencia a Specialty).
//
// Deliberadamente NO es llamado por New, y vive deliberadamente en su propio paquete
// en lugar de un archivo en el paquete raíz: nada en la ruta de compilación WASM de
// una aplicación consumidora (su view.go, que importa el paquete raíz servicecatalog para
// servicecatalog.NewView) importa jamás "github.com/veltylabs/service_catalog/migrate" — así
// webtyp.com/ddl nunca entra en ese grafo de compilación, independientemente de las etiquetas
// de compilación en el lado del consumidor.
//
// conn es un ddl.Execer, no un *orm.DB, por lo que un transporte en tiempo de despliegue que
// solo puede ejecutar DDL lo satisface. El RawConn() de un *orm.DB también lo satisface,
// para invocadores locales/de prueba.
// Ejemplo: err := migrate.Migrate(conn, compiler)
func Migrate(conn ddl.Execer, ddlCompiler ddl.Compiler) error {
	// Sync y no CreateTable: crea la tabla que falta y además agrega las columnas que un modelo
	// ganó, así un campo nuevo llega a una base ya instalada. Flujo: https://github.com/webtyp/ddl/blob/main/docs/diagrams/SYNC_FLOW.md
	return ddl.New(conn, ddlCompiler).Sync(
		&servicecatalog.Specialty{},
		&servicecatalog.Service{},
		&servicecatalog.ServiceSupply{},
		&servicecatalog.ServiceAgreement{},
	)
}
