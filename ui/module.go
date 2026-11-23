package ui

import "webtyp.com/svg"

// ID es la identidad de este módulo: el prefijo del recurso RBAC en el servidor,
// la ruta de navegación en el cliente. "service_catalog" es el nombre que el resto de la
// aplicación ya utiliza (RBAC, nav, tests), independiente del nombre propio de este paquete.
const ID = "service_catalog"

// Identidad de cada pantalla, para que la app las ponga en su menú.
//
// Ejemplo: Screens: []string{ui.ServicesID, ui.SpecialtiesID}
const (
	ServicesID    = ID + ".services"
	SpecialtiesID = ID + ".specialties"
)

// Label es el texto a mostrar en el ítem de navegación.
const Label = "Catálogo"

const (
	IconCatalog     = svg.Icon("service_catalog")
	IconServices    = svg.Icon("service_catalog_services")
	IconSpecialties = svg.Icon("service_catalog_specialties")
)
