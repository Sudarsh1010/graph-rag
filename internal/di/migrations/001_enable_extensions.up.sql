-- Prerequisites:
-- 1. pgvector extension must be installed on PostgreSQL server
--    Ubuntu/Debian: apt-get install postgresql-16-pgvector
-- 2. Apache AGE extension must be installed and configured
--    Requires: shared_preload_libraries = 'age' in postgresql.conf
-- 3. Application must connect as superuser to run CREATE EXTENSION

CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS age;
