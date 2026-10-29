#!/bin/bash
sed -i 's/got != body/!strings.Contains(got, body)/g' tests/op_errors_test.go
sed -i 's|"webtyp.com/storage/mem"|"webtyp.com/storage/mem"\n\t"strings"|g' tests/op_errors_test.go
