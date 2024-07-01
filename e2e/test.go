package e2e

import (
	"context"
	"time"

	"github.com/andyday/depot"
	"github.com/andyday/depot/transform"
	"github.com/andyday/depot/types"
	"github.com/andyday/go-log"
	"github.com/stretchr/testify/suite"
)

var (
	testWidgetKey = Widget{TenantID: "tenant", ID: "widget"}
	testWidget    = Widget{
		TenantID:  "tenant",
		ID:        "widget",
		Name:      "Widget",
		CreatedAt: time.Now().UTC(),
	}
)

type Widget struct {
	TenantID    string      `depot:"tenantId,pk"`
	ID          string      `depot:"id,sk"`
	Name        string      `depot:"name"`
	Description string      `depot:"description,omitempty"`
	Count       interface{} `depot:"count,omitempty"`
	Total       interface{} `depot:"total,omitempty"`
	CreatedAt   time.Time   `depot:"createdAt"`
}

func (w Widget) CountInt() int { return types.NumberAs[int](w.Count) }
func (w Widget) TotalInt() int { return types.NumberAs[int](w.Total) }

type Suite struct {
	suite.Suite
	db     depot.Database
	widget *depot.Table
	ctx    context.Context
}

func (s *Suite) SetupTest() {
	s.ctx = context.Background()
	s.widget = s.db.Table("depot-widget")
	err := s.widget.Delete(s.ctx, testWidgetKey)
	s.NoError(err)
	err = s.widget.Get(s.ctx, testWidgetKey)
	s.ErrorIs(err, types.ErrEntityNotFound)
}

func TestPut(s *Suite) {
	expected := testWidget
	err := s.widget.Put(s.ctx, expected)
	s.NoError(err)
	widget := testWidgetKey
	err = s.widget.Get(s.ctx, &widget)
	s.NoError(err)
	log.Infof(s.ctx, "Widget: %+v", widget)
	s.Equal(expected, widget)
	expected.Name = "Widget (edited)"
	err = s.widget.Put(s.ctx, expected)
	s.NoError(err)
	err = s.widget.Get(s.ctx, &widget)
	s.NoError(err)
	log.Infof(s.ctx, "Widget: %+v", widget)
	s.Equal(expected, widget)
}

func TestCreate(s *Suite) {
	expected := testWidget
	err := s.widget.Create(s.ctx, expected)
	s.NoError(err)
	widget := testWidgetKey
	err = s.widget.Get(s.ctx, &widget)
	s.NoError(err)
	log.Infof(s.ctx, "Widget: %+v", widget)
	s.Equal(expected, widget)
	err = s.widget.Create(s.ctx, expected)
	s.ErrorIs(err, types.ErrEntityAlreadyExists)
}

func TestUpdate(s *Suite) {
	expected := testWidget
	err := s.widget.Create(s.ctx, expected)
	s.NoError(err)
	err = s.widget.Update(s.ctx, &Widget{
		TenantID:    expected.TenantID,
		ID:          expected.ID,
		Description: "New Description",
		Count:       transform.Add(5),
		Total:       transform.Subtract(5),
	})
	s.NoError(err)
	widget := testWidgetKey
	err = s.widget.Get(s.ctx, &widget)
	s.NoError(err)
	s.Equal(expected.TenantID, widget.TenantID)
	s.Equal(expected.ID, widget.ID)
	s.Equal(expected.Name, widget.Name)
	s.Equal("New Description", widget.Description)
	s.Equal(5, widget.CountInt())
	s.Equal(-5, widget.TotalInt())
	s.Equal(expected.CreatedAt, widget.CreatedAt)
}
