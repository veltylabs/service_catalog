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
	d := ddl.New(conn, ddlCompiler)
	if err := d.CreateTable(&servicecatalog.Specialty{}); err != nil {
		return err
	}
	if err := d.CreateTable(&servicecatalog.Service{}); err != nil {
		return err
	}
	if err := d.CreateTable(&servicecatalog.ServiceSupply{}); err != nil {
		return err
	}
	if err := d.CreateTable(&servicecatalog.ServiceAgreement{}); err != nil {
		return err
	}
	return nil
}
