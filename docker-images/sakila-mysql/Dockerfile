# https://daveb.medium.com/using-a-pre-seeded-mysqlql-docker-image-for-faster-ci-cd-ccf35305e996
FROM debian:buster-slim AS sakila
# Download the Sakila MySQL example database
RUN apt-get update && apt-get install -y --no-install-recommends wget ca-certificates unzip
RUN wget https://github.com/jOOQ/sakila/archive/refs/heads/main.zip && unzip main.zip 'sakila-main/mysql-sakila-db/*'
# Copy the schema and data .sql files from the Sakila database into
# /docker-entrypoint-initdb.d so they are executed when the
# container initializes
FROM mysql:8.0 AS sakila-preseed
COPY --from=sakila /sakila-main/mysql-sakila-db/mysql-sakila-schema.sql /docker-entrypoint-initdb.d/01-mysql-sakila-schema.sql
COPY --from=sakila /sakila-main/mysql-sakila-db/mysql-sakila-insert-data.sql /docker-entrypoint-initdb.d/02-mysql-sakila-insert-data.sql
# Modify MySQL's original docker-entrypoint.sh file slightly so
# that MySQL doesn't attempt to start back up after the
#  "initdb" process is complete
#RUN grep -v 'exec "$@"' /usr/local/bin/docker-entrypoint.sh > /docker-entrypoint.sh && chmod 755 /docker-entrypoint.sh
# "trust" should never be used in a production setting, but should
# be safe enough for this example
ENV MYSQL_ROOT_PASSWORD password
ENV MYSQL_DATABASE sakila
# Run mysql, it will stop after the .sql files have been installed
# into the database
RUN /entrypoint.sh mysqld & sleep 160 && mysqladmin --password=password shutdown
# Copy just the pre-seeded DATA directory from the previous
# command in our new image
FROM mysql:8.0
RUN sed -i 's|/var/lib/mysql|/var/lib/mysql2|g' /etc/my.cnf
COPY --from=sakila-preseed /var/lib/mysql /var/lib/mysql2