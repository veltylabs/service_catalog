package ui

import (
	servicecatalog "github.com/veltylabs/service_catalog"
	"webtyp.com/lang"
	"webtyp.com/model"
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
				Permission: model.ResourceOf(&servicecatalog.ServiceModel),
				Presenter:  servicecatalog.NewView(caller, view.WithDescription(descServices)),
			},
			{
				ID:         ID + ".specialties",
				Label:      "Especialidades",
				Icon:       string(IconSpecialties),
				Permission: model.ResourceOf(&servicecatalog.SpecialtyModel),
				Presenter:  servicecatalog.NewSpecialtyView(caller, view.WithDescription(descSpecialties)),
			},
		},
	}}
}
