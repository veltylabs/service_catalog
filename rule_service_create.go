package servicecatalog

import (
	"webtyp.com/crud"
	"webtyp.com/fmt"
	"webtyp.com/model"
	"webtyp.com/storage"
	"webtyp.com/orm"
)

func (p *Service) Create(a crud.CurrentUser, tx crud.Tx, m *Module) (crud.RuleAnswer, error) {
	validCurrency := false
	for _, c := range m.currencies {
		if p.Currency == c {
			validCurrency = true
			break
		}
	}
	if !validCurrency {
		return crud.RuleAnswer{}, ErrUnsupportedCurrency
	}
	if p.TaxTreatment != TaxTaxed && p.TaxTreatment != TaxExempt {
		return crud.RuleAnswer{}, &model.FieldError{Field: "tax_treatment", Err: fmt.Err("invalid tax treatment")}
	}
	if err := checkPrice(*p); err != nil {
		return crud.RuleAnswer{}, err
	}
	conflict := false
	err := tx.ReadAll(func() model.Model { return &Service{} }, crud.Query{Conditions: []storage.Condition{orm.Eq(Service_.Sku, p.Sku)}}, func(row model.Model) {
		conflict = true
	})
	if err != nil {
		return crud.RuleAnswer{}, err
	}
	if conflict {
		return crud.RuleAnswer{}, ErrAlreadyExists
	}
	return crud.ApplyAndPublish(), nil
}
