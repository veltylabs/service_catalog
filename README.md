# service_catalog
<img src="docs/img/badges.svg">

## ¿Por qué existe?

`service_catalog` es el catálogo de servicios y especialidades de una clínica o centro de salud. Lo encuentras cuando otro
módulo, como `appointment_booking`, necesita saber qué servicios se ofrecen, a qué precio y con qué insumos.

- **Sin esto:** cada módulo guarda su propia lista de servicios y los precios se desincronizan.
- **Con esto:** hay una sola fuente de verdad; la persistencia y el transporte los hace [`crud`](https://github.com/webtyp/crud)
  y las reglas del negocio (SKU único, precio pendiente, especialidad en uso) viven en los archivos `rule_*.go`
  ([arquitectura](docs/ARCHITECTURE.md)).

## Antes de comenzar

- Go instalado y acceso a los módulos `webtyp.com/*` (`GOPRIVATE=webtyp.com/*,github.com/veltylabs/*`).
- Una base de datos conectada como `*orm.DB`; las tablas las crea [`migrate.Migrate`](migrate/).
- Un generador de ids (`model.IDGenerator`), por ejemplo `webtyp.com/unixid`, y el id del inquilino (`TenantID`).

## Guía paso a paso

### Paso 1: crear el módulo

**Objetivo:** conectar el catálogo a una base en memoria.

```go
package main

import (
	servicecatalog "github.com/veltylabs/service_catalog"
	"webtyp.com/fmt"
	"webtyp.com/orm"
	"webtyp.com/storage/mem"
	"webtyp.com/unixid"
)

func main() {
	ids, err := unixid.NewUnixID()
	if err != nil {
		fmt.Println(err)
		return
	}
	mod, err := servicecatalog.New(orm.New(mem.New()), servicecatalog.Deps{IDs: ids, TenantID: "acme"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(mod.ModelName())
}
```

**Resultado esperado:** imprime `service_catalog`.

### Paso 2: operar el catálogo

**Objetivo:** crear un servicio y buscarlo por SKU.

Las operaciones son las de [`crud`](https://github.com/webtyp/crud), pedidas con `router.Operation{Module, Entity, Action}`
(por ejemplo `service.create`). Los ejemplos completos, con sus resultados, están en [tests/catalog_test.go](tests/catalog_test.go).

**Resultado esperado:** el servicio queda guardado; un SKU repetido se rechaza con un error en el campo `sku`.

### Paso 3: ver el catálogo en pantalla

**Objetivo:** abrir la demo con datos en memoria. La demo está en [seed/demo.go](seed/demo.go) y se abre con `webtyp` en la raíz del módulo.

**Resultado esperado:** una lista de servicios con búsqueda.

## Operaciones: sin y con `crud`

| Antes (manejadores a mano) | Después (declaración + reglas) |
|---|---|
| `list_specialties / get_specialty` | `specialty.read` |
| `upsert_specialty` | `specialty.create` / `specialty.update` |
| `delete_specialty` | `specialty.delete` |
| `list_services / get_service` | `service.read` |
| `find_service_by_sku` | `service.read` con filtro `{sku}` |
| `create_service / update_service / upsert_service` | `service.create` / `service.update` |
| `deactivate_service` | `service.update` con `["is_active"]` |
| `delete_service` | `service.delete` |
| `list_currencies` | modelo de lectura `CurrencyOption` |
| `list_agreements / upsert_agreement / delete_agreement` | `service_agreement.read`, `.create`, `.update`, `.delete` |
| `list_service_supplies / upsert_service_supply / delete_service_supply` | `service_supply.read`, `.create`, `.update`, `.delete` |

## Mapa de documentación

| Orden | Documento | Qué pregunta responde y cuándo leerlo |
|---|---|---|
| 1 | [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | ¿Qué entidades hay y quién decide qué? Léelo antes de tocar una regla. |
| 2 | [docs/TESTS.md](docs/TESTS.md) | ¿Qué comportamiento está probado? Léelo antes de cambiar una regla. |
| 3 | [docs/diagrams/database.md](docs/diagrams/database.md) | ¿Cómo se relacionan las tablas? Léelo antes de cambiar el esquema. |
| 4 | [docs/diagrams/service_card.md](docs/diagrams/service_card.md) | ¿Cómo se ve la ficha de un servicio? Léelo antes de cambiar la pantalla. |
