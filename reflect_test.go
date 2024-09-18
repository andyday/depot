package depot

import (
	"testing"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/stretchr/testify/assert"
)

type WidgetStatus string

type Widget struct {
	TenantID            string                     `depot:"tenantId,pk,index:created:pk,index:named:pk,index:category:pk"`
	ID                  string                     `depot:"id,sk"`
	Name                string                     `depot:"name,index:named:sk"`
	Category            string                     `depot:"category,omitempty,index:category:sk"`
	Description         string                     `depot:"desc,omitempty"`
	Count               int64                      `depot:"count,omitempty"`
	Total               int64                      `depot:"total,omitempty"`
	Refs                []string                   `depot:"refs,omitempty"`
	Preferences         map[string]map[string]bool `depot:"preferences,omitempty"`
	Data                map[string]interface{}     `depot:"data,omitempty"`
	TTL                 int64                      `depot:"ttl,ttl"`
	Exclude             string                     `depot:"-"`
	Version             int64                      `depot:"version,omitempty"`
	Status              WidgetStatus               `depot:"status"`
	ExpirationPartition int64                      `depot:"expirationPartition,omitempty,index:expired:pk"`
	Expiration          *time.Time                 `depot:"expiration,omitempty,index:expired:sk"`
	CreatedAt           time.Time                  `depot:"createdAt,index:created:sk"`
	UpdatedAt           time.Time                  `depot:"updatedAt"`
}

func TestEntityKey(t *testing.T) {
	k, err := EntityKey(struct {
		A string `depot:"a,pk"`
	}{A: "val"})
	assert.NoError(t, err)
	assert.Equal(t, Key{Partition: KeyPart{Name: "a", Value: "val"}}, k)
	assert.Equal(t, "val", k.String())

	k, err = EntityKey(struct {
		A string `depot:"a,pk"`
		B int    `depot:"b,sk"`
	}{A: "av", B: 123})
	assert.NoError(t, err)
	assert.Equal(t, Key{
		Partition: KeyPart{Name: "a", Value: "av"},
		Sort:      KeyPart{Name: "b", Value: 123},
	}, k)
	assert.Equal(t, "av:123", k.String())

	_, err = EntityKey("key")
	assert.ErrorIs(t, err, ErrInvalidEntityType)
}

func TestEntityProperties(t *testing.T) {
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Hour)
	props, err := EntityProperties(Widget{
		TenantID:  "tv",
		ID:        "iv",
		Name:      "nv",
		TTL:       123,
		Status:    "sv",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})
	assert.NoError(t, err)
	assert.Equal(t, []Property{
		{Name: "tenantId", Value: "tv", Index: true},
		{Name: "id", Value: "iv", Index: true},
		{Name: "name", Value: "nv", Index: true},
		{Name: "ttl", Value: int64(123)},
		{Name: "status", Value: WidgetStatus("sv")},
		{Name: "createdAt", Value: createdAt, Index: true},
		{Name: "updatedAt", Value: updatedAt},
	}, props)

	_, err = EntityProperties("entity")
	assert.ErrorIs(t, err, ErrInvalidEntityType)
}

func TestEntityMap(t *testing.T) {
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Hour)
	props, err := EntityMap(Widget{
		TenantID:  "tv",
		ID:        "iv",
		Name:      "nv",
		TTL:       123,
		Status:    "sv",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, false)
	assert.NoError(t, err)
	assert.Equal(t, map[string]any{
		"tenantId":  "tv",
		"id":        "iv",
		"name":      "nv",
		"ttl":       int64(123),
		"status":    WidgetStatus("sv"),
		"createdAt": createdAt,
		"updatedAt": updatedAt,
	}, props)

	props, err = EntityMap(Widget{
		TenantID:  "tv",
		ID:        "iv",
		Name:      "nv",
		TTL:       123,
		Status:    "sv",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, true)
	assert.NoError(t, err)
	assert.Equal(t, map[string]any{
		"tenantId":  "tv",
		"id":        "iv",
		"name":      "nv",
		"ttl":       time.Unix(123, 0),
		"status":    WidgetStatus("sv"),
		"createdAt": createdAt,
		"updatedAt": updatedAt,
	}, props)

	_, err = EntityMap("entity", false)
	assert.ErrorIs(t, err, ErrInvalidEntityType)
}

func TestEntityFromMap(t *testing.T) {
	var entity1, entity2 Widget
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Hour)
	exp := createdAt.Add(2 * time.Hour)

	err := EntityFromMap(map[string]any{
		"tenantId": "tv",
		"id":       "iv",
		"name":     "nv",
		"ttl":      int64(123),
		"status":   WidgetStatus("sv"),
		"refs":     []any{"a", "b"},
		"preferences": map[string]any{
			"a": map[string]any{"aa": true},
			"b": map[string]any{"bb": true},
		},
		"createdAt": createdAt,
		"updatedAt": updatedAt,
	}, &entity1, false)
	assert.NoError(t, err)
	assert.Equal(t, Widget{
		TenantID: "tv",
		ID:       "iv",
		Name:     "nv",
		TTL:      123,
		Status:   "sv",
		Refs:     []string{"a", "b"},
		Preferences: map[string]map[string]bool{
			"a": {"aa": true},
			"b": {"bb": true},
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, entity1)

	err = EntityFromMap(map[string]any{
		"tenantId":   "tv",
		"id":         "iv",
		"name":       "nv",
		"ttl":        time.Unix(123, 0),
		"status":     WidgetStatus("sv"),
		"expiration": exp,
		"createdAt":  createdAt,
		"updatedAt":  updatedAt,
	}, &entity2, true)
	assert.NoError(t, err)
	assert.Equal(t, Widget{
		TenantID:   "tv",
		ID:         "iv",
		Name:       "nv",
		TTL:        123,
		Status:     "sv",
		Expiration: &exp,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, entity2)

	err = EntityFromMap(map[string]any{}, ptr.String("string"), false)
	assert.ErrorIs(t, err, ErrInvalidEntityType)
}

func TestEntityFromPropertyMap(t *testing.T) {
	var entity Widget
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Hour)

	err := EntityFromPropertyMap(map[string]Property{
		"tenantId": {Name: "tenantId", Value: "tv", Index: true},
		"id":       {Name: "id", Value: "iv", Index: true},
		"name":     {Name: "name", Value: "nv", Index: true},
		"ttl":      {Name: "ttl", Value: int64(123)},
		"status":   {Name: "status", Value: WidgetStatus("sv")},
		"refs":     {Name: "refs", Value: []any{"a", "b"}},
		"preferences": {Name: "preferences", Value: map[string]any{
			"a": map[string]any{"aa": true},
			"b": map[string]any{"bb": true},
		}},
		"createdAt": {Name: "createdAt", Value: createdAt},
		"updatedAt": {Name: "updatedAt", Value: updatedAt},
	}, &entity)
	assert.NoError(t, err)
	assert.Equal(t, Widget{
		TenantID: "tv",
		ID:       "iv",
		Name:     "nv",
		TTL:      123,
		Status:   "sv",
		Refs:     []string{"a", "b"},
		Preferences: map[string]map[string]bool{
			"a": {"aa": true},
			"b": {"bb": true},
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, entity)

	err = EntityFromPropertyMap(map[string]Property{}, ptr.String("string"))
	assert.ErrorIs(t, err, ErrInvalidEntityType)
}

func TestEntityFromProperties(t *testing.T) {
	var entity Widget
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Hour)

	err := EntityFromProperties([]Property{
		{Name: "tenantId", Value: "tv", Index: true},
		{Name: "id", Value: "iv", Index: true},
		{Name: "name", Value: "nv", Index: true},
		{Name: "ttl", Value: int64(123)},
		{Name: "status", Value: WidgetStatus("sv")},
		{Name: "refs", Value: []any{"a", "b"}},
		{Name: "preferences", Value: map[string]any{
			"a": map[string]any{"aa": true},
			"b": map[string]any{"bb": true},
		}},
		{Name: "createdAt", Value: createdAt},
		{Name: "updatedAt", Value: updatedAt},
	}, &entity)
	assert.NoError(t, err)
	assert.Equal(t, Widget{
		TenantID: "tv",
		ID:       "iv",
		Name:     "nv",
		TTL:      123,
		Status:   "sv",
		Refs:     []string{"a", "b"},
		Preferences: map[string]map[string]bool{
			"a": {"aa": true},
			"b": {"bb": true},
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, entity)

	err = EntityFromPropertyMap(map[string]Property{}, ptr.String("string"))
	assert.ErrorIs(t, err, ErrInvalidEntityType)
}

func TestEntityUpdates(t *testing.T) {
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Hour)
	add := Add("total")
	updates, err := EntityUpdates(Widget{
		TenantID:  "tv",
		ID:        "iv",
		Name:      "nv",
		Total:     5,
		TTL:       123,
		Status:    "sv",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, []UpdateOp{add})
	assert.NoError(t, err)
	assert.Equal(t, []Update{
		{Name: "name", Value: "nv"},
		{Name: "total", Value: int64(5), Op: add},
		{Name: "ttl", Value: int64(123)},
		{Name: "status", Value: WidgetStatus("sv")},
		{Name: "createdAt", Value: createdAt},
		{Name: "updatedAt", Value: updatedAt},
	}, updates)

	_, err = EntityUpdates("entity", nil)
	assert.ErrorIs(t, err, ErrInvalidEntityType)
}

func TestEntityConditions(t *testing.T) {
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Hour)
	gt := GreaterThan("createdAt")
	desc := Desc()
	sortField, updates, err := EntityConditions("created", Widget{
		TenantID:  "tv",
		ID:        "iv",
		Name:      "nv",
		TTL:       123,
		Status:    "sv",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, []QueryOp{gt, desc})
	assert.NoError(t, err)
	assert.Equal(t, "createdAt", sortField)
	assert.Equal(t, []EntityCondition{
		{Name: "tenantId", Value: "tv", KeyType: KeyTypePartition},
		{Name: "id", Value: "iv", KeyType: KeyTypeSort},
		{Name: "name", Value: "nv"},
		{Name: "ttl", Value: int64(123)},
		{Name: "status", Value: WidgetStatus("sv")},
		{Name: "createdAt", Value: createdAt, KeyType: KeyTypeSort, Op: gt},
		{Name: "updatedAt", Value: updatedAt},
	}, updates)

	_, err = EntityUpdates("entity", nil)
	assert.ErrorIs(t, err, ErrInvalidEntityType)
}

func TestRealSlice(t *testing.T) {
	assert.Equal(t, []string{"a", "b"}, RealSlice([]interface{}{"a", "b"}))
	assert.Equal(t, []int{1, 2}, RealSlice([]interface{}{1, 2}))
	assert.Equal(t, []int64{1, 2}, RealSlice([]interface{}{int64(1), int64(2)}))
}
