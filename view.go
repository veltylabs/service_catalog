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
