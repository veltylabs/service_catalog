package tests

import (
	"errors"
	"testing"

	"webtyp.com/ddl"
	"webtyp.com/model"

	servicecatalog "github.com/veltylabs/service_catalog"
	"github.com/veltylabs/service_catalog/migrate"
)

// migrateConn es la base de datos que ve la migración: Exec no hace nada, y
// TableColumns responde qué columnas tiene ya cada tabla.
type migrateConn struct {
	existing  []tableColumns
	failTable string
}

type tableColumns struct {
	table   string
	columns []string
}

var errIntrospect = errors.New("introspection failed")

func (c *migrateConn) Exec(string, ...any) error { return nil }

func (c *migrateConn) TableColumns(table string) ([]string, error) {
	if table == c.failTable {
		return nil, errIntrospect
	}
	for _, t := range c.existing {
		if t.table == table {
			return t.columns, nil
		}
	}
	return nil, nil
}

// recordingCompiler guarda, en orden, cada sentencia que pidió la migración.
type recordingCompiler struct{ stmts []ddl.Stmt }

func (r *recordingCompiler) CompileDDL(s ddl.Stmt, _ model.Model) (string, []any, error) {
	r.stmts = append(r.stmts, s)
	return "", nil, nil
}

// migrationModels son las tablas de este repo en orden de llave foránea: el orden que usa Migrate.
func migrationModels() []model.Model {
	return []model.Model{&servicecatalog.Specialty{}, &servicecatalog.Service{}, &servicecatalog.ServiceSupply{}, &servicecatalog.ServiceAgreement{}}
}

func columnsOf(m model.Model) []string {
	var out []string
	for _, f := range m.Schema() {
		out = append(out, f.Name)
	}
	return out
}

func upToDate() []tableColumns {
	var out []tableColumns
	for _, m := range migrationModels() {
		out = append(out, tableColumns{table: m.ModelName(), columns: columnsOf(m)})
	}
	return out
}

// Rama "listado funciona, no falta nada": solo el CREATE TABLE idempotente por tabla, en orden.
func TestMigrate_UpToDate_OnlyEnsuresTables(t *testing.T) {
	comp := &recordingCompiler{}
	if err := migrate.Migrate(&migrateConn{existing: upToDate()}, comp); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	models := migrationModels()
	if len(comp.stmts) != len(models) {
		t.Fatalf("got %d statements, want %d: %+v", len(comp.stmts), len(models), comp.stmts)
	}
	for i, m := range models {
		if comp.stmts[i].Op != ddl.OpCreateTable || comp.stmts[i].Table != m.ModelName() {
			t.Errorf("statement %d: got %+v, want CREATE TABLE %s (foreign-key order)", i, comp.stmts[i], m.ModelName())
		}
	}
}

// Rama "falta un campo en la tabla": exactamente un ADD COLUMN, para ese campo de esa tabla.
func TestMigrate_NewField_AddsOnlyThatColumn(t *testing.T) {
	models := migrationModels()
	last := models[len(models)-1]
	cols := columnsOf(last)
	existing := upToDate()
	existing[len(existing)-1].columns = cols[:len(cols)-1]

	comp := &recordingCompiler{}
	if err := migrate.Migrate(&migrateConn{existing: existing}, comp); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	var adds []ddl.Stmt
	for _, s := range comp.stmts {
		if s.Op == ddl.OpAddColumn {
			adds = append(adds, s)
		}
	}
	if len(adds) != 1 {
		t.Fatalf("got %d ADD COLUMN, want 1: %+v", len(adds), adds)
	}
	if adds[0].Table != last.ModelName() || adds[0].Column == nil || adds[0].Column.Name != cols[len(cols)-1] {
		t.Fatalf("got %+v, want ADD COLUMN %s.%s", adds[0], last.ModelName(), cols[len(cols)-1])
	}
}

// Rama "listar las columnas falló": el error regresa y no se toca ninguna tabla posterior.
func TestMigrate_IntrospectionFails_ReturnsErrorAndStops(t *testing.T) {
	first := migrationModels()[0]
	comp := &recordingCompiler{}
	err := migrate.Migrate(&migrateConn{existing: upToDate(), failTable: first.ModelName()}, comp)
	if !errors.Is(err, errIntrospect) {
		t.Fatalf("got %v, want the introspection error", err)
	}
	for _, s := range comp.stmts {
		if s.Table != first.ModelName() {
			t.Fatalf("touched %s after the failure on %s", s.Table, first.ModelName())
		}
	}
}
