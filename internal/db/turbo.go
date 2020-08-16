package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/caring/go-packages/pkg/errors"
	"github.com/google/uuid"

	"github.com/caring/porsche-turbo/pb"
)



// turboService provides an API for interacting with the turbos table
type turboService struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

// Turbo is a struct representation of a row in the turbos table
type Turbo struct {
	ID  	uuid.UUID
	Name  string
}

// protoTurbo is an interface that most proto turbo objects will satisfy
type protoTurbo interface {
	GetName() string
}

// NewTurbo is a convenience helper cast a proto turbo to it's DB layer struct
func NewTurbo(ID string, proto protoTurbo) (*Turbo, error) {
	mID, err := ParseUUID(ID)
	if err != nil {
		return nil, err
	}

	return &Turbo{
		ID:  	mID,
		Name: proto.GetName(),
	}, nil
}

// ToProto casts a db turbo into a proto response object
func (m *Turbo) ToProto() *pb.TurboResponse {
	return &pb.TurboResponse{
		Id:  				m.ID.String(),
		Name:       m.Name,
	}
}

// Get fetches a single turbo from the db
func (svc *turboService) Get(ctx context.Context, ID uuid.UUID) (*Turbo, error) {
	return svc.get(ctx, false, ID)
}

// GetTx fetches a single turbo from the db inside of a tx from ctx
func (svc *turboService) GetTx(ctx context.Context, ID uuid.UUID) (*Turbo, error) {
	return svc.get(ctx, true, ID)
}

// get fetches a single turbo from the db
func (svc *turboService) get(ctx context.Context, useTx bool, ID uuid.UUID) (*Turbo, error) {
	errMsg := func() string { return "Error executing get turbo - " + fmt.Sprint(ID) }

	var (
		stmt *sql.Stmt
		err  error
		tx   *sql.Tx
	)

	if useTx {

		if tx, err = FromCtx(ctx); err != nil {
			return nil, err
		}

		stmt = tx.Stmt(svc.stmts["get-turbo"])
	} else {
		stmt = svc.stmts["get-turbo"]
	}

	p := Turbo{}

	err = stmt.QueryRowContext(ctx, ID).
		Scan(&m.TurboID, &m.Name)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(ErrNotFound, errMsg())
		}

		return nil, errors.Wrap(err, errMsg())
	}

	return &p, nil
}

// Create a new turbo
func (svc *turboService) Create(ctx context.Context, input *Turbo) error {
	return svc.create(ctx, false, input)
}

// CreateTx creates a new turbo withing a tx from ctx
func (svc *turboService) CreateTx(ctx context.Context, input *Turbo) error {
	return svc.create(ctx, true, input)
}

// create a new turbo. if useTx = true then it will attempt to create the turbo within a transaction
// from context.
func (svc *turboService) create(ctx context.Context, useTx bool, input *Turbo) error {
	errMsg := func() string { return "Error executing create turbo - " + fmt.Sprint(input) }

	var (
		stmt *sql.Stmt
		err  error
		tx   *sql.Tx
	)

	if useTx {

		if tx, err = FromCtx(ctx); err != nil {
			return err
		}

		stmt = tx.Stmt(svc.stmts["create-turbo"])
	} else {
		stmt = svc.stmts["create-turbo"]
	}

	result, err := stmt.ExecContext(ctx, input.TurboID, input.Name)
	if err != nil {
		return errors.Wrap(err, errMsg())
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, errMsg())
	}

	if rowCount == 0 {
		return errors.Wrap(ErrNotCreated, errMsg())
	}

	return nil
}

// Update updates a single turbo row in the DB
func (svc *turboService) Update(ctx context.Context, input *Turbo) error {
	return svc.update(ctx, false, input)
}

// UpdateTx updates a single turbo row in the DB within a tx from ctx
func (svc *turboService) UpdateTx(ctx context.Context, input *Turbo) error {
	return svc.update(ctx, true, input)
}

// update a turbo. if useTx = true then it will attempt to update the turbo within a transaction
// from context.
func (svc *turboService) update(ctx context.Context, useTx bool, input *Turbo) error {
	errMsg := func() string { return "Error executing update turbo - " + fmt.Sprint(input) }

	var (
		stmt *sql.Stmt
		err  error
		tx   *sql.Tx
	)

	if useTx {

		if tx, err = FromCtx(ctx); err != nil {
			return err
		}

		stmt = tx.Stmt(svc.stmts["update-turbo"])
	} else {
		stmt = svc.stmts["update-turbo"]
	}

	result, err := stmt.ExecContext(ctx, input.Name, input.TurboID)
	if err != nil {
		return errors.Wrap(err, errMsg())
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, errMsg())
	}

	if rowCount == 0 {
		return errors.Wrap(ErrNoRowsAffected, errMsg())
	}

	return nil
}

// Delete sets deleted_at for a single turbos row
func (svc *turboService) Delete(ctx context.Context, ID uuid.UUID) error {
	return svc.delete(ctx, false, ID)
}

// DeleteTx sets deleted_at for a single turbos row within a tx from ctx
func (svc *turboService) DeleteTx(ctx context.Context, ID uuid.UUID) error {
	return svc.delete(ctx, true, ID)
}

// delete a turbo by setting deleted at. if useTx = true then it will attempt to delete the turbo within a transaction
// from context.
func (svc *turboService) delete(ctx context.Context, useTx bool, ID uuid.UUID) error {
	errMsg := func() string { return "Error executing delete turbo - " + ID.String() }

	var (
		stmt *sql.Stmt
		err  error
		tx   *sql.Tx
	)

	if useTx {

		if tx, err = FromCtx(ctx); err != nil {
			return err
		}

		stmt = tx.Stmt(svc.stmts["delete-turbo"])
	} else {
		stmt = svc.stmts["delete-turbo"]
	}

	result, err := stmt.ExecContext(ctx, ID)
	if err != nil {
		return errors.Wrap(err, errMsg())
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, errMsg())
	}

	if rowCount == 0 {
		return errors.Wrap(ErrNotFound, errMsg())
	}

	return nil
}

