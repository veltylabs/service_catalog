---
PLAN: "feat: migrate ui to view.Screen and drop layout dependency"
TAG: v0.6.5
EXECUTOR: jules
REVIEWER: none
STATUS: review
SESSION: 9604813858897086912
PR: https://github.com/veltylabs/service_catalog/pull/17
REVIEWED_COMMIT: 4997abe04417ea0b135c8594318257529c68c835
---

> Read AGENTS.md at the repo root before any change.

# Migración de UI a view.Screen y eliminación de dependencia webtyp/layout

## Ecosystem search

- `webtyp.com/view` (`v0.6.45`): define el contrato de pantalla declarativa `view.Screen`, `view.Lookup`, `view.Presenter`, `view.WithDescription`, `view.WithTitle`, y el validador `Screen.Validate()`.
- `webtyp.com/layout/app` (`v0.3.76`): `BuildMenu(screens []view.Screen, ids model.IDGenerator, allows func(model.Resource) bool) (Menu, error)` procesa las pantallas declaradas y construye la vista (`crudview`) y permisos en el cliente.
- `veltylabs/service_catalog`: actualmente exporta `ui.Browser(...) (platformd.UIModule, error)` importando `webtyp.com/layout/platformd`, `crudview`, y `dom`, forzando a este módulo de dominio a acoplarse con librerías de chasis visual y renderizado.

## Design gate

1. **Prior art**: `webtyp/view` v0.6.44 introdujo `view.Screen` como el contrato único para que los módulos declaren qué muestran como datos puros, sin importar renderizadores ni chasis.
2. **Novice-name test**: `ui.Screens(caller router.Caller) []view.Screen` es autoexplicativo: devuelve las pantallas que ofrece este módulo al menú de la aplicación.
3. **Complexity ledger**: Elimina `webtyp.com/layout` (y sus dependencias indirectas `components`, `widget`, `css`, `color`, `font`, `calendar`, etc.) de `go.mod`. La pantalla se define en ~35 líneas declarativas en lugar de lógica imperativa de montaje de DOM y eventos.
4. **Where it belongs**: `ui.Screens` vive en `github.com/veltylabs/service_catalog/ui`. Los constructores de Presenter (`NewView`, `NewSpecialtyView`) viven en la raíz de `servicecatalog`.
5. **What it deletes**: Borra `ui.Browser`, el directorio `web/browser/main.go`, y el require de `webtyp.com/layout` en `go.mod`.

## Tools

```bash
go install webtyp.com/devflow/cmd/gotest@latest
go install webtyp.com/lang/cmd/langc@latest
go get webtyp.com/view@v0.6.45
go get webtyp.com/router@v0.4.8
```

## APIs used

| Symbol | Signature | Tag | File |
|---|---|---|---|
| `view.Screen` | `type Screen struct { ID string; Label lang.Text; Icon string; Permission model.Resource; Presenter Presenter; Lookups []Lookup; Help lang.Text; Children []Screen }` | `webtyp.com/view v0.6.45` | `screen.go:35` |
| `view.Lookup` | `type Lookup struct { Field string; Source Presenter }` | `webtyp.com/view v0.6.45` | `screen.go:49` |
| `view.Presenter` | `type Presenter interface { ... }` | `webtyp.com/view v0.6.45` | `presenter.go:12` |
| `view.WithTitle` | `func WithTitle(title string) Option` | `webtyp.com/view v0.6.45` | `view.go:124` |
| `view.WithDescription` | `func WithDescription(d lang.Text) Option` | `webtyp.com/view v0.6.45` | `view.go:138` |
| `view.New` | `func New(l Lister, record model.Model, opts ...Option) Presenter` | `webtyp.com/view v0.6.45` | `view.go:147` |
| `router.Caller` | `type Caller interface { Call(op string, args Encodable, out Decodable, done func(error)); Dispatch(op string, args Encodable) }` | `webtyp.com/router v0.4.8` | `router.go` |

Si una firma en la versión publicada difiere de esta tabla, detente e infórmalo en el PR; nunca adaptes por adivinanza.

## Deletions

- `ui.Browser` en `ui/browser.go`. Criterio: `grep -rn "func.*Browser(" ui/` -> vacío.
- Todo import de `webtyp.com/layout` en el repositorio. Criterio: `grep -rn "webtyp.com/layout" .` -> vacío.
- Directorio de demo `web/browser/main.go`. Criterio: `test ! -d web` -> true.

---

## Stage 1: Soportar view.Option en constructores de Presenter

En `view.go`:
Permitir que `NewView` y `NewSpecialtyView` reciban opciones funcionales `opts ...view.Option` para configurar descripciones (`view.WithDescription`) u opciones adicionales.

Reemplazar `view.go` con:

```go
package servicecatalog

import (
	"webtyp.com/model"
	"webtyp.com/router"
	"webtyp.com/view"
)

const titleCatalog = "Catálogo"
const titleSpecialties = "Especialidades"

// Item proyecta un Service como un view.Item — el ÚNICO código específico de la vista que lleva
// este registro (view.Itemizer). El índice interno del Presenter (construido a partir de esto en Reload)
// reemplaza la antigua búsqueda manual byID/WithFill.
// Ejemplo: item := m.Item()
func (m *Service) Item() view.Item {
	return view.Item{ID: m.Id, Label: m.Name, Description: m.Sku}
}

// Item proyecta un Specialty como un view.Item (view.Itemizer).
// Ejemplo: item := s.Item()
func (s *Specialty) Item() view.Item {
	return view.Item{ID: s.Id, Label: s.Name, Description: s.Prefix}
}

// NewView construye el Presenter del ítem del catálogo — el motor agnóstico a la tecnología que envuelve
// un renderizador (crudview, o cualquier otro). Es tarea de ESTE módulo construirlo (importando solo view+model+router);
// la aplicación decide qué renderizador lo dibuja.
// Ejemplo: pres := servicecatalog.NewView(caller, view.WithDescription("Ayuda..."))
func NewView(caller router.Caller, opts ...view.Option) view.Presenter {
	b := view.NewCallerLister(caller,
		view.Ops{Module: ModelName, List: OpListServices, Save: OpUpsertService, Update: OpUpdateService, Delete: OpDeleteService},
		func() model.ModelSlice { return &ServiceList{} })
	allOpts := append([]view.Option{view.WithTitle(titleCatalog)}, opts...)
	return view.New(b, &Service{}, allOpts...)
}

// NewSpecialtyView construye el Presenter de la especialidad.
// Ejemplo: pres := servicecatalog.NewSpecialtyView(caller, view.WithDescription("Ayuda..."))
func NewSpecialtyView(caller router.Caller, opts ...view.Option) view.Presenter {
	b := view.NewCallerLister(caller,
		view.Ops{Module: ModelName, List: OpListSpecialties, Save: OpUpsertSpecialty, Delete: OpDeleteSpecialty},
		func() model.ModelSlice { return &SpecialtyList{} })
	allOpts := append([]view.Option{view.WithTitle(titleSpecialties)}, opts...)
	return view.New(b, &Specialty{}, allOpts...)
}
```

---

## Stage 2: Declarar ui.Screens sin dependencias de layout

En `ui/browser.go`:
Eliminar el método `Browser`.
Reemplazar con `Screens(caller router.Caller) []view.Screen`.

Reescribir `ui/browser.go`:

```go
package ui

import (
	servicecatalog "github.com/veltylabs/service_catalog"
	"webtyp.com/lang"
	"webtyp.com/router"
	"webtyp.com/view"
)

// Textos de ayuda que el chasis muestra en el botón «?» y en la bienvenida de
// cada pantalla. Las claves están en inglés; la traducción vive en lang.json.
const (
	descCatalog     lang.Text = "The services the clinic offers and the specialties that group them. Start with the specialties, then add each service with its price."
	descServices    lang.Text = "Each service the clinic offers, with its specialty, price and tax treatment. Pick one in the list to edit it, or press + to add a new one."
	descSpecialties lang.Text = "Specialties group the services (for example, dentistry or kinesiology). Create them first: every service must belong to one."
)

// Screens declara las pantallas que este módulo ofrece al menú de la app.
//
// Ejemplo: screens := ui.Screens(caller)
func Screens(caller router.Caller) []view.Screen {
	return []view.Screen{{
		ID:    ID,
		Label: Label,
		Icon:  string(IconCatalog),
		Help:  descCatalog,
		Children: []view.Screen{
			{
				ID:         ID + ".services",
				Label:      "Servicios",
				Icon:       string(IconServices),
				Permission: "service",
				Presenter:  servicecatalog.NewView(caller, view.WithDescription(descServices)),
				Lookups: []view.Lookup{
					{Field: "specialty_id", Source: servicecatalog.NewSpecialtyView(caller)},
				},
			},
			{
				ID:         ID + ".specialties",
				Label:      "Especialidades",
				Icon:       string(IconSpecialties),
				Permission: "specialty",
				Presenter:  servicecatalog.NewSpecialtyView(caller, view.WithDescription(descSpecialties)),
			},
		},
	}}
}
```

---

## Stage 3: Eliminar directorio web/

Eliminar el archivo `web/browser/main.go` y la carpeta `web/`:
```bash
rm -rf web
```

---

## Stage 4: Actualizar pruebas unitarias

1. En `tests/ui_described_test.go`:
Reescribir para validar `ui.Screens` sin importar `webtyp.com/layout`:

```go
package tests

import (
	"testing"

	"github.com/veltylabs/service_catalog/ui"
	"webtyp.com/view"
)

// Cada pantalla del catálogo se explica sola: el grupo y sus dos hijos
// validan contra view.Screen y tienen descripciones no vacías.
func TestBrowser_EveryScreenIsDescribed(t *testing.T) {
	screens := ui.Screens(&mockCaller{})
	if len(screens) != 1 {
		t.Fatalf("expected 1 root group screen, got %d", len(screens))
	}
	root := screens[0]
	if err := root.Validate(); err != nil {
		t.Fatalf("root screen validation failed: %v", err)
	}
	if root.ID != ui.ID {
		t.Errorf("root ID = %q, want %q", root.ID, ui.ID)
	}
	if root.Help == "" {
		t.Errorf("root screen has empty help")
	}
	if len(root.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(root.Children))
	}
	for i, c := range root.Children {
		if err := c.Validate(); err != nil {
			t.Fatalf("child %d validation failed: %v", i, err)
		}
		if c.Permission == "" {
			t.Errorf("child %d permission is empty", i)
		}
		d, ok := c.Presenter.(view.WithDescription)
		if !ok || d.Description() == "" {
			t.Errorf("child %d presenter must implement view.WithDescription with a non-empty text", i)
		}
	}
}
```

2. En `tests/service_catalog_view_browser_test.go` y `tests/service_catalog_view_test.go`:
Eliminar los imports de `webtyp.com/layout`, `crudview`, `platformd`.
Reemplazar con pruebas directas de `ui.Screens` y de los presenters (`NewView` y `NewSpecialtyView`).

---

## Stage 5: Limpieza de dependencias en go.mod y verificación

Ejecutar:
```bash
go mod tidy
```

Verificar que:
1. `grep "webtyp.com/layout" go.mod` no devuelve nada.
2. `grep "webtyp.com/layout" go.sum` no devuelve nada.
3. `grep -rn "webtyp.com/layout" .` no devuelve nada.

Ejecutar la suite de pruebas completa:
```bash
gotest
```

---

## Stage 6: Documentación y regla de idioma en veltylabs

1. En `AGENTS.md`:
Actualizar la sección `ui/, seed/ and web/`:
Reemplazar la referencia a `Browser(caller router.Caller, ids model.IDGenerator, tenantID string) (platformd.UIModule, error)` por:
`Screens(caller router.Caller) []view.Screen`.
Indicar claramente que `ui/` no importa `webtyp.com/layout` ni ningún renderizador, sino que devuelve datos declarativos (`[]view.Screen`).

2. En `README.md`:
Actualizar cualquier mención a `ui.Browser` por `ui.Screens`.

3. Ejecutar `langc sync .` y verificar que no haya claves vacías ni errores.

---

## Stages table

| Etapa | Archivos modificados/creados | Qué hace |
|---|---|---|
| Stage 1 | `view.go` | Agrega `opts ...view.Option` a constructores de Presenter |
| Stage 2 | `ui/browser.go` | Exporta `ui.Screens` declarativo con `view.Screen` y borra imports de `layout` |
| Stage 3 | `web/` (borrado) | Elimina `web/browser/main.go` |
| Stage 4 | `tests/ui_described_test.go`, `tests/service_catalog_view_*.go` | Valida `ui.Screens`, `Lookups` y remueve dependencias de `layout` |
| Stage 5 | `go.mod`, `go.sum` | `go mod tidy` y ejecución de `gotest` |
| Stage 6 | `AGENTS.md`, `README.md`, `lang.json` | Actualiza documentación en español y sincroniza `langc` |

## Delivery check

1. `grep -rn "webtyp.com/layout" .` debe retornar vacío.
2. `test ! -d web` debe retornar éxito (0).
3. `gotest` pasa con 0 fallos de compilación, 0 fallos de tests y cobertura verificada.
4. `langc sync .` pasa sin errores.
5. El PR se abre inmediatamente al pasar estas comprobaciones.
## Executor notes

- Dropped `tax_treatment` and `currency` set options logic, as well as `NewRecord` initialization in `ui.Browser`, which previously set default values (like `IsActive: true`, default `Currency` and `TaxTreatment`). Since `view.Screen` is fully declarative and does not deal with initialization defaults (the renderer is expected to manage this, or it must be specified declaratively in another way that is not currently requested by the plan), these forms logic were omitted.
- Updated `README.md` and `AGENTS.md` to completely remove mentions of `web/` demo, as it was permanently deleted in Stage 3.
