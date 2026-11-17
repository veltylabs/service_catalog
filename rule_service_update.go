package servicecatalog

import (
	"webtyp.com/crud"
	"webtyp.com/fmt"
	"webtyp.com/model"
	"webtyp.com/storage"
	"webtyp.com/orm"
)

func (next *Service) Update(a crud.CurrentUser, tx crud.Tx, old *Service, fields []string, m *Module) (crud.RuleAnswer, error) {
	validCurrency := false
	for _, c := range m.currencies {
		if next.Currency == c {
			validCurrency = true
			break
		}
	}
	if !validCurrency {
		return crud.RuleAnswer{}, ErrUnsupportedCurrency
	}
	if next.TaxTreatment != TaxTaxed && next.TaxTreatment != TaxExempt {
		return crud.RuleAnswer{}, &model.FieldError{Field: "tax_treatment", Err: fmt.Err("invalid tax treatment")}
	}
	if err := checkPrice(*next); err != nil {
		return crud.RuleAnswer{}, err
	}
	conflict := false
	err := tx.ReadAll(func() model.Model { return &Service{} }, crud.Query{Conditions: []storage.Condition{orm.Eq(Service_.Sku, next.Sku)}}, func(row model.Model) {
		if row.(*Service).Id != old.Id {
			conflict = true
		}
	})
	if err != nil {
		return crud.RuleAnswer{}, err
	}
	if conflict {
		return crud.RuleAnswer{}, ErrAlreadyExists
	}
	return crud.ApplyAndPublish(), nil
}
