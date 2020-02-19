docker run -d --name timescaledb -p 5432:5432 -e POSTGRES_USER=user \
-e POSTGRES_PASSWORD=pwd -e TS_TUNE_MEMORY=6GB -e TS_TUNE_NUM_CPUS=2 \
timescale/timescaledb-postgis:latest-pg10