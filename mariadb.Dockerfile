FROM mariadb:latest
EXPOSE 3306

ENV   MYSQL_ROOT_PASSWORD=admin
ENV   MYSQL_DATABASE=banking
ENV   TZ=America/New_York
