package servicecatalog

import (
	"webtyp.com/crud"
	"webtyp.com/model"
)

func (p *CurrencyOption) Read(a crud.CurrentUser, filter *CurrencyOption, args crud.ReadArgs, m *Module) (model.ModelSlice, error) {
	var list CurrencyOptionList
	for _, code := range m.currencies {
		list = append(list, &CurrencyOption{Code: code})
	}
	return &list, nil
}
