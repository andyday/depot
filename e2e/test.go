package e2e

import (
	"context"
	"time"

	"github.com/andyday/depot"
	"github.com/andyday/go-log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/suite"
)

var (
	testWidgetKey = Widget{TenantID: "tenant", ID: "widget"}
	testWidget    = Widget{
		TenantID:    "tenant",
		ID:          "widget",
		Name:        "Widget",
		Description: "Test Widget",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	testWidgets = []Widget{
		{
			TenantID:            "tenant",
			ID:                  "widget1",
			Name:                "Widget 1",
			Description:         "Test Widget 1",
			Category:            "category1",
			ExpirationPartition: 1,
			Expiration:          aws.Time(time.Now().UTC().Add(-time.Hour)),
			CreatedAt:           time.Now().UTC().Add(-time.Hour),
			UpdatedAt:           time.Now().UTC(),
		},
		{
			TenantID:            "tenant",
			ID:                  "widget2",
			Name:                "Widget 2",
			Description:         "Test Widget 2",
			Category:            "category2",
			ExpirationPartition: 1,
			Expiration:          aws.Time(time.Now().UTC().Add(-time.Hour)),
			CreatedAt:           time.Now().UTC().Add(-time.Hour),
			UpdatedAt:           time.Now().UTC(),
		},
		{
			TenantID:            "tenant",
			ID:                  "widget3",
			Name:                "The Widget 3",
			Description:         "Test Widget 3",
			Category:            "category1",
			ExpirationPartition: 1,
			Expiration:          aws.Time(time.Now().UTC().Add(time.Hour)),
			CreatedAt:           time.Now().UTC().Add(time.Minute),
			UpdatedAt:           time.Now().UTC(),
		},
		{
			TenantID:            "tenant",
			ID:                  "widget4",
			Name:                "Widget 4",
			Description:         "Test Widget 4",
			Category:            "category3",
			ExpirationPartition: 1,
			Expiration:          aws.Time(time.Now().UTC().Add(time.Hour)),
			CreatedAt:           time.Now().UTC().Add(2 * time.Minute),
			UpdatedAt:           time.Now().UTC(),
		},
		{
			TenantID:    "tenant",
			ID:          "widget5",
			Name:        "Widget 5",
			Description: "Test Widget 5",
			CreatedAt:   time.Now().UTC().Add(3 * time.Minute),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			TenantID:    "tenant",
			ID:          "widget6",
			Name:        "Widget 6",
			Description: "Test Widget 6",
			CreatedAt:   time.Now().UTC().Add(4 * time.Minute),
			UpdatedAt:   time.Now().UTC(),
		},
	}
)

type Widget struct {
	TenantID            string     `depot:"tenantId,pk,index:created:pk,index:named:pk,index:category:pk"`
	ID                  string     `depot:"id,sk"`
	Name                string     `depot:"name,index:named:sk"`
	Category            string     `depot:"category,omitempty,index:category:sk"`
	Description         string     `depot:"description,omitempty"`
	Count               int64      `depot:"count,omitempty"`
	Total               int64      `depot:"total,omitempty"`
	ExpirationPartition int64      `depot:"expirationPartition,omitempty,index:expired:pk"`
	Expiration          *time.Time `depot:"expiration,omitempty,index:expired:sk"`
	CreatedAt           time.Time  `depot:"createdAt,index:created:sk"`
	UpdatedAt           time.Time  `depot:"updatedAt"`
}

type Suite struct {
	suite.Suite
	db     depot.Database
	widget depot.Table[Widget]
	ctx    context.Context
}

func (s *Suite) SetupTest() {
	s.ctx = context.Background()
	s.widget = depot.NewTable[Widget](s.db, "depot-widget")
	_, err := s.widget.Delete(s.ctx, testWidgetKey)
	s.NoError(err)
	_, err = s.widget.Get(s.ctx, testWidgetKey)
	s.ErrorIs(err, depot.ErrEntityNotFound)
}

func TestPut(s *Suite) {
	var widget Widget
	expected := testWidget
	_, err := s.widget.Put(s.ctx, expected)
	s.NoError(err)
	widget, err = s.widget.Get(s.ctx, testWidgetKey)
	s.NoError(err)
	log.Infof(s.ctx, "Widget: %+v", widget)
	s.Equal(expected, widget)
	expected.Name = "Widget (edited)"
	_, err = s.widget.Put(s.ctx, expected)
	s.NoError(err)
	widget, err = s.widget.Get(s.ctx, testWidgetKey)
	s.NoError(err)
	log.Infof(s.ctx, "Widget: %+v", widget)
	s.Equal(expected, widget)
}

func TestCreate(s *Suite) {
	var widget Widget
	expected := testWidget
	_, err := s.widget.Create(s.ctx, expected)
	s.NoError(err)
	widget, err = s.widget.Get(s.ctx, testWidgetKey)
	s.NoError(err)
	log.Infof(s.ctx, "Widget: %+v", widget)
	s.Equal(expected, widget)
	_, err = s.widget.Create(s.ctx, expected)
	s.ErrorIs(err, depot.ErrEntityAlreadyExists)
}

func TestUpdate(s *Suite) {
	var widget Widget
	expected := testWidget
	_, err := s.widget.Create(s.ctx, expected)
	s.NoError(err)
	_, err = s.widget.Update(s.ctx, Widget{
		TenantID:    expected.TenantID,
		ID:          expected.ID,
		Description: "New Description",
		Count:       5,
		Total:       5,
		UpdatedAt:   time.Now().UTC(),
	}, depot.Add("count"), depot.Subtract("total"))
	s.NoError(err)
	widget, err = s.widget.Get(s.ctx, testWidgetKey)
	s.NoError(err)
	s.Equal(expected.TenantID, widget.TenantID)
	s.Equal(expected.ID, widget.ID)
	s.Equal(expected.Name, widget.Name)
	s.Equal("New Description", widget.Description)
	s.Equal(int64(5), widget.Count)
	s.Equal(int64(-5), widget.Total)
	s.Equal(expected.CreatedAt, widget.CreatedAt)
}

func TestQueryCreatedIndex(s *Suite) {
	for _, w := range testWidgets {
		_, err := s.widget.Put(s.ctx, w)
		s.NoError(err)
	}

	widgets, page, err := s.widget.Query(s.ctx,
		"created",
		Widget{TenantID: "tenant", CreatedAt: time.Now().UTC().Add(-30 * time.Minute)},
		depot.GreaterThan("createdAt"),
		depot.Desc())
	s.NoError(err)
	s.Empty(page)
	s.Equal(4, len(widgets))
	s.Equal(testWidgets[5], widgets[0])
	s.Equal(testWidgets[4], widgets[1])
	s.Equal(testWidgets[3], widgets[2])
	s.Equal(testWidgets[2], widgets[3])
}

func TestQueryNamedIndex(s *Suite) {
	for _, w := range testWidgets {
		_, err := s.widget.Put(s.ctx, w)
		s.NoError(err)
	}
	widgets, page, err := s.widget.Query(s.ctx,
		"named",
		Widget{TenantID: "tenant", Name: "Widget"},
		depot.Prefix("name"),
		depot.Exists("expiration"),
		depot.Limit(2))

	s.NoError(err)
	s.NotEmpty(page)
	s.Equal(2, len(widgets))
	s.Equal(testWidgets[0], widgets[0])
	s.Equal(testWidgets[1], widgets[1])

	widgets, page, err = s.widget.Query(s.ctx,
		"named",
		Widget{TenantID: "tenant", Name: "Widget"},
		depot.Prefix("name"),
		depot.Exists("expiration"),
		depot.Limit(2),
		depot.Page(page))

	s.NoError(err)
	s.NotEmpty(page)
	s.Equal(1, len(widgets))
	s.Equal(testWidgets[3], widgets[0])

	widgets, page, err = s.widget.Query(s.ctx,
		"named",
		Widget{TenantID: "tenant", Name: "Widget"},
		depot.Prefix("name"),
		depot.Exists("expiration"),
		depot.Limit(2),
		depot.Page(page))

	s.NoError(err)
	s.Empty(page)
	s.Equal(0, len(widgets))
}
