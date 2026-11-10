# Tests
| Test | Rule | Under |
|---|---|---|
| [TestTenant_ReadOtherTenantIsNotFound](../tests/catalog_test.go) | read | TestTenant_ReadOtherTenantIsNotFound |
| [TestService_DuplicateSkuRejected](../tests/catalog_test.go) | unique sku | TestService_DuplicateSkuRejected |
| [TestService_FindBySkuFilter](../tests/catalog_test.go) | filter sku | TestService_FindBySkuFilter |
| [TestSpecialty_DeleteInUseRejected](../tests/catalog_test.go) | in use | TestSpecialty_DeleteInUseRejected |
| [TestReadModel_Currencies](../tests/currency_test.go) | read currencies | TestReadModel_Currencies |
| [TestMount_MatchesDefinition](../tests/conformance_test.go) | definition | TestMount_MatchesDefinition |
| [TestCurrencyValidation](../tests/currency_test.go) | currency validation | TestCurrencyValidation |
| [TestDemo_Valid](../tests/demo_test.go) | demo valid | TestDemo_Valid |
| [TestSentinelErrorKinds](../tests/error_test.go) | error kinds | TestSentinelErrorKinds |
| [TestSentinelErrorTexts](../tests/error_test.go) | error texts | TestSentinelErrorTexts |
| [TestItemsDeclaredOnFields](../tests/items_test.go) | field items | TestItemsDeclaredOnFields |
| [TestPricePending](../tests/price_pending_test.go) | price pending | TestPricePending |
| [TestServiceRecipe_LinesInOrder](../tests/recipe_test.go) | lines in order | TestServiceRecipe_LinesInOrder |
| [TestServiceRecipe_Rules](../tests/recipe_test.go) | recipe rules | TestServiceRecipe_Rules |
| [TestSeedLoad](../tests/seed_test.go) | seed load | TestSeedLoad |
| [TestMigrate_UpToDate_OnlyEnsuresTables](../tests/migrate_test.go) | con el esquema al día solo se aseguran las tablas | TestMigrate_UpToDate_OnlyEnsuresTables |
| [TestMigrate_NewField_AddsOnlyThatColumn](../tests/migrate_test.go) | un campo nuevo agrega solo esa columna | TestMigrate_NewField_AddsOnlyThatColumn |
| [TestMigrate_IntrospectionFails_ReturnsErrorAndStops](../tests/migrate_test.go) | si falla la introspección, la migración devuelve el error y se detiene | TestMigrate_IntrospectionFails_ReturnsErrorAndStops |
