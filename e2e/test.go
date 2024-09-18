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
		Refs:        []string{"ref1", "ref2"},
		Preferences: map[string]map[string]bool{
			"a": {
				"1": true,
				"2": false,
			},
			"b": {
				"1": false,
				"2": true,
			},
		},
		Data: map[string]interface{}{
			"c": "d",
			"e": "f",
		},
		Version:   123,
		Status:    "Active",
		TTL:       time.Now().Add(time.Hour).Unix(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
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

	testMessages = []Message{
		{
			TenantID: "tenant",
			ID:       1,
			Body:     "Message 1",
		},
		{
			TenantID: "tenant",
			ID:       2,
			Body:     "Message 2",
		},
		{
			TenantID: "tenant",
			ID:       3,
			Body:     "Message 3",
		},
		{
			TenantID: "tenant",
			ID:       4,
			Body:     "Message 4",
		},
		{
			TenantID: "tenant",
			ID:       5,
			Body:     "Message 5",
		},
		{
			TenantID: "tenant",
			ID:       6,
			Body:     "Message 6",
		},
		{
			TenantID: "tenant",
			ID:       7,
			Body:     "Message 7",
		},
	}
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
	Version             int64                      `depot:"version,omitempty"`
	Status              WidgetStatus               `depot:"status"`
	ExpirationPartition int64                      `depot:"expirationPartition,omitempty,index:expired:pk"`
	Expiration          *time.Time                 `depot:"expiration,omitempty,index:expired:sk"`
	CreatedAt           time.Time                  `depot:"createdAt,index:created:sk"`
	UpdatedAt           time.Time                  `depot:"updatedAt"`
}

type Message struct {
	TenantID string `depot:"tenantId,pk"`
	ID       int64  `depot:"id,sk"`
	Body     string `depot:"body"`
}

type Suite struct {
	suite.Suite
	db       depot.Database
	widgets  depot.Table[Widget]
	messages depot.Table[Message]
	ctx      context.Context
}

func (s *Suite) SetupTest() {
	s.ctx = context.Background()
	s.widgets = depot.NewTable[Widget](s.db, "depot-widget")
	_, err := s.widgets.Delete(s.ctx, testWidgetKey)
	s.NoError(err)
	_, err = s.widgets.Get(s.ctx, testWidgetKey)
	s.ErrorIs(err, depot.ErrEntityNotFound)
	s.messages = depot.NewTable[Message](s.db, "depot-message")
}

func TestPut(s *Suite) {
	var widget Widget
	expected := testWidget
	_, err := s.widgets.Put(s.ctx, expected)
	s.NoError(err)
	widget, err = s.widgets.Get(s.ctx, testWidgetKey)
	s.NoError(err)
	log.Infof(s.ctx, "Widget: %+v", widget)
	s.equals(expected, widget)
	expected.Name = "Widget (edited)"
	_, err = s.widgets.Put(s.ctx, expected)
	s.NoError(err)
	widget, err = s.widgets.Get(s.ctx, testWidgetKey)
	s.NoError(err)
	log.Infof(s.ctx, "Widget: %+v", widget)
	s.equals(expected, widget)
}

func TestCreate(s *Suite) {
	var widget Widget
	expected := testWidget
	_, err := s.widgets.Create(s.ctx, expected)
	s.NoError(err)
	widget, err = s.widgets.Get(s.ctx, testWidgetKey)
	s.NoError(err)
	log.Infof(s.ctx, "Widget: %+v", widget)
	s.equals(expected, widget)

	_, err = s.widgets.Create(s.ctx, expected)
	s.ErrorIs(err, depot.ErrEntityAlreadyExists)
}

func TestUpdate(s *Suite) {
	var widget Widget
	expected := testWidget
	_, err := s.widgets.Create(s.ctx, expected)
	s.NoError(err)
	_, err = s.widgets.Update(s.ctx, Widget{
		TenantID:    expected.TenantID,
		ID:          expected.ID,
		Description: "New Description",
		Count:       5,
		Total:       5,
		Version:     123,
		UpdatedAt:   time.Now().UTC(),
	}, depot.Add("count"), depot.Subtract("total"))
	s.NoError(err)
	widget, err = s.widgets.Get(s.ctx, testWidgetKey)
	s.NoError(err)
	s.Equal(expected.TenantID, widget.TenantID)
	s.Equal(expected.ID, widget.ID)
	s.Equal(expected.Name, widget.Name)
	s.Equal("New Description", widget.Description)
	s.Equal(int64(5), widget.Count)
	s.Equal(int64(-5), widget.Total)
	s.WithinDuration(expected.CreatedAt, widget.CreatedAt, time.Millisecond)
}

func TestQueryCreatedIndex(s *Suite) {
	for _, w := range testWidgets {
		_, err := s.widgets.Put(s.ctx, w)
		s.NoError(err)
	}

	widgets, page, err := s.widgets.Query(s.ctx,
		"created",
		Widget{TenantID: "tenant", CreatedAt: time.Now().UTC().Add(-30 * time.Minute)},
		depot.GreaterThan("createdAt"),
		depot.Desc())
	s.NoError(err)
	s.Empty(page)
	s.Equal(4, len(widgets))
	s.equals(testWidgets[5], widgets[0])
	s.equals(testWidgets[4], widgets[1])
	s.equals(testWidgets[3], widgets[2])
	s.equals(testWidgets[2], widgets[3])
}

func TestQueryNamedIndex(s *Suite) {
	for _, w := range testWidgets {
		_, err := s.widgets.Put(s.ctx, w)
		s.NoError(err)
	}
	widgets, page, err := s.widgets.Query(s.ctx,
		"named",
		Widget{TenantID: "tenant", Name: "Widget"},
		depot.GreaterThan("name"),
		depot.Exists("expiration"),
		depot.Limit(2))

	s.NoError(err)
	s.NotEmpty(page)
	s.Equal(2, len(widgets))
	s.equals(testWidgets[0], widgets[0])
	s.equals(testWidgets[1], widgets[1])

	widgets, page, err = s.widgets.Query(s.ctx,
		"named",
		Widget{TenantID: "tenant", Name: "Widget"},
		depot.GreaterThan("name"),
		depot.Exists("expiration"),
		depot.Limit(2),
		depot.Page(page))

	s.NoError(err)
	s.Equal(1, len(widgets))
	s.equals(testWidgets[3], widgets[0])

	if page != "" {
		widgets, page, err = s.widgets.Query(s.ctx,
			"named",
			Widget{TenantID: "tenant", Name: "Widget"},
			depot.GreaterThan("name"),
			depot.Exists("expiration"),
			depot.Limit(2),
			depot.Page(page))

		s.NoError(err)
		s.Empty(page)
		s.Equal(0, len(widgets))
	}
}

func TestQueryMessage(s *Suite) {
	for _, m := range testMessages {
		_, err := s.messages.Put(s.ctx, m)
		s.NoError(err)
	}

	messages, page, err := s.messages.Query(s.ctx, "",
		Message{TenantID: "tenant"},
		depot.Limit(3),
		depot.Desc())
	s.NoError(err)
	s.NotEmpty(page)
	s.Equal(3, len(messages))
	s.Equal(testMessages[6], messages[0])
	s.Equal(testMessages[5], messages[1])
	s.Equal(testMessages[4], messages[2])

	messages, page, err = s.messages.Query(s.ctx, "",
		Message{TenantID: "tenant"},
		depot.Limit(3),
		depot.Page(page),
		depot.Desc())
	s.NoError(err)
	s.NotEmpty(page)
	s.Equal(3, len(messages))
	s.Equal(testMessages[3], messages[0])
	s.Equal(testMessages[2], messages[1])
	s.Equal(testMessages[1], messages[2])

	messages, page, err = s.messages.Query(s.ctx, "",
		Message{TenantID: "tenant"},
		depot.Limit(3),
		depot.Page(page),
		depot.Desc())
	s.NoError(err)
	s.Empty(page)
	s.Equal(1, len(messages))
	s.Equal(testMessages[0], messages[0])
}

func (s *Suite) equals(a, b Widget) {
	s.WithinDuration(a.CreatedAt, b.CreatedAt, time.Millisecond)
	s.WithinDuration(a.UpdatedAt, b.UpdatedAt, time.Millisecond)
	if a.Expiration != nil {
		s.WithinDuration(*a.Expiration, *b.Expiration, time.Millisecond)
	}
	a.CreatedAt = b.CreatedAt
	a.UpdatedAt = b.UpdatedAt
	a.Expiration = b.Expiration

	s.Equal(a, b)
}
