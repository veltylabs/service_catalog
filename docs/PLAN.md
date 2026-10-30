---
PLAN: "fix: migrate with ddl.Sync so a new field reaches an existing table"
EXECUTOR: jules
REVIEWER: none
STATUS: review
SESSION: 4817017239307385509
PR: https://github.com/veltylabs/service_catalog/pull/23
REVIEWED_COMMIT: b734d2436ee61d1ccdffe35d7f942404f648a0be
---

> Read AGENTS.md at the repo root before any change.

# Migrar con `ddl.Sync`

## Contexto

`migrate.Migrate` crea las tablas de este repositorio con `ddl.DB.CreateTable`, que compila a `CREATE TABLE IF NOT EXISTS`:
contra una tabla que ya existe no hace nada. Entonces, cuando un modelo gana un campo, una base de datos que ya está
instalada nunca recibe la nueva columna, y la primera lectura o escritura de ese campo falla en producción.
`ddl.DB.Sync` crea la tabla que falta **y además** agrega las columnas faltantes (también renombra una columna que el modelo declara
como renombrada, y elimina una columna solo cuando no contiene datos). Su flujo, rama por rama, es este diagrama — las pruebas
a continuación recorren sus ramas: https://github.com/webtyp/ddl/blob/main/docs/diagrams/SYNC_FLOW.md

Regla decidida para todo el ecosistema: nadie escribe archivos de migración; cada biblioteca y módulo migra con `Sync`.

## Herramientas

```bash
go install webtyp.com/devflow/cmd/gotest@latest
go get webtyp.com/ddl@v0.0.36
go mod tidy
```

Entorno: `GOPRIVATE=webtyp.com/*,github.com/veltylabs/*`, `GOFLAGS=-mod=mod`.

## APIs utilizadas

| Símbolo | Firma exacta | Etiqueta |
|---|---|---|
| `ddl.New` | `func New(conn Execer, ddlCompiler Compiler) *DB` | `webtyp.com/ddl` v0.0.36 |
| `(*ddl.DB).Sync` | `func (d *DB) Sync(models ...model.Model) error` — modelos en el orden dado; una conexión sin transacciones sincroniza modelo por modelo | v0.0.36 |
| `ddl.Execer` | `interface { Exec(query string, args ...any) error }` | v0.0.36 |
| `ddl.Compiler` | `interface { CompileDDL(s Stmt, m model.Model) (query string, args []any, err error) }` | v0.0.36 |
| `ddl.TableIntrospector` | `interface { TableColumns(table string) ([]string, error) }` — cuando la conexión lo tiene, `Sync` agrega solo las columnas faltantes; cuando `TableColumns` falla, `Sync` devuelve ese mismo error y no toca ningún modelo posterior | v0.0.36 |
| `ddl.Stmt` | `struct { Op Op; Table string; Database string; Column *model.Field; OldName string; ColumnName string }` | v0.0.36 |
| `ddl.OpCreateTable`, `ddl.OpAddColumn` | constantes de tipo `ddl.Op` | v0.0.36 |
| `model.Model` | tiene `ModelName() string` and `Schema() []model.Field`; `model.Field.Name string` | `webtyp.com/model` (ya requerido) |

Si una firma en la etiqueta publicada difiere de esta tabla, detente y nombra el símbolo en el PR; nunca lo adaptes adivinando.

## Etapa 1 — pruebas primero: `tests/migrate_test.go` (new file)

Usa el nombre del paquete que ya usan los otros archivos en `tests/` (`package tests` si la carpeta es nueva). Las pruebas se
escriben **antes** del código: ejecuta `gotest` después de esta etapa y observa cómo falla `TestMigrate_NewField_AddsOnlyThatColumn`
(hoy la migración nunca agrega una columna). Escribe exactamente este archivo:

```go
package tests

import (
	"errors"
	"testing"

	"webtyp.com/ddl"
	"webtyp.com/model"

	servicecatalog "github.com/veltylabs/service_catalog"
	"github.com/veltylabs/service_catalog/migrate"
)

// migrateConn is the database the migration sees: Exec does nothing, and
// TableColumns answers which columns each table already tiene.
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

// recordingCompiler keeps, in order, every statement the migration asked for.
type recordingCompiler struct{ stmts []ddl.Stmt }

func (r *recordingCompiler) CompileDDL(s ddl.Stmt, _ model.Model) (string, []any, error) {
	r.stmts = append(r.stmts, s)
	return "", nil, nil
}

// migrationModels are this repo's tables in foreign-key order: the order Migrate uses.
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

// Branch "listing works, nothing missing": only the idempotent CREATE TABLE per table, in order.
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

// Branch "field missing from the table": exactly one ADD COLUMN, for that field of that table.
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

// Branch "listing the columns failed": the error comes back and no later table is touched.
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
```

- Elimina `migrate/migrate_test.go`: afirmaba un `Exec` por tabla, lo cual `Sync` ya no es. Sus casos se mueven al nuevo archivo. Comprueba: `test ! -e migrate/migrate_test.go`.

Haz un commit después de esta etapa (la prueba en rojo se espera aquí).

## Etapa 2 — `migrate/migrate.go`

Reemplaza **solo el cuerpo** de `Migrate`, manteniendo su firma, los nombres de sus parámetros y los comentarios de documentación existentes
(paquete y función). El cuerpo se convierte exactamente en:

```go
	// Sync y no CreateTable: crea la tabla que falta y además agrega las columnas que un modelo
	// ganó, así un campo nuevo llega a una base ya instalada. Flujo: https://github.com/webtyp/ddl/blob/main/docs/diagrams/SYNC_FLOW.md
	return ddl.New(conn, ddlCompiler).Sync(
		&servicecatalog.Specialty{},
		&servicecatalog.Service{},
		&servicecatalog.ServiceSupply{},
		&servicecatalog.ServiceAgreement{},
	)
```

Elimina cualquier importación que quede sin uso. Ejecuta `gotest`: las tres pruebas pasan. Haz un commit después de esta etapa.

## Etapa 3 — docs

Si `README.md` o cualquier `docs/*.md` dice que las tablas se crean con `CreateTable` (o "crea las tablas" sin
mencionar que se agregan nuevas columnas), cámbialo para decir que `migrate.Migrate` usa `ddl.Sync`, en español, y
enlaza el diagrama https://github.com/webtyp/ddl/blob/main/docs/diagrams/SYNC_FLOW.md. Si ningún documento menciona la migración, agrega una línea a `README.md` bajo su sección de uso o
instalación diciéndolo. Haz un commit después de esta etapa.

## Etapa 4 — language

Traduce al español cada archivo `README.md`, `docs/*`, diagrama y comentario de código de este repositorio que todavía esté en inglés
(solo esos; nunca reescribas lo que ya está en español), y agrega un breve ejemplo de uso (`// Ejemplo: …`) a cada comentario
de símbolo exportado que no tenga uno. El código, los identificadores y las claves de traducción permanecen en inglés. Haz un commit después de esta etapa.

## Comprobación de entrega

1. `gotest` pasa (compila, revisa y ejecuta todas las pruebas). Un PR que no compila, o cuyas pruebas no pudieron ejecutarse, es una entrega fallida.
2. `grep -rn "CreateTable(" --include=*.go . | grep -v _test.go` → vacío.
3. `git diff --stat main` enumera `tests/migrate_test.go`, `migrate/migrate.go`, `go.mod`, `go.sum`, los documentos tocados, el archivo de prueba eliminado — nada más eliminado.

## Reglas de entrega

- Nunca te detengas para pedir permiso o informar progreso: cuando pase la verificación de entrega, haz commit, push y abre el PR. Atrapado
  por 30 minutos: abre el PR de todos modos, listando lo que falla.
- Haz commit después de cada etapa. Nunca ejecutes `git reset`, `git restore`, `git checkout -- <archivo>` o `git stash` en tu propio trabajo.
- Nunca edites un archivo generado (`model_orm.go`): un error de compilación allí significa que el código escrito a mano está mal.
- No cambies los modelos, su orden o cualquier otro paquete.

## Etapas

| # | Archivos | Resultado |
|---|---|---|
| 1 | `tests/migrate_test.go`, delete migrate/migrate_test.go | Tres pruebas, una por rama; la prueba de campo nuevo es roja |
| 2 | `migrate/migrate.go` | `Migrate` usa `Sync`; todo verde |
| 3 | `README.md` / `docs/*` | La migración está documentada con el diagrama |
| 4 | cada `README.md`, `docs/*`, comentario aún en inglés | Español |
