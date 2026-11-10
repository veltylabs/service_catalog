# Service Catalog Module
<img src="docs/img/badges.svg">

El módulo `service-catalog` gestiona el catálogo de servicios y sus especialidades para el ecosistema Velty.
Actúa como la fuente de verdad de los servicios ofertados por una clínica o centro de salud, los cuales son consumidos
por módulos como `appointment-booking`.

## Entidades y Operaciones (crud)

La persistencia y el transporte los maneja el framework `crud` bajo las reglas del dominio escritas en `rule_*.go`.
No hay lógica de persistencia manual; la tabla mapea a:

| Today | After |
|---|---|
| `list_specialties / get_specialty` | `specialty.read` |
| `upsert_specialty` | `specialty.create` / `specialty.update` |
| `delete_specialty` | `specialty.delete` |
| `list_services / get_service` | `service.read` |
| `find_service_by_sku` | `service.read` with filter `{sku}` |
| `create_service / update_service / upsert_service` | `service.create` / `service.update` |
| `deactivate_service` | `service.update` with `["is_active"]` |
| `delete_service` | `service.delete` |
| `list_currencies` | read model `CurrencyOption` |
| `list_agreements / upsert_agreement / delete_agreement` | `service_agreement.read`, `service_agreement.create`, `service_agreement.update`, `service_agreement.delete` |
| `list_service_supplies / upsert_service_supply / delete_service_supply` | `service_supply.read`, `service_supply.create`, `service_supply.update`, `service_supply.delete` |

### Sin/Con Framework

| Antes | Con crud-by-declaration |
|---|---|
| Handlers manuales (`op*.go`), `store.Create`, `store.Update` repetitivos | Declaración en `definition.go` |
| Lógica de dominio acoplada a DB | Lógica de negocio (unicidad, precio) en `rule_*.go` |
| Operaciones duplicadas por verbos | Rutas automáticas REST/MCP basadas en acciones |

## Dependencias (Hexagonal)

El módulo requiere ser inicializado mediante la inyección de `IDGenerator`, y `Events` de manera opcional para publicar sus eventos de creación y edición.
