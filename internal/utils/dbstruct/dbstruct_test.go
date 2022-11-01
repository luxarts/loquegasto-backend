package dbstruct

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Field1 string     `db:"field_1"`
	Field2 int        `db:"field_2"`
	Field3 *time.Time `db:"field_3"`
}

func TestGetColumns(t *testing.T) {
	// Given
	ts := testStruct{
		Field1: "test",
		Field2: 123,
		Field3: nil,
	}

	// When
	c := GetColumns(&ts)

	// Then
	require.Equal(t, []string{"field_1", "field_2", "field_3"}, c)
}
func TestGetValues(t *testing.T) {
	// Given
	ts := testStruct{
		Field1: "test",
		Field2: 123,
		Field3: nil,
	}

	// When
	v := GetValues(&ts)

	// Then
	require.Equal(t, "[test 123 <nil>]", fmt.Sprintf("%v", v))
}
