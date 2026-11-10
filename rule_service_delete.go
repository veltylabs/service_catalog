package servicecatalog

import "webtyp.com/crud"

func (p *Service) Delete(a crud.CurrentUser, tx crud.Tx, m *Module) (crud.RuleAnswer, error) {
	return crud.ApplyAndPublish, nil
}
