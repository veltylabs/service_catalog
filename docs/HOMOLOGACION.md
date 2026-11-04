# Homologación de servicios — qué falló en Pa100T y cómo se evita

## El problema

En Pa100T cada profesional creaba **su propia copia** del servicio que ofrecía. Si cinco médicos hacían consultas,
existían cinco filas `service`, cada una con su código, su nombre y su duración:

| profesional | código | nombre | duración |
|---|---|---|---|
| Dra. A | `dcon` | consulta | 30 |
| Dr. B | `con` | Consulta | 20 |
| Dra. C | `consulta` | medicina general | 15 |
| Dr. D | `dcon` | consulta dental | 30 |

El resultado, con los años: 131 filas para lo que en realidad son 100 servicios, nombres con errores de tipeo
(«meddicina general»), un mismo código usado para cosas distintas (`dcon` = consulta general **y** consulta dental) y
ninguna forma de saber cuántas consultas se hicieron en total, ni de ponerles un precio común.

## La causa

Se mezclaron dos cosas distintas en una sola fila:

1. **Qué es el servicio** (una consulta general): igual para todos, con un solo nombre, un código y un precio.
2. **Cómo lo atiende cada profesional** (cuánto tarda): distinto para cada uno, y está bien que lo sea.

Como la duración era distinta, parecía que hacía falta un servicio distinto por profesional. No hacía falta.

## Cómo queda en mjosefa-cms

| Qué | Dónde | Cuántos |
|---|---|---|
| El servicio (nombre, SKU, precio, especialidad) | `service_catalog` (este módulo) — el catálogo | **uno** por servicio real |
| La duración de cada profesional para ese servicio | `appointment_booking` — `employee_service_config` | uno por profesional y servicio |

Así, la Dra. A y el Dr. B ofrecen **el mismo** servicio `mdcong` (Consulta general), cada uno con su duración (30 y
20 minutos). Se puede contar, cobrar y reportar el servicio una sola vez.

## Cómo se pasó de Pa100T al catálogo

1. Se agruparon las 131 filas por lo que realmente eran ([CATALOGO_SERVICIOS.md](https://github.com/veltylabs/mjosefa-cms/blob/main/docs/service_catalog/CATALOGO_SERVICIOS.md) (en `mjosefa-cms`, primer cliente), columna
   «código(s) origen»). Los casos dudosos quedaron marcados ⚠️ y los que nadie supo explicar, 🟡.
2. El resultado es [migrate/pastdb/catalog.csv](https://github.com/veltylabs/mjosefa-cms/blob/main/migrate/pastdb/catalog.csv): 100 servicios.
3. Cada fila antigua apunta a su servicio en `migrate/pastdb/service_map.csv` (`código,nombre → sku`), que se revisa a
   mano antes de importar: una fila sin equivalencia detiene la importación.

## La regla para el futuro

Antes de crear un servicio nuevo, buscar si ya existe en el catálogo. Si existe, al profesional se le **asigna** ese
servicio con su duración; nunca se crea una copia.
