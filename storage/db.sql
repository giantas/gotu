-- query: InitDb
BEGIN TRANSACTION;
DROP TABLE IF EXISTS files;
CREATE TABLE files(
    id INTEGER NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    UNIQUE(name, path)
);
COMMIT;

-- query: FileInsertOrReplaceMany
INSERT OR REPLACE INTO files(name, path) VALUES ($1, $2) RETURNING id;

-- query: FileInsertOrReplace
INSERT OR REPLACE INTO files(name, path) VALUES($1, $2) RETURNING id;

-- query: FileDeleteOne
DELETE FROM files WHERE id = $1;

-- query: FileReadOne
SELECT id, name, path FROM files WHERE id = $1;

-- query: FileReadMany
SELECT id, name, path FROM files ORDER BY id DESC LIMIT $1 OFFSET $2;
