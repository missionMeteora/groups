package groups

import (
	"github.com/itsmontoya/turtle"
	"github.com/missionMeteora/toolkit/errors"
)

const (
	// ErrInvalidType is returned when an invalid type is stored for an id
	ErrInvalidType = errors.Error("invalid type")
	// ErrGroupAlreadySet is returned when a group has been attempted to be set, but already exists for an id
	ErrGroupAlreadySet = errors.Error("cannot set group, group has already been set for this id")
	// ErrGroupNotSet is returned when a group has been attempted to be removed, but doesn't belong for an id
	ErrGroupNotSet = errors.Error("cannot remove group, group does not belong to this id")
)

// New will return a new instance of groups
func New(path string) (gp *Groups, err error) {
	var g Groups
	if g.db, err = turtle.New("groups", path, marshal, unmarshal); err != nil {
		return
	}

	gp = &g
	return
}

// Groups manages user groups
type Groups struct {
	db *turtle.Turtle
}

func (g *Groups) get(txn turtle.Txn, id string) (gm groupMap, err error) {
	var (
		val turtle.Value
		ok  bool
	)

	if val, err = txn.Get(id); err != nil {
		return
	}

	if gm, ok = val.(groupMap); !ok {
		err = ErrInvalidType
		return
	}

	return
}

// Get will get the group slice by id
func (g *Groups) Get(id string) (gs []string, err error) {
	var gm groupMap
	err = g.db.Read(func(txn turtle.Txn) (err error) {
		if gm, err = g.get(txn, id); err != nil {
			return
		}

		gs = gm.Slice()
		return
	})

	return
}

// Has will confirm if an id has a given group
func (g *Groups) Has(id, group string) (has bool) {
	var gm groupMap
	g.db.Read(func(txn turtle.Txn) (err error) {
		if gm, err = g.get(txn, id); err != nil {
			return
		}

		has = gm.Has(group)
		return
	})

	return
}

// Set will set a group to a given id
func (g *Groups) Set(id string, groups ...string) (gs []string, err error) {
	var (
		gm   groupMap
		errs errors.ErrorList
	)

	errs.Push(g.db.Update(func(txn turtle.Txn) (err error) {
		if gm, err = g.get(txn, id); err != nil {
			gm = make(groupMap)
			err = nil
		}

		for _, group := range groups {
			if !gm.Set(group) {
				errs.Push(ErrGroupAlreadySet)
				err = nil
				return
			}

			if err = txn.Put(id, gm.Dup()); err != nil {
				return
			}
		}

		return
	}))

	return gm.Slice(), errs.Err()
}

// Remove will remove a group from a given id
func (g *Groups) Remove(id, group string) (gs []string, err error) {
	var gm groupMap
	err = g.db.Update(func(txn turtle.Txn) (err error) {
		if gm, err = g.get(txn, id); err != nil {
			gm = make(groupMap)
			err = nil
		}

		if !gm.Remove(group) {
			return ErrGroupNotSet
		}

		return txn.Put(id, gm.Dup())
	})

	gs = gm.Slice()
	return
}

// Close will close groups
func (g *Groups) Close() (err error) {
	return g.db.Close()
}
