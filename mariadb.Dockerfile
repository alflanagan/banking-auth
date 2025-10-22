FROM mariadb:latest
EXPOSE 3307

ENV
      MYSQL_ROOT_PASSWORD=admin
      MYSQL_DATABASE=banking
      TZ=America/New_York
