package dynamo

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

func TestEncodePage(t *testing.T) {
	expected := map[string]types.AttributeValue{
		"a1": &types.AttributeValueMemberN{Value: "123"},
		"a2": &types.AttributeValueMemberS{Value: "s1"},
	}
	encoded, err := EncodePage(expected)
	assert.NoError(t, err)

	decoded, err := DecodePage(encoded)
	assert.NoError(t, err)
	assert.Equal(t, expected, decoded)
}
