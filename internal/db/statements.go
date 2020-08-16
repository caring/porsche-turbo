package db

var statements = map[string]string{
  // inserts a new row into the turbos table
  "create-turbo": `
  INSERT INTO turbos (turbo_id, name)
    values(UUID_TO_BIN(?), ?)
  `,
  // soft deletes a turbo by id
  "delete-turbo": `
  UPDATE
    turbos
  SET
    deleted_at = NOW()
  WHERE
    turbo_id = UUID_TO_BIN(?)
    AND deleted_at IS NULL
  `,
  // gets a single turbo row by id
  "get-turbo": `
  SELECT
    turbo_id, name
  FROM
    turbos
  WHERE
    turbo_id = UUID_TO_BIN(?)
    AND deleted_at IS NULL
  `,
  // update a single turbo row by ID
  "update-turbo": `
  UPDATE
    turbos
  SET
    name = ?
  WHERE
    turbo_id = UUID_TO_BIN(?)
    AND deleted_at IS NULL
  `,
}
