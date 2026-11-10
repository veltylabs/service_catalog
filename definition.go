package servicecatalog

import (
	"webtyp.com/input"
	"webtyp.com/model"
)

// Helpers para satisfacer el parser de tipos estático AST de ormc v0.9.24:
var (
	Checkbox_FieldBool = input.Checkbox()
	Decimal_FieldFloat = input.Decimal()

	BaseBool_FieldBool = model.Bool()
	BaseInt_FieldInt   = model.Int()
)

// CurrencyCLP es la moneda por defecto del catálogo.
const CurrencyCLP = "CLP"

const (
	TaxTaxed  = "TAXED"
	TaxExempt = "EXEMPT"
)

var SpecialtyModel = model.Definition{
	Name: "specialty",
	Allow: model.Create | model.Read | model.Update | model.Delete,
	Rules: model.Rules{Create: model.AddMyBusinessRule, Update: model.AddMyBusinessRule, Delete: model.AddMyBusinessRule, Read: model.HaveNoBusinessRule},
	Fields: model.Fields{
		{Name: "id", Type: model.Text(), DB: &model.FieldDB{PK: true}, OmitEmpty: true, ServerOnly: true},
		{Name: "tenant_id", Type: model.Text(), NotNull: true, ServerOnly: true},
		{Name: "prefix", Type: input.Text(), Label: "Reference code", Help: "Two letters every service code (SKU) of this specialty starts with, e.g. MD for MD001.", NotNull: true, Item: model.ItemDescription,
			Permitted: model.Permitted{Letters: true, Minimum: 2, Maximum: 2}},
		// input.Text()'s own baseline Extra ('.', ',', '(', ')') no incluye
		// '+' — "Dental (Odontología + Ortodoncia)" lo necesita, así que
		// este campo declara su propio whitelist para reemplazar el piso
		// del Kind (Field.Validate: un whitelist positivo propio reemplaza
		// el floor de Text, no lo suma). Tilde cubre las vocales acentuadas
		// sin listarlas una a una.
		{Name: "name", Type: input.Text(), Label: "Name", NotNull: true, Item: model.ItemLabel,
			Permitted: model.Permitted{Letters: true, Tilde: true, Spaces: true, Extra: []rune{'+', '(', ')'}, Minimum: 1, Maximum: 100}},
		{Name: "description", Type: input.Textarea(), Label: "Description", OmitEmpty: true},
		{Name: "position", Type: model.Int(), OmitEmpty: true},
		// NotNull deliberadamente ausente: para un bool, NotNull rechaza el
		// zero value (ValidateFields: field.NotNull && IsZeroPtr, y
		// IsZeroPtr para FieldBool es !*p) — con NotNull:true sería
		// imposible crear una especialidad SIN publicar, justo lo opuesto
		// de "cerrado por defecto". El false-por-defecto ya lo garantiza
		// el zero value de Go, no una restricción de esquema.
		{Name: "is_published", Type: input.Checkbox(), Label: "Published"},
		{Name: "updated_at", Type: model.NanoTime(), OmitEmpty: true, ServerOnly: true},
	},
}

var ServiceModel = model.Definition{
	Name: "service",
	Allow: model.Create | model.Read | model.Update | model.Delete,
	Rules: model.Rules{Create: model.AddMyBusinessRule, Update: model.AddMyBusinessRule, Delete: model.AddMyBusinessRule, Read: model.AddMyBusinessRule},
	Fields: model.Fields{
		{Name: "id", Type: model.Text(), DB: &model.FieldDB{PK: true}, OmitEmpty: true, ServerOnly: true},
		{Name: "tenant_id", Type: model.Text(), NotNull: true, ServerOnly: true},
		{Name: "specialty_id", Type: input.Select(), Label: "Specialty", Ref: &SpecialtyModel, DB: &model.FieldDB{RefColumn: "id"}, NotNull: true},
		{Name: "sku", Type: input.Text(), Label: "SKU", NotNull: true, Item: model.ItemDescription, Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune{'-'}, Minimum: 1, Maximum: 50}},
		// Nombres clínicos reales llevan guion y barra: "Tórax simple (AP-PA)", "Columna lumbar AP/LAT".
		{Name: "name", Type: input.Text(), Label: "Name", NotNull: true, Item: model.ItemLabel,
			Permitted: model.Permitted{Letters: true, Tilde: true, Numbers: true, Spaces: true,
				Extra: []rune{'.', ',', '(', ')', '-', '/', '+'}, Minimum: 1, Maximum: 255}},
		{Name: "description", Type: input.Textarea(), Label: "Description", Help: "Public information: patients read it when booking. Say what the service includes, how long it takes or how to prepare.", OmitEmpty: true},
		{Name: "tax_treatment", Type: input.Select(), Label: "Tax treatment", Help: "Exempt: a health service without IVA (consultations and most exams). Taxed: it carries 19% IVA.", NotNull: true},
		// Sin NotNull: 0 es válido (precio pendiente). La regla real la aplica checkPrice: un servicio activo necesita precio.
		{Name: "price", Type: input.Decimal(), Label: "Price"},
		{Name: "currency", Type: input.Select(), Label: "Currency", NotNull: true, Permitted: model.Permitted{Letters: true, Minimum: 3, Maximum: 3}},
		// Sin NotNull: para un bool rechazaría false y no se podría guardar un servicio inactivo (ver is_published).
		{Name: "is_active", Type: input.Checkbox(), Label: "Active", Help: "Only active services can be booked and charged. Turn it off to stop offering it without losing its history."},
		{Name: "updated_at", Type: model.NanoTime(), OmitEmpty: true, ServerOnly: true},
		{Name: "active_only", Type: model.Bool(), OmitEmpty: true, DB: nil},
		{Name: "limit", Type: model.Int(), OmitEmpty: true, DB: nil},
		{Name: "offset", Type: model.Int(), OmitEmpty: true, DB: nil},
		{Name: "search_text", Type: model.Text(), OmitEmpty: true, DB: nil},
		{Name: "after", Type: model.Text(), OmitEmpty: true, DB: nil},
	},
}

var ServiceAgreementModel = model.Definition{
	Name: "service_agreement",
	Allow: model.Create | model.Read | model.Update | model.Delete,
	Rules: model.Rules{Create: model.AddMyBusinessRule, Update: model.AddMyBusinessRule, Delete: model.HaveNoBusinessRule, Read: model.HaveNoBusinessRule},
	Fields: model.Fields{
		{Name: "id", Type: model.Text(), DB: &model.FieldDB{PK: true}, OmitEmpty: true, ServerOnly: true},
		{Name: "tenant_id", Type: model.Text(), NotNull: true, ServerOnly: true},
		{Name: "service_id", Type: model.Text(), Ref: &ServiceModel, DB: &model.FieldDB{RefColumn: "id"}, NotNull: true},
		{Name: "insurer", Type: input.Text(), Label: "Insurer", NotNull: true, Permitted: model.Permitted{Letters: true, Spaces: true, Minimum: 1, Maximum: 100}},
		{Name: "code", Type: input.Text(), Label: "Code", OmitEmpty: true, Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune{'-'}, Minimum: 1, Maximum: 50}},
		{Name: "price", Type: input.Decimal(), Label: "Price", OmitEmpty: true},
		{Name: "is_active", Type: input.Checkbox(), Label: "Active"}, // sin NotNull: rechazaría false (ver is_published)
		{Name: "updated_at", Type: model.NanoTime(), OmitEmpty: true, ServerOnly: true},
	},
}

// ServiceSupplyModel es una línea de la receta de insumos de un servicio: lo que usa una atención. item_id es una
// referencia suave a un ítem de supply_manager (este módulo nunca lo importa); quantity está en la unidad del ítem.
var ServiceSupplyModel = model.Definition{
	Name: "service_supply",
	Allow: model.Create | model.Read | model.Update | model.Delete,
	Rules: model.Rules{Create: model.AddMyBusinessRule, Update: model.AddMyBusinessRule, Delete: model.HaveNoBusinessRule, Read: model.HaveNoBusinessRule},
	Fields: model.Fields{
		{Name: "id", Type: model.Text(), DB: &model.FieldDB{PK: true}, OmitEmpty: true, ServerOnly: true},
		{Name: "tenant_id", Type: model.Text(), NotNull: true, ServerOnly: true},
		{Name: "service_id", Type: model.Text(), Ref: &ServiceModel, DB: &model.FieldDB{RefColumn: "id"}, NotNull: true},
		{Name: "item_id", Type: model.Text(), NotNull: true},
		{Name: "quantity", Type: input.Decimal(), Label: "Quantity", NotNull: true},
		{Name: "updated_at", Type: model.NanoTime(), OmitEmpty: true, ServerOnly: true},
	},
}

var CurrencyOptionModel = model.Definition{
	Name: "currency_option",
	Allow: model.Read,
	Rules: model.Rules{Read: model.AddMyBusinessRule},
	Fields: model.Fields{
		{Name: "code", Type: model.Text()},
	},
}


var Definition = model.ModuleDefinition{
	Name: "service_catalog",
	Entities: []model.Definition{SpecialtyModel, ServiceModel, ServiceAgreementModel, ServiceSupplyModel, CurrencyOptionModel},
}
