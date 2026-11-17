package servicecatalog

import (
	"webtyp.com/crud"
	"webtyp.com/model"
)

func (next *Specialty) Update(a crud.CurrentUser, tx crud.Tx, old *Specialty, fields []string, m *Module) (crud.RuleAnswer, error) {
	conflict := false
	err := tx.ReadAll(func() model.Model { return &Specialty{} }, crud.Query{}, func(row model.Model) {
		if row.(*Specialty).Prefix == next.Prefix && row.(*Specialty).Id != old.Id {
			conflict = true
		}
	})
	if err != nil {
		return crud.RuleAnswer{}, err
	}
	if conflict {
		return crud.RuleAnswer{}, ErrSpecialtyPrefixExists
	}
	return crud.ApplyAndPublish(), nil
}
