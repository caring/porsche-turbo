package db

import (
  "context"
  "database/sql"
  "database/sql/driver"
  "testing"

  "github.com/DATA-DOG/go-sqlmock"
  "github.com/google/uuid"
  "github.com/stretchr/testify/assert"

  "github.com/caring/porsche-turbo/pb"
)



// ensures that casting from proto to store structs occurs correctly
func TestNewTurbo(t *testing.T) {
  turboID := uuid.MustParse("72bc87f3-4a9f-4d05-93fe-844d3cd94c65")
  proto := pb.CreateTurboRequest{
    Name:       "Foobar",
  }

  r, err := NewTurbo(turboID.String(), &proto)

  assert.NoError(t, err, "Expected NewCategory not to error")
  assert.Equal(t, turboID, r.ID, "Expected UUIDs to match")
  assert.Equal(t, proto.Name, r.Name, "Expected name to be correctly assigned")
}

// ensures that casting from store to proto response occurs correctly
func TestTurbo_ToProto(t *testing.T) {
  turboID := uuid.MustParse("72bc87f3-4a9f-4d05-93fe-844d3cd94c65")

  turbo := &Turbo{
    ID:  turboID,
    Name:       "foobar",
  }

  r := turbo.ToProto()

  assert.Equal(t, turboID.String(), r.TurboId, "Expected field to be mapped back to proto object correctly")
  assert.Equal(t, "foobar", r.Name, "Expected field to be mapped back to proto object correctly")
}

func TestTurboService_get(t *testing.T) {
  turboID := uuid.MustParse("72bc87f3-4a9f-4d05-93fe-844d3cd94c65")
  stmt := map[string]string{
    "get-turbo": "SELECT turbos",
  }
  args := []driver.Value{
    "72bc87f3-4a9f-4d05-93fe-844d3cd94c65",
  }

  // ensures execution within a transaction occurs without error and the correct result is returned
  t.Run("With a provided transaction", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectBegin()
    mock.ExpectQuery("SELECT turbos").
      WithArgs(args...).
      WillReturnRows(
        sqlmock.NewRows([]string{"turbo_id", "name"}).
          AddRow(turboID, "Foobar"),
      )

    tx, err := store.GetTx()
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "transaction setup failed")
    }

    result, err := store.Turbo.GetTx(ToCtx(context.Background(), tx), turboID)
    assert.NoError(t, err, "Expecting no query error")

    assert.Equal(t, turboID, r.ID, "Expected correct turbo ID to be returned")
    assert.Equal(t, "Foobar", r.Name, "Expected correct name to be returned")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })

  // ensures that execution outside of transaction occurs without error and the correct result is returned
  t.Run("Without a provided transaction", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectQuery("SELECT turbos").
      WithArgs(args...).
      WillReturnRows(
        sqlmock.NewRows([]string{"turbo_id", "name"}).
          AddRow(turboID, "Foobar"),
      )

    result, err := store.Turbo.Get(context.Background(), turboID)
    assert.NoError(t, err, "Expecting no query error")

    assert.Equal(t, turboID, r.ID, "Expected correct turbo ID to be returned")
    assert.Equal(t, "Foobar", r.Name, "Expected correct name to be returned")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })

  // ensures a record not found is handled correctly
  t.Run("No rows returned", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectQuery("SELECT turbos").
      WithArgs(args...).WillReturnError(sql.ErrNoRows)

    _, err = store.Turbo.Get(context.Background(), turboID)
    assert.EqualError(t, err, "Error executing get turbo - 72bc87f3-4a9f-4d05-93fe-844d3cd94c65: the record you are attempting to find or update is not found", "Expecting no query error")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })
}

func TestTurboService_create(t *testing.T) {
  turboID := uuid.MustParse("72bc87f3-4a9f-4d05-93fe-844d3cd94c65")
  stmt := map[string]string{
    "create-turbo": "INSERT turbos",
  }
  input := &Turbo{
    ID:   turboID,
    Name: "Foobar",
  }
  args := []driver.Value{
    "72bc87f3-4a9f-4d05-93fe-844d3cd94c65",
    "Foobar",
  }

  // ensures that execution within a transaction occurs without error
  t.Run("With a provided transaction", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectBegin()
    mock.ExpectExec("INSERT turbos").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 1))

    tx, err := store.GetTx()
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "transaction setup failed")
    }

    err = store.Turbo.CreateTx(ToCtx(context.Background(), tx), input)
    assert.NoError(t, err, "Expecting no query error")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })

  // ensures that execution outside of a transaction occurs without error
  t.Run("Without a provided transaction", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectExec("INSERT turbos").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 1))

    err = store.Turbo.Create(context.Background(), input)
    assert.NoError(t, err, "Expecting no query error")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })

  // ensures that a failed record create is handled correctly
  t.Run("Failed record create", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectExec("INSERT turbos").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 0))

    err = store.Turbo.Create(context.Background(), input)
    assert.EqualError(t, err, "Error executing create turbo - &{72bc87f3-4a9f-4d05-93fe-844d3cd94c65 Foobar}: no new rows were created", "Expecting no query error")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })
}

func TestTurboService_update(t *testing.T) {
  turboID := uuid.MustParse("72bc87f3-4a9f-4d05-93fe-844d3cd94c65")
  stmt := map[string]string{
    "update-turbo": "UPDATE turbos",
  }
  input := &Turbo{
    ID:   turboID,
    Name: "Foobar",
  }
  args := []driver.Value{
    "Foobar",
    "72bc87f3-4a9f-4d05-93fe-844d3cd94c65",
  }

  // ensures that execution within a transaction occurs without error
  t.Run("With a provided transaction", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectBegin()
    mock.ExpectExec("UPDATE turbos").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 1))

    tx, err := store.GetTx()
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "transaction setup failed")
    }

    err = store.Turbo.UpdateTx(ToCtx(context.Background(), tx), input)
    assert.NoError(t, err, "Expecting no query error")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })

  // ensures execution out of a transaction occurs without error
  t.Run("Without a provided transaction", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectExec("UPDATE turbos").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 1))

    err = store.Turbo.Update(context.Background(), input)
    assert.NoError(t, err, "Expecting no query error")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })

  // ensures correct error to be returned when no rows are updated
  t.Run("No updates occurred", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectExec("UPDATE turbos").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 0))

    err = store.Turbo.Update(context.Background(), input)
    assert.EqualError(t, err, "Error executing update turbo - &{72bc87f3-4a9f-4d05-93fe-844d3cd94c65 Foobar}: no rows affected", "Expecting no query error")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })
}

func TestTurboService_delete(t *testing.T) {
  turboID := uuid.MustParse("72bc87f3-4a9f-4d05-93fe-844d3cd94c65")
  stmt := map[string]string{
    "delete-turbo": "UPDATE turbos",
  }
  args := []driver.Value{
    "72bc87f3-4a9f-4d05-93fe-844d3cd94c65",
  }

  // ensures that execution withing a transaction occurs without error
  t.Run("With a provided transaction", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectBegin()
    mock.ExpectExec("UPDATE turbos").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 1))

    tx, err := store.GetTx()
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "transaction setup failed")
    }

    err = store.Turbo.DeleteTx(ToCtx(context.Background(), tx), turboID)
    assert.NoError(t, err, "Expecting no query error")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })

  // ensures that execution outside of a transaction occurs without error
  t.Run("Without a provided transaction", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectExec("UPDATE turbos").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 1))

    err = store.Turbo.Delete(context.Background(), turboID)
    assert.NoError(t, err, "Expecting no query error")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })

  // ensures that deleting a non existent record is handled correctly
  t.Run("Deleting a non existent record", func(t *testing.T) {
    store, mock, err := NewTestDB(stmt)
    if ok := assert.NoError(t, err, "Expected no error"); !ok {
      assert.FailNow(t, "test setup failed")
    }

    mock.ExpectExec("UPDATE turbos").
      WithArgs(args...).
      WillReturnResult(sqlmock.NewResult(0, 0))

    err = store.Turbo.Delete(context.Background(), turboID)
    assert.EqualError(t, err, "Error executing delete turbo - 72bc87f3-4a9f-4d05-93fe-844d3cd94c65: the record you are attempting to find or update is not found", "Expecting not found error")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err, "Expecting all mock conditions to be met")
  })
}