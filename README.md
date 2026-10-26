# service-catalog
<img src="docs/img/badges.svg">

Catálogo de los **servicios** que ofrece una empresa y las **especialidades** que los clasifican. Diseñado para el ecosistema Velty.
Los productos y suministros son otros módulos (`pharmacy_manager`, `supply_manager`).

## Inicio Rápido

```go
import (
	"webtyp.com/orm"
	"webtyp.com/storage/mem"
	servicecatalog "github.com/veltylabs/service_catalog"
)

// En una aplicación real, inyecta la conexión a la base de datos elegida por la app compositora
db := orm.New(mem.New())
idGenerator := &MyIDGenerator{} // implementa model.IDGenerator

module, _ := servicecatalog.New(db, servicecatalog.Deps{
	IDs:        idGenerator,
	Currencies: []string{servicecatalog.CurrencyCLP}, // Si no se indica, el catálogo acepta solo CLP
})

// Crear especialidad
specialty, _ := module.UpsertSpecialty(servicecatalog.Specialty{
	TenantId: "my-tenant",
	Prefix:   "md",
	Name:     "Medicina General",
})

// Crear un servicio
service, _ := module.CreateService(servicecatalog.Service{
	TenantId:    "my-tenant",
	SpecialtyId: specialty.Id,
	Sku:         "md0001",
	Name:        "Consulta General",
	Price:       100,
	Currency:    "USD",
	IsActive:    true,
})
```

## Operaciones
| Nombre Op | Descripción |
|-----------|-------------|
| `list_specialties` | Listar especialidades para un inquilino |
| `list_currencies` | Listar monedas aceptadas |
| `get_specialty` | Obtener especialidad por ID |
| `upsert_specialty` | Crear o actualizar especialidad |
| `delete_specialty` | Eliminar especialidad |
| `list_services` | Listar servicios para un inquilino |
| `get_service` | Obtener servicio por ID |
| `find_service_by_sku` | Encontrar servicio por SKU |
| `create_service` | Crear nuevo servicio |
| `update_service` | Actualizar servicio existente |
| `upsert_service` | Crear o actualizar servicio |
| `deactivate_service` | Borrado lógico de servicio |
| `delete_service` | Borrado físico de servicio |
| `list_agreements` | Listar convenios de facturación |
| `upsert_agreement` | Crear o actualizar convenio |
| `delete_agreement` | Eliminar convenio |

## Lectores (Inter-Módulo)
- `ServiceExists(tenantId, serviceId string) (bool, error)`: Comprueba si un servicio existe y está activo.
- `ServiceCard(tenantId, serviceId string) (name string, price float64, currency string, active bool, err error)`: Obtiene detalles del servicio para facturación/reservas.

## Vista y demo

Este módulo contiene su propio componente UI (`ui/`) y datos semilla (`seed/`).

- **`ui/`**: Exporta `ID` (`catalog_item`), `Label` (`Catálogo`), y `Screens(caller)`. La especialidad se configura como un `view.Lookup`; sin selector junto al título. Cada pantalla lleva una descripción de sí misma a través de opciones como `view.WithDescription`.
- **`seed/`**: Exporta `seed.Load(module, tenantID)`, poblando 3 especialidades canónicas y 1 servicio en cada una.

## Archivos Clave
| Archivo | Propósito |
|------|---------|
| `model.go` | Estructuras de datos y esquemas |
| `module.go` | `Deps`, `New`, los métodos de dominio y los lectores (`ServiceExists`, `ServiceCard`) |
| `errors.go` | Constantes de error, `statusOf`, `writeError` |
| `ops.go` | Nombres de operaciones, las formas en red de la vista, `MountOperations` y los manejadores |
| `recipe.go` | La receta de insumos de un servicio y su lector `ServiceRecipe` |
| `migration.go` | Migración de prefijo SKU a especialidad |
| `model_orm.go` | Ayudantes ORM generados |
| `ui/` | Presentador de vista del módulo y componentes |
| `seed/` | Cargador de datos de ejemplo para demo |
| `tests/` | Pruebas funcionales y unitarias |
