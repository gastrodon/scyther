package storage

const (
	TABLE_QUEUES   = "queues"
	TABLE_MESSAGES = "messages"

	FIELDS_QUEUES   = "id, name, capacity, size"
	FIELDS_MESSAGES = "queue, data"

	SCHEMA_QUEUES = `
	        id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
	        name CHAR(255) UNIQUE,
	        capacity BIGINT UNSIGNED,
	        size BIGINT UNSIGNED NOT NULL,
	        ordered BIGINT UNSIGNED UNIQUE NOT NULL AUTO_INCREMENT`
	SCHEMA_MESSAGES = `
			queue CHAR(36) NOT NULL,
			data BINARY(255) NOT NULL,
			ordered BIGINT UNSIGNED UNIQUE NOT NULL AUTO_INCREMENT`

	CREATE_TABLE_QUEUES   = "CREATE TABLE IF NOT EXISTS " + TABLE_QUEUES + " (" + SCHEMA_QUEUES + ")"
	CREATE_TABLE_MESSAGES = "CREATE TABLE IF NOT EXISTS " + TABLE_MESSAGES + " (" + SCHEMA_MESSAGES + ")"

	CLEAR_TABLE_QUEUES   = "DELETE FROM " + TABLE_QUEUES
	CLEAR_TABLE_MESSAGES = "DELETE FROM " + TABLE_MESSAGES

	COUNT_QUEUES  = "SELECT count(*) FROM " + TABLE_QUEUES
	READ_QUEUES   = "SELECT " + FIELDS_QUEUES + " FROM " + TABLE_QUEUES + " ORDER BY ordered DESC"
	READ_QUEUE    = "SELECT " + FIELDS_QUEUES + " FROM " + TABLE_QUEUES + " WHERE id=? LIMIT 1"
	READ_QUEUE_ID = "SELECT id FROM " + TABLE_QUEUES + " WHERE name=? LIMIT 1"
	WRITE_QUEUE   = "INSERT INTO " + TABLE_QUEUES + "(" + FIELDS_QUEUES + ") VALUES (?, ?, ?, 0)"
	DELETE_QUEUE  = "DELETE FROM " + TABLE_QUEUES + " WHERE id=? LIMIT 1"

	DELETE_MESSAGES = "DELETE FROM " + TABLE_MESSAGES + " WHERE queue=?"

	TEST_READ_QUEUE_ID_EXISTS = "SELECT count(id)=1 FROM " + TABLE_QUEUES + " WHERE id=?"
)
