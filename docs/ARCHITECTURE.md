# Arquitectura del Catálogo de Servicios

## Alcance del Dominio
El módulo `service-catalog` gestiona el catálogo de servicios y sus especialidades para el ecosistema Velty (los productos y suministros viven en `pharmacy_manager` y `supply_manager`). Sirve como fuente de verdad para otros módulos, como `appointment-booking`.

## Entidades
- **Specialty** (Especialidad): taxonomía con un código de referencia de dos letras (prefijo de los SKU).
- **Service** (Servicio): Un servicio que la empresa proporciona. Referencia a un `Specialty` mediante `specialty_id`. Contiene SKU, nombre, descripción, precio, moneda y un `tax_treatment` opcional (`TAXED` o `EXEMPT`; la tasa no se almacena aquí).
  - *Nota de Diseño*: otros módulos (como `appointment-booking`) necesitan la *capacidad* de ser reservable, mantenida en un solo lugar: `ServiceExists`.
- **Agreement** (Convenio): Representa un convenio de facturación/aseguradora asociado a un servicio. Un servicio puede tener múltiples convenios, cada uno especificando una aseguradora (ej., FONASA, Isapre X), código, precio, y estado activo.

## Patrones
- **Arnés de módulo reutilizable**: el módulo está acoplado solo a contratos publicados, no a infraestructura concreta — ver `AGENTS.md` (la raíz de este repo) para la lista completa de permitidos/prohibidos a la que debe adherirse este módulo:
    - `orm.DB` para el almacenamiento (agnóstico al backend — envuelve cualquier `storage.Conn` que la app inyecte).
    - `ddl.CreateTable`/`ddl.Sync` (sobre `db.RawConn()`) para la migración del propio esquema del módulo en `New()` — reemplaza al obsoleto `orm.DB.CreateTable`.
    - `router.OpModule` (`ModelName()` + `MountOps(reg router.OpRegistry)`) para el transporte — el módulo nunca ve `router.Router`/`router.APIModule`, y nunca importa `tinywasm/mcp`.
    - `model.IDGenerator` para la identidad (`Deps.IDs`, requerido — el módulo nunca construye el suyo propio).
    - `events.Publisher` para actualizaciones controladas por eventos (`Deps.Publisher`, opcional — `nil` desactiva la publicación silenciosamente).
    - `view.Presenter` (`NewView(caller router.Caller) view.Presenter`, `NewSpecialtyView(caller router.Caller) view.Presenter`) para UI, construido solo con `view`+`model`+`router` — la aplicación elige el renderizador. La especialidad es un campo de formulario cargado con `SetOptions`; sin selector junto al título.
    - Las pruebas se ejecutan contra `storage/mem`, nunca `tinywasm/sqlite`.
- **Multi-inquilino**: Toda especialidad, servicio y convenio está asociado con un `tenant_id`. La unicidad de SKU y prefijo se imponen por inquilino.
- **Monedas Aceptadas**: `Deps.Currencies` establece las monedas que acepta el módulo (ISO 4217). Si está vacío, por defecto es CLP.
- **Eventos tipados**: cada evento publicado lleva un payload `model.Encodable` (`&Specialty`/`&Service`/`&Agreement`), nunca un simple `map`.

## Operaciones (vía `MountOps`)
| Op | Acción | Recurso | Descripción |
|-----------|--------|----------|-------------|
| `list_specialties` | `r` | `specialty` | Listar especialidades para un inquilino |
| `get_specialty` | `r` | `specialty` | Obtener especialidad por ID |
| `upsert_specialty` | `c`\|`u` | `specialty` | Crear o actualizar una especialidad |
| `delete_specialty` | `d` | `specialty` | Eliminar una especialidad |
| `list_services` | `r` | `service` | Listar servicios para un inquilino |
| `get_service` | `r` | `service` | Obtener servicio por ID |
| `find_service_by_sku` | `r` | `service` | Encontrar servicio por SKU |
| `create_service` | `c` | `service` | Crear nuevo servicio |
| `update_service` | `u` | `service` | Actualizar servicio existente |
| `upsert_service` | `c`\|`u` | `service` | Crear o actualizar servicio |
| `deactivate_service` | `u` | `service` | Borrado lógico de servicio |
| `delete_service` | `d` | `service` | Borrado físico de servicio |
| `list_agreements` | `r` | `service_agreement` | Listar convenios de un servicio |
| `upsert_agreement` | `c`\|`u` | `service_agreement` | Crear o actualizar un convenio |
| `delete_agreement` | `d` | `service_agreement` | Eliminar un convenio |

## Ejemplo de Raíz de Composición
```go
cat, _ := servicecatalog.New(db, servicecatalog.Deps{
    IDs:       idGenerator,   // model.IDGenerator
    Publisher: eventPublisher, // events.Publisher, nil desactiva la publicación
})
cat.MountOps(opRegistry) // router.OpRegistry
view := cat.NewView(caller) // router.Caller -> view.Presenter
// Catalog implementa la interfaz CatalogReader para otros módulos
appt := appointmentbooking.New(db, appointmentbooking.Deps{Catalog: cat, ...})
```

### Nota de Base de Datos
No hay datos de producción. Las bases de datos de desarrollo se recrean: `Migrate` crea `specialty`, `service` y `service_agreement` (en ese orden de dependencias FK). Nunca descarta las tablas antiguas `catalog_item` / `catalog_agreement`.

## Receta de insumos

Un servicio tiene una **receta**: los insumos que usa una atención (`service_supply`: servicio, `item_id` de `supply_manager` como una
referencia blanda, cantidad > 0, cada insumo una vez por servicio). `ServiceRecipe(tenantId, serviceId) (itemIds []string,
quantities []float64, err error)` la devuelve con tipos simples, así que `clinical_encounter` la declara en una interfaz estrecha propia
y la aplicación le pasa este módulo. Este módulo nunca importa `supply_manager`; el selector de insumos de la
pantalla del servicio viene con las pantallas que siguen las reglas de guía del módulo.
