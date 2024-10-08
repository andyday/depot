package e2e

import (
	"context"
	"os"
	"testing"

	"github.com/andyday/depot/datastore"
	"github.com/stretchr/testify/suite"
)

type DatastoreSuite struct {
	Suite
}

func TestDatastoreSuite(t *testing.T) {
	t.Skip()
	suite.Run(t, new(DatastoreSuite))
}

func (s *DatastoreSuite) SetupSuite() {
	var err error
	s.db, err = datastore.NewDatabase(context.Background(), os.Getenv("DATASTORE_PROJECT_ID"), "")
	s.NoError(err)
}

func (s *DatastoreSuite) TestPut() {
	TestPut(&s.Suite)
}

func (s *DatastoreSuite) TestCreate() {
	TestCreate(&s.Suite)
}

func (s *DatastoreSuite) TestUpdate() {
	TestUpdate(&s.Suite)
}

func (s *DatastoreSuite) TestQueryCreatedIndex() {
	TestQueryCreatedIndex(&s.Suite)
}

func (s *DatastoreSuite) TestQueryNamedIndex() {
	TestQueryNamedIndex(&s.Suite)
}

func (s *DatastoreSuite) TestQueryMessage() {
	TestQueryMessage(&s.Suite)
}
