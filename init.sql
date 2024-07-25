DO
$do$
BEGIN
   IF NOT EXISTS (
      SELECT
      FROM   pg_catalog.pg_database
      WHERE  datname = 'messagio') THEN

      PERFORM dblink_exec('dbname=postgres', 'CREATE DATABASE messagio');
   END IF;
END
$do$;