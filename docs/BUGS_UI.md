# Bugs de UI y mejoras — `service_catalog`

> Auditoría en vivo sobre `#service_catalog.services` y `#service_catalog.specialties` (2026-11-19).

## Bugs encontrados

| # | Sev | Pantalla | Qué se ve en pantalla | Causa raíz verificada |
|---|---|---|---|---|
| SC-1 | 🟡 | `#service_catalog.services` | **El selector «Tratamiento tributario» se muestra en blanco** al abrir registros importados como «Consulta general» (en vez de mostrar «Selecciona una opción» o su valor). | El registro en BD tiene `tax_treatment = ""` y el `<select>` renderizado en `form` cuando recibe un valor vacío tras haber renderizado o sin coincidencia muestra la caja en blanco sin placeholder. |
| SC-2 | 🟡 | `#service_catalog.services` | **Lista de 117 servicios sin filtro por «Especialidad».** Para encontrar o auditar los servicios de una especialidad hay que recorrer 117 tarjetas o adivinar por nombre. | `ServiceModel` (`service_catalog/definition.go:62`) no declara `Filter: true` en `specialty_id` ni `crudview` tiene filtro de especialidad montado (ver U3 de `APP_USABILITY_MASTER_PLAN.md`). |
| SC-3 | 🟡 | Menú | **Claves de pantalla en español literal en Go** (`Label: "Servicios"`, `Label: "Especialidades"`). | `service_catalog/ui/browser.go:31,38` rompe la regla de `translations` (claves en inglés en Go, traducción en `lang.json`), impidiendo traducción automática a otros idiomas. |
| SC-4 | ⚪ | `#service_catalog.services` | **En la lista de 117 servicios no se muestra el precio ni la especialidad en la fila**, solo el nombre y el SKU (`mdcong`). | `ServiceModel` solo marca `name` (`ItemLabel`) y `sku` (`ItemDescription`), ocultando precio y estado activo/inactivo hasta abrir uno por uno. |

## Sugerencias de mejora justificadas

1. **Agregar filtro por Especialidad sobre la lista de Servicios:** montar `selectsearch` de especialidades arriba del buscador (patrón U3). *Justificación:* con 117+ prestaciones clínicas, filtrar por especialidad es la forma natural en que administración revisa precios y vigencia.
2. **Mostrar precio y estado inactivo en la tarjeta de lista (`Item`):** incluir precio formateado (p. ej. `$30.000`) y una marca sutil si está inactivo. *Justificación:* evita abrir decenas de fichas para detectar servicios sin precio (`Precio: 0`, como ocurre hoy en «Consulta general») o desactivados.
3. **Corregir las claves de `ui/browser.go` a inglés (`"Services"`, `"Specialties"`)** y completar el tratamiento tributario por defecto o mostrar la opción prompt en el `<select>`. *Justificación:* elimina errores de `langc` y evita selects visualmente vacíos.
