package validate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequired(t *testing.T) {
	req := require.New(t)

	req.NotNil(Required("value", nil))

	req.Nil(Required("value", "test"))
	req.Nil(Required("value", 0))
	var emptyStruct struct{}
	req.Nil(Required("value", emptyStruct))
}

func TestNotEmpty(t *testing.T) {
	req := require.New(t)

	req.NotNil(NotEmpty("value", ""))

	req.Nil(NotEmpty("value", " "))
	req.Nil(NotEmpty("value", "test"))
	req.Nil(NotEmpty("value", " test "))
}

func TestNotBlank(t *testing.T) {
	req := require.New(t)

	req.NotNil(NotBlank("value", ""))
	req.NotNil(NotBlank("value", " "))

	req.Nil(NotBlank("value", "test"))
	req.Nil(NotBlank("value", " test "))
}

func TestEmail(t *testing.T) {
	req := require.New(t)

	req.NotNil(Email("value", "test"))
	req.NotNil(Email("value", " test "))
	req.NotNil(Email("value", "test@"))

	req.Nil(Email("value", "test@test"))
	req.Nil(Email("value", "test@test.org"))
}

func TestNotNegative(t *testing.T) {
	req := require.New(t)

	req.NotNil(NotNegative("value", -1))
	req.NotNil(NotNegative("value", -999))

	req.Nil(NotNegative("value", 0))
	req.Nil(NotNegative("value", 1))
	req.Nil(NotNegative("value", 999))
}

func TestNotZero(t *testing.T) {
	req := require.New(t)

	req.NotNil(NotZero("value", -1))
	req.NotNil(NotZero("value", 1))

	req.Nil(NotZero("value", 0))
}

func TestNotZeroOrPositive(t *testing.T) {
	req := require.New(t)

	req.NotNil(NotZeroOrPositive("value", 0))
	req.NotNil(NotZeroOrPositive("value", 1))

	req.Nil(NotZeroOrPositive("value", -1))
	req.Nil(NotZeroOrPositive("value", -999))
}

func TestNotPositive(t *testing.T) {
	req := require.New(t)

	req.NotNil(NotPositive("value", 1))
	req.NotNil(NotPositive("value", 999))

	req.Nil(NotPositive("value", 0))
	req.Nil(NotPositive("value", -1))
	req.Nil(NotPositive("value", -999))
}
