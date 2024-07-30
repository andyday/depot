package e2e

import (
	"context"
	"os"
	"testing"

	"github.com/andyday/depot/firestore"
	"github.com/stretchr/testify/suite"
)

type FirestoreSuite struct {
	Suite
}

func TestFirestoreSuite(t *testing.T) {
	suite.Run(t, new(FirestoreSuite))
}

func (s *FirestoreSuite) SetupSuite() {
	var err error
	s.db, err = firestore.NewDatabase(context.Background(), os.Getenv("FIRESTORE_PROJECT_ID"), "")
	s.NoError(err)
}

func (s *FirestoreSuite) TestPut() {
	TestPut(&s.Suite)
}

func (s *FirestoreSuite) TestCreate() {
	TestCreate(&s.Suite)
}

func (s *FirestoreSuite) TestUpdate() {
	TestUpdate(&s.Suite)
}

func (s *FirestoreSuite) TestQueryCreatedIndex() {
	TestQueryCreatedIndex(&s.Suite)
}

func (s *FirestoreSuite) TestQueryNamedIndex() {
	TestQueryNamedIndex(&s.Suite)
}

func (s *FirestoreSuite) TestQueryMessage() {
	TestQueryMessage(&s.Suite)
}
