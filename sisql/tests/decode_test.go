package sisql_test

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
	"github.com/wonksing/si/v2/tests/testmodels"
)

var (
	bmap = map[string]interface{}{
		"nil":   "",
		"int2_": 123,
	}
)

func BenchmarkDecode_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		table := testmodels.Table{}
		byt, _ := json.Marshal(bmap)
		err := json.Unmarshal(byt, &table)
		require.Nil(b, err)
	}
}
func BenchmarkDecode_Mapstructure(b *testing.B) {
	for i := 0; i < b.N; i++ {
		table := testmodels.Table{}
		err := mapstructure.Decode(bmap, &table)
		require.Nil(b, err)
	}
}
