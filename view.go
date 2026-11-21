package servicecatalog

import (
	"webtyp.com/fmt"
	"webtyp.com/model"
	"webtyp.com/router"
	"webtyp.com/view"
)

// Item construye la fila visual de una especialidad con el color asignado desde la paleta en BD.
func (m *Specialty) Item() model.Item {
	return model.Item{ID: m.Id, Label: m.Name, Description: m.Prefix, Tone: m.Tone}
}

const titleCatalog = "Catálogo"
const titleSpecialties = "Especialidades"

// NewView construye el Presenter del ítem del catálogo — el motor agnóstico a la tecnología que envuelve
// un renderizador (crudview, o cualquier otro). Es tarea de ESTE módulo construirlo (importando solo view+model+router);
// la aplicación decide qué renderizador lo dibuja.
// Ejemplo: pres := servicecatalog.NewView(caller, view.WithDescription("Ayuda..."))
func NewView(caller router.Caller, opts ...view.Option) view.Presenter {
	b := view.NewCallerLister(view.NewRecords().Caller(caller),
		view.Ops{Module: ModelName, Page: view.Page{First: view.All}},
		func() model.TableSlice { return &ServiceList{} })
	allOpts := append([]view.Option{
		view.WithTitle(titleCatalog),
		view.WithOptions(Service_.TaxTreatment,
			fmt.KeyValue{Key: TaxExempt, Value: "Exento"},
			fmt.KeyValue{Key: TaxTaxed, Value: "Afecto (19% IVA)"},
		),
		view.WithOptions(Service_.Currency,
			fmt.KeyValue{Key: CurrencyCLP, Value: "CLP"},
		),
	}, opts...)
	return view.New(b, &Service{}, allOpts...)
}

// NewSpecialtyView construye el Presenter de la especialidad.
// Ejemplo: pres := servicecatalog.NewSpecialtyView(caller, view.WithDescription("Ayuda..."))
func NewSpecialtyView(caller router.Caller, opts ...view.Option) view.Presenter {
	b := view.NewCallerLister(view.NewRecords().Caller(caller),
		view.Ops{Module: ModelName, Page: view.Page{First: view.All}},
		func() model.TableSlice { return &SpecialtyList{} })
	allOpts := append([]view.Option{view.WithTitle(titleSpecialties)}, opts...)
	return view.New(b, &Specialty{}, allOpts...)
}
