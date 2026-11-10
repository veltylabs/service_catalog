# Arquitectura del Catálogo de Servicios

## Alcance del Dominio
El módulo `service-catalog` gestiona el catálogo de servicios y sus especialidades para el ecosistema Velty (los productos y suministros viven en `pharmacy_manager` y `supply_manager`). Sirve como fuente de verdad para otros módulos, como `appointment-booking`.

## Entidades
- **Specialty** (Especialidad): taxonomía con un código de referencia de dos letras (prefijo de los SKU).
- **Service** (Servicio): Un servicio que la empresa proporciona. Referencia a un `Specialty` mediante `specialty_id`. Contiene SKU, nombre, descripción, precio, moneda y un `tax_treatment` opcional (`TAXED` o `EXEMPT`; la tasa no se almacena aquí).
- **Agreement** (Convenio): Representa un convenio de facturación/aseguradora asociado a un servicio. Un servicio puede tener múltiples convenios.

## Patrones
- **Arnés de módulo reutilizable**: el módulo está acoplado solo a contratos publicados, no a infraestructura concreta:
    - La persistencia está manejada exclusivamente por el framework `crud` mediante las declaraciones de modelo y tienda compartida (`store.go`).
    - El módulo ya NO escribe lógicas transaccionales a mano: solo declara reglas de negocio (`rule_*.go`).
    - Las tablas se crean y migran usando `ddl.Sync`.
    - Las pruebas se ejecutan contra `storage/mem`, verificando reglas sobre un loopback (vía `crud.LoopbackActingAs`).

## Operaciones
Las operaciones son automáticamente definidas según el esquema `Allow` en `definition.go` y expuestas con los esquemas estándar del framework crud:

- `<entidad>.read`
- `<entidad>.create`
- `<entidad>.update`
- `<entidad>.delete`

Cada escritura emite el evento correspondiente (`<entidad>.created`, `<entidad>.updated`, `<entidad>.deleted`) gracias a la integración estándar de `crud` con el publicador inyectado en la inicialización (`events.Publisher`).

### Nota de Base de Datos
No hay datos de producción. Las bases de datos de desarrollo se recrean y migran usando `ddl.Step`.
