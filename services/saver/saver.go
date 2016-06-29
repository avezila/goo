package saver

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Saver struct {
	db *mgo.Collection
}

func New(server *grpc.Server) (*Saver, error) {
	session, err := mgo.Dial("saver-mongo")
	if err != nil {
		return nil, err
	}

	saver := &Saver{
		db: session.DB("saver").C("bytes"),
	}
	RegisterSaverServer(server, saver)
	return saver, nil
}

func (s *Saver) FilterRows(ctx context.Context, req *FilterRowsReq) (res *FilterRowsRes, err error) {
	res = &FilterRowsRes{}
	f := s.db.Find(bson.M{"id": bson.M{"$regex": "^" + req.Id + ".*"}})
	n, err := f.Count()

	res.Max = int32(n)
	err = f.Skip(int(req.From)).Limit(int(req.Count)).All(&res.Rows)

	return res, err
}

func (s *Saver) PushRow(ctx context.Context, req *PushRowReq) (*PushRowRes, error) {
	res := &PushRowRes{}

	_, err := s.db.Upsert(bson.M{"id": req.Row.Id}, bson.M{"$set": req.Row})

	return res, err
}
