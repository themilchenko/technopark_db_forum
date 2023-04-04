FROM postgres:latest
ENV POSTGRES_PASSWORD=forum
ENV POSTGRES_USER=forum
ENV POSTGRES_DB=forum
COPY db.sql /docker-entrypoint-initdb.d/
EXPOSE 5432
