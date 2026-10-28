package servicecatalog

import (
	"webtyp.com/model"
	"webtyp.com/router"
	"webtyp.com/view"
)

const titleCatalog = "Catálogo"
const titleSpecialties = "Especialidades"

// NewView construye el Presenter del ítem del catálogo — el motor agnóstico a la tecnología que envuelve
// un renderizador (crudview, o cualquier otro). Es tarea de ESTE módulo construirlo (importando solo view+model+router);
// la aplicación decide qué renderizador lo dibuja.
// Ejemplo: pres := servicecatalog.NewView(caller, view.WithDescription("Ayuda..."))
func NewView(caller router.Caller, opts ...view.Option) view.Presenter {
	b := view.NewCallerLister(view.NewRecords().Caller(caller),
		view.Ops{Module: ModelName, List: OpListServices, Save: OpUpsertService, Update: OpUpdateService, Delete: OpDeleteService},
		func() model.ModelSlice { return &ServiceList{} })
	allOpts := append([]view.Option{view.WithTitle(titleCatalog)}, opts...)
	return view.New(b, &Service{}, allOpts...)
}

// NewSpecialtyView construye el Presenter de la especialidad.
// Ejemplo: pres := servicecatalog.NewSpecialtyView(caller, view.WithDescription("Ayuda..."))
func NewSpecialtyView(caller router.Caller, opts ...view.Option) view.Presenter {
	b := view.NewCallerLister(view.NewRecords().Caller(caller),
		view.Ops{Module: ModelName, List: OpListSpecialties, Save: OpUpsertSpecialty, Delete: OpDeleteSpecialty},
		func() model.ModelSlice { return &SpecialtyList{} })
	allOpts := append([]view.Option{view.WithTitle(titleSpecialties)}, opts...)
	return view.New(b, &Specialty{}, allOpts...)
}
