CREATE DATABASE snippets;
CREATE USER 'web'@'localhost' IDENTIFIED BY 'pass';
GRANT ALL PRIVILEGES ON *.* TO 'web'@'localhost' WITH GRANT OPTION;
