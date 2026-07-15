# PostgreSQL Notes (PSQL.md)

> Personal PostgreSQL cheat sheet for backend development.

---

# Table of Contents

- Linux Commands
- Connect to PostgreSQL
- Databases
- Tables
- CRUD
- ALTER TABLE
- Constraints
- Indexes
- Users & Roles
- Help Commands
- Common Data Types
- SQL Execution Order
- Best Practices

---

# Linux Commands

## Start PostgreSQL

```bash
sudo systemctl start postgresql
```

## Stop PostgreSQL

```bash
sudo systemctl stop postgresql
```

## Restart PostgreSQL

```bash
sudo systemctl restart postgresql
```

## Check Status

```bash
sudo systemctl status postgresql
```

## Enable on Boot

```bash
sudo systemctl enable postgresql
```

## List PostgreSQL Clusters

```bash
pg_lsclusters
```

Example:

```text
Ver Cluster Port Status Owner
18  main    5432 online postgres
```

---

# Connect to PostgreSQL

## Peer Authentication (Linux User)

```bash
sudo -u postgres psql
```

## Password Authentication

```bash
psql -h localhost -U postgres
```

## Connect to Specific Database

```bash
psql -h localhost -U postgres -d url_shortener
```

## PostgreSQL Version

```bash
psql --version
```

## Find psql

```bash
which psql
```

---

# Databases

## List Databases

```sql
\l
```

## Current Database

```sql
SELECT current_database();
```

## Create Database

```sql
CREATE DATABASE url_shortener;
```

## Connect Database

```sql
\c url_shortener
```

## Delete Database

```sql
DROP DATABASE url_shortener;
```

---

# Tables

## List Tables

```sql
\dt
```

## Describe Table

```sql
\d urls
```

## Detailed Table Information

```sql
\d+ urls
```

## Create Table

```sql
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    click_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## Delete Table

```sql
DROP TABLE urls;
```

---

# CRUD

## INSERT

```sql
INSERT INTO urls
(original_url, short_code)
VALUES
('https://google.com', 'abc123');
```

---

## SELECT All

```sql
SELECT *
FROM urls;
```

---

## SELECT One

```sql
SELECT *
FROM urls
WHERE id = 1;
```

---

## SELECT Specific Columns

```sql
SELECT
    short_code,
    original_url
FROM urls;
```

---

## UPDATE

```sql
UPDATE urls
SET click_count = click_count + 1
WHERE id = 1;
```

---

## DELETE

```sql
DELETE
FROM urls
WHERE id = 1;
```

---

## Count Rows

```sql
SELECT COUNT(*)
FROM urls;
```

---

# ALTER TABLE

## Add Column

```sql
ALTER TABLE urls
ADD COLUMN expires_at TIMESTAMP;
```

---

## Remove Column

```sql
ALTER TABLE urls
DROP COLUMN expires_at;
```

---

## Rename Column

```sql
ALTER TABLE urls
RENAME COLUMN short_code TO code;
```

---

## Change Column Type

```sql
ALTER TABLE urls
ALTER COLUMN code TYPE VARCHAR(20);
```

---

## Rename Table

```sql
ALTER TABLE urls
RENAME TO links;
```

---

# Constraints

## PRIMARY KEY

```sql
PRIMARY KEY
```

Each row has a unique identifier.

---

## UNIQUE

```sql
UNIQUE
```

No duplicate values allowed.

Example:

```sql
short_code VARCHAR(10) UNIQUE
```

---

## NOT NULL

```sql
NOT NULL
```

Column cannot be empty.

---

## DEFAULT

```sql
DEFAULT 0
```

Provides a default value.

---

## FOREIGN KEY

```sql
REFERENCES
```

Example:

```sql
user_id INTEGER REFERENCES users(id)
```

---

# Indexes

## Create Index

```sql
CREATE INDEX idx_short_code
ON urls(short_code);
```

Indexes make searching much faster.

---

# Users & Roles

## List Users

```sql
\du
```

---

## Create User

```sql
CREATE USER shamil
WITH PASSWORD 'password';
```

---

## Grant Privileges

```sql
GRANT ALL PRIVILEGES
ON DATABASE url_shortener
TO shamil;
```

---

# Help Commands

## PostgreSQL Commands

```sql
\?
```

---

## SQL Commands

```sql
\h
```

---

## CREATE TABLE Help

```sql
\h CREATE TABLE
```

---

## ALTER TABLE Help

```sql
\h ALTER TABLE
```

---

## INSERT Help

```sql
\h INSERT
```

---

## SELECT Help

```sql
\h SELECT
```

---

# Exit

```sql
\q
```

---

# Common PostgreSQL Data Types

| Type | Description | Example |
|------|-------------|---------|
| SERIAL | Auto Increment Integer | id |
| INTEGER | Whole Number | click_count |
| BIGINT | Large Integer | views |
| BOOLEAN | true / false | completed |
| TEXT | Long Text | original_url |
| VARCHAR(255) | Limited Text | username |
| TIMESTAMP | Date & Time | created_at |
| DATE | Date Only | birthday |
| UUID | Unique Identifier | user_id |

---

# SQL Execution Order

Typical workflow:

```text
Start PostgreSQL
        ↓
Connect (psql)
        ↓
Choose Database
        ↓
Create Tables
        ↓
INSERT
        ↓
SELECT
        ↓
UPDATE
        ↓
DELETE
        ↓
Exit
```

---

# PostgreSQL Hierarchy

```text
Cluster
│
├── Database
│   ├── Table
│   │   ├── Row (Record)
│   │   │   ├── Column (Field)
```

Example:

```text
Cluster
│
├── url_shortener
│   ├── urls
│   │   ├── id
│   │   ├── original_url
│   │   ├── short_code
│   │   ├── click_count
│   │   └── created_at
```

---

# Common SQL Keywords

```sql
CREATE
ALTER
DROP
INSERT
SELECT
UPDATE
DELETE
WHERE
ORDER BY
LIMIT
OFFSET
GROUP BY
HAVING
JOIN
LEFT JOIN
RIGHT JOIN
INNER JOIN
COUNT
DISTINCT
AS
```

---

# Best Practices

- Always use `PRIMARY KEY`.
- Use `NOT NULL` where appropriate.
- Add `UNIQUE` for values that must not repeat.
- Avoid `SELECT *` in production code.
- Use indexes for frequently searched columns.
- Use foreign keys for relationships.
- Use transactions for multiple dependent operations.
- Keep SQL readable and properly formatted.

---

# Learning Roadmap

## Beginner

- [x] CREATE TABLE
- [x] INSERT
- [x] SELECT
- [x] UPDATE
- [x] DELETE
- [x] WHERE
- [x] PRIMARY KEY
- [x] UNIQUE
- [x] FOREIGN KEY

## Intermediate

- [ ] JOIN
- [ ] GROUP BY
- [ ] HAVING
- [ ] INDEXES
- [ ] VIEWS
- [ ] TRANSACTIONS
- [ ] EXPLAIN

## Advanced

- [ ] Stored Procedures
- [ ] Triggers
- [ ] Functions
- [ ] CTE (WITH)
- [ ] Partitioning
- [ ] Replication
- [ ] Performance Tuning


<!-- 

docker run -it --rm \
  --network url-shortener_backend \
  --env-file .env \
  -e DB_HOST=postgres \
  -v $(pwd):/app \
  -w /app \
  golang:1.26-alpine \
  go run scripts/create_admin.go


 -->