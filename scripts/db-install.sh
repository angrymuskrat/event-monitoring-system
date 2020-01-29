# for ubuntu

# install postgres 10
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt/ $(lsb_release -sc)-pgdg main" > /etc/apt/sources.list.d/PostgreSQL.list'
sudo apt update
sudo apt-get install postgresql-10 postgresql-server-dev-10 -y

# instal requerments of postgis
sudo apt install gcc
sudo apt install libproj-dev proj-data proj-bin -y
sudo apt install libgeos-dev -y
sudo apt install libxml2-dev -y
sudo apt install libjson0 libjson0-dev -y
sudo add-apt-repository ppa:ubuntugis/ppa && sudo apt-get update
sudo apt install gdal-bin libgdal-dev -y


# instal postgis
curl -O https://download.osgeo.org/postgis/source/postgis-2.4.8.tar.gz
tar xvfz postgis-2.4.8.tar.gz
cd postgis-2.4.8/ || exit
./configure
make
sudo make install
sudo make clean
cd ..
sudo rm -r postgis-2.4.8/
rm postgis-2.4.8.tar.gz

# instal timescaleDB
# if there is problem here use:  sudo apt install software-properties-common
sudo add-apt-repository ppa:timescale/timescaledb-ppa
sudo apt-get update

# Now install appropriate package for PG version
sudo apt install timescaledb-postgresql-10 -y
sudo timescaledb-tune
sudo service postgresql restart
