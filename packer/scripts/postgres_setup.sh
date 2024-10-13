!#/bin/bash

sudo apt install postgresql postgresql-contrib -y

sudo systemctl start postgresql
sudo systemctl enable postgresql

sudo -u postgres psql -c "CREATE DATABASE webapp;"
sudo -u postgres psql -c "CREATE USER sid WITH PASSWORD 'sidd';"
sudo -u postgres psql -c "ALTER DATABASE webapp OWNER TO sid;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE webapp TO sid;"

