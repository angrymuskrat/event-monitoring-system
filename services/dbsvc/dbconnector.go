package dbsvc

import (
	"database/sql"
	"fmt"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	_ "github.com/lib/pq"
)


type DBConnector struct {
	db *sql.DB
}

func NewDBConnector(config string) (*DBConnector, error) {
	db, err := sql.Open("postgres", config)
	if err != nil {
		return nil, ErrDBConnector
	}
	dbc := &DBConnector{db: db}

	_, err = dbc.db.Exec(`CREATE TABLE IF NOT EXISTS posts(
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
		Lat real,
		Lon real 
	)`)
	if err != nil {
		errClose := dbc.db.Close();
		if errClose != nil {

		}
		return nil, ErrDBConnector
	}
	return dbc, nil
}

func (c *DBConnector) Push(posts []data.Post) error {
	sql := `
		INSERT INTO posts (ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, 
			LikesCount, IsAd, AuthorID, LocationID, Lat, Lon)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	for _, v := range posts {
		_, err := c.db.Exec(sql, v.ID, v.Shortcode, v.ImageURL, v.IsVideo, v.Caption, v.CommentsCount, v.Timestamp,
			v.LikesCount, v.IsAd, v.AuthorID, v.LocationID, v.Lat, v.Lon)
		if err != nil {
			return err
		}
	}
	return nil
}

func strBetween(col string, min interface{}, max interface{}) string {
	return fmt.Sprintf("(%v BETWEEN %v AND %v)", col, min, max);
}

func (c DBConnector) Select(irv data.SpatioTemporalInterval) (posts []data.Post, err error) {
	sql := `
		SELECT * 
		FROM posts
		WHERE (Timestamp BETWEEN $1 AND $2)
			AND (Lat BETWEEN $3 AND $4)
			AND (Lon BETWEEN $5 AND $6)
	`

	rows, err := c.db.Query(sql, irv.MinTime, irv.MaxTime, irv.MinLat, irv.MaxLat, irv.MinLon, irv.MaxLon)
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
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
	return err
}
