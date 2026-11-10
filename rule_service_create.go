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
		return 0, ErrUnsupportedCurrency
	}
	if p.TaxTreatment != TaxTaxed && p.TaxTreatment != TaxExempt {
		return 0, &model.FieldError{Field: "tax_treatment", Err: fmt.Err("invalid tax treatment")}
	}
	if err := checkPrice(*p); err != nil {
		return 0, err
	}
	conflict := false
	err := tx.ReadAll(func() model.Model { return &Service{} }, crud.Query{Conditions: []storage.Condition{orm.Eq(Service_.Sku, p.Sku)}}, func(row model.Model) {
		conflict = true
	})
	if err != nil {
		return 0, err
	}
	if conflict {
		return 0, ErrAlreadyExists
	}
	return crud.ApplyAndPublish, nil
}
