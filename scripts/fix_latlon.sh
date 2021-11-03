city='london'

PGPASSWORD='cnjhfl;' psql  -U storage -d ${city} -h 10.32.15.30 -c 'UPDATE posts SET location = ST_SetSRID( ST_Point(ST_Y(location), ST_X(location)), 4326);'
PGPASSWORD='cnjhfl;' psql  -U storage -d ${city} -h 10.32.15.30 -c 'UPDATE locations SET position = ST_SetSRID( ST_Point(ST_Y(position), ST_X(position)), 4326);'
PGPASSWORD='cnjhfl;' psql  -U storage -d ${city} -h 10.32.15.30 -c 'UPDATE events SET center = ST_SetSRID( ST_Point(ST_Y(center), ST_X(center)), 4326);'
PGPASSWORD='cnjhfl;' psql  -U storage -d ${city} -h 10.32.15.30 -c 'UPDATE events_3_no_filter SET center = ST_SetSRID( ST_Point(ST_Y(center), ST_X(center)), 4326);'
PGPASSWORD='cnjhfl;' psql  -U storage -d ${city} -h 10.32.15.30 -c 'UPDATE events_6_no_filter SET center = ST_SetSRID( ST_Point(ST_Y(center), ST_X(center)), 4326);'
PGPASSWORD='cnjhfl;' psql  -U storage -d ${city} -h 10.32.15.30 -c 'UPDATE events_12_no_filter SET center = ST_SetSRID( ST_Point(ST_Y(center), ST_X(center)), 4326);'

