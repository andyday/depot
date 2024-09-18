package depot_test

import (
	"context"
	"errors"
	"testing"

	"github.com/andyday/depot"
	"github.com/andyday/depot/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var errTest = errors.New("test error")

type Record struct {
	Name string
}

type DepotSuite struct {
	suite.Suite
	ctx context.Context
	db  *mocks.Database
	tbl depot.Table[Record]
}

func TestDepotSuite(t *testing.T) {
	suite.Run(t, new(DepotSuite))
}

func (s *DepotSuite) SetupTest() {
	s.ctx = context.Background()
	s.db = mocks.NewDatabase(s.T())
	s.tbl = depot.NewTable[Record](s.db, "record")

}

func (s *DepotSuite) TestGet() {
	var (
		in  Record
		out Record
		err error
	)
	s.db.On("Get", s.ctx, "record", &in).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*Record)
		*arg = Record{Name: "record"}
	}).Once()

	out, err = s.tbl.Get(s.ctx, in)
	s.NoError(err)
	s.Equal(Record{Name: "record"}, out)

	s.db.On("Get", s.ctx, "record", &in).Return(errTest).Once()
	_, err = s.tbl.Get(s.ctx, in)
	s.ErrorIs(err, errTest)
}

func (s *DepotSuite) TestPut() {
	var (
		in  = Record{Name: "record"}
		out Record
		err error
	)
	s.db.On("Put", s.ctx, "record", &in).Return(nil).Once()

	out, err = s.tbl.Put(s.ctx, in)
	s.NoError(err)
	s.Equal(in, out)

	s.db.On("Put", s.ctx, "record", &in).Return(errTest).Once()
	_, err = s.tbl.Put(s.ctx, in)
	s.ErrorIs(err, errTest)
}

func (s *DepotSuite) TestCreate() {
	var (
		in  = Record{Name: "record"}
		out Record
		err error
	)
	s.db.On("Create", s.ctx, "record", &in).Return(nil).Once()

	out, err = s.tbl.Create(s.ctx, in)
	s.NoError(err)
	s.Equal(in, out)

	s.db.On("Create", s.ctx, "record", &in).Return(errTest).Once()
	_, err = s.tbl.Create(s.ctx, in)
	s.ErrorIs(err, errTest)
}

func (s *DepotSuite) TestUpdate() {
	var (
		in  = Record{Name: "record"}
		op  = depot.Add("field")
		out Record
		err error
	)
	s.db.On("Update", s.ctx, "record", &in, op).Return(nil).Once()

	out, err = s.tbl.Update(s.ctx, in, op)
	s.NoError(err)
	s.Equal(in, out)

	s.db.On("Update", s.ctx, "record", &in).Return(errTest).Once()
	_, err = s.tbl.Update(s.ctx, in)
	s.ErrorIs(err, errTest)
}

func (s *DepotSuite) TestDelete() {
	var (
		in  = Record{Name: "record"}
		out Record
		err error
	)
	s.db.On("Delete", s.ctx, "record", &in).Return(nil).Once()

	out, err = s.tbl.Delete(s.ctx, in)
	s.NoError(err)
	s.Equal(in, out)

	s.db.On("Delete", s.ctx, "record", &in).Return(errTest).Once()
	_, err = s.tbl.Delete(s.ctx, in)
	s.ErrorIs(err, errTest)
}

func (s *DepotSuite) TestQuery() {
	var (
		in       = Record{Name: "record"}
		inList   []Record
		entities = []Record{{Name: "record"}}
		out      []Record
		next     string
		op       = depot.Equal("field")
		err      error
	)
	s.db.On("Query", s.ctx, "record", "kind", &in, &inList, op).
		Return("next", nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(4).(*[]Record)
			*arg = entities
		}).
		Once()

	out, next, err = s.tbl.Query(s.ctx, "kind", in, op)
	s.NoError(err)
	s.Equal("next", next)
	s.Equal(entities, out)
}
