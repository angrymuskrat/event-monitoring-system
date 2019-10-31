package dbsvc

import (
	"database/sql"
	"fmt"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/lib/pq"
)

type DBConnector struct {
	db     *sql.DB
	logger log.Logger
}

func NewDBConnector(config string, logger log.Logger) (*DBConnector, error) {
	db, err := sql.Open("postgres", config)
	if err != nil {
		return nil, err
	}
	dbc := &DBConnector{db: db, logger: log.With(logger, "dbConnector")}

	create := `CREATE TABLE IF NOT EXISTS posts(
		ID varchar (120) not null primary key,
		Shortcode varchar (120),
		ImageURL varchar (120),
		IsVideo boolean not null,
		Caption varchar (120),
		CommentsCount bigint,
		Timestamp bigint,
		LikesCount bigint,
		IsAd boolean,
		AuthorID varchar (120),
		LocationID varchar (120),
		--Lat real,
		--Lon real
		Location geometry 
	)`

	_, err = dbc.db.Exec(create)
	if err != nil {
		errClose := dbc.db.Close();
		if errClose != nil {
			level.Error(dbc.logger).Log("err", errClose.Error())
		}
		return nil, err
	}
	return dbc, nil
}

func (c *DBConnector) Push(posts []data.Post) error {
	err := c.db.Ping()
	if err != nil {
		level.Error(c.logger).Log("errPush", err);
		return ErrDBConnecting
	}

	for _, v := range posts {
		insert := fmt.Sprintf(`
			INSERT INTO posts (ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, 
				AuthorID, LocationID, Location)
			VALUES ('%v', '%v', '%v', %v, '%v', %v, %v, %v, %v, '%v', '%v', ST_GeometryFromText('POINT(%v %v)'));`,
			v.ID, v.Shortcode, v.ImageURL, v.IsVideo, v.Caption, v.CommentsCount, v.Timestamp,
			v.LikesCount, v.IsAd, v.AuthorID, v.LocationID, v.Lat, v.Lon)
		_, err = c.db.Exec(insert);
		if err != nil {
			level.Error(c.logger).Log("errPush", err);
		}
	}

	return err
}

func (c DBConnector) Select(irv data.SpatioTemporalInterval) (posts []data.Post, err error) {
	poly := fmt.Sprintf("ST_GeometryFromText('POLYGON((%v %v, %v %v, %v %v, %v %v, %v %v))')",
		irv.MinLat, irv.MinLon, irv.MaxLat, irv.MinLon, irv.MaxLat, irv.MaxLon, irv.MinLat, irv.MaxLon,
		irv.MinLat, irv.MinLon)

	statement := fmt.Sprintf(`
		SELECT ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, 
			LocationID, ST_X(Location) as Lat, ST_Y(Location) as Lon
		FROM posts
		WHERE ST_Contains(%v, Location) AND (Timestamp BETWEEN %v AND %v)
	`, poly, irv.MinTime, irv.MaxTime)


	err = c.db.Ping()
	if err != nil {
		level.Error(c.logger).Log("errSelect", err);
		return nil, ErrDBConnecting
	}

	rows, err := c.db.Query(statement)
	if err != nil {
		return nil, err;
	}

	defer func() {
		errClose := rows.Close();
		if errClose != nil {
			level.Error(c.logger).Log("err", errClose.Error())
		}
	}()

	for rows.Next() {
		p := new(data.Post)
		//(ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, LocationID, Location)
		err = rows.Scan(&p.ID, &p.Shortcode, &p.ImageURL, &p.IsVideo, &p.Caption, &p.CommentsCount, &p.Timestamp,
			&p.LikesCount, &p.IsAd, &p.AuthorID, &p.LocationID, &p.Lat, &p.Lon)
		if err != nil {
			return nil, err
		}
		posts = append(posts, *p)
	}
	return posts, nil
}

func (c *DBConnector) Close() error {
	err := c.db.Close();
	if err != nil {
		level.Error(c.logger).Log("err", err.Error())
	}
	return err
}
