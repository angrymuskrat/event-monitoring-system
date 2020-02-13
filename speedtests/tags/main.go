package main

import (
	"context"
	"fmt"
	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/storage"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/angrymuskrat/event-monitoring-system/utils/rand"
	"github.com/jackc/pgx/v4"
	"time"
)

type Storage struct {
	conn *pgx.Conn
}

type Test struct {
	Events []data.Event
	Time   int64
}

const AuthToken = "database=%v user=myuser password=mypass sslmode=disable host=localhost port=5432"
const NameTestDB = "test_tags"

func NewStorage(ctx context.Context) (*Storage, error) {
	connConfig, err := pgx.ParseConfig(fmt.Sprintf(AuthToken, storage.PostgresDBName))
	if err != nil {
		return nil, err
	}
	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}

	//tmp
	/*_, err = conn.Exec(ctx, "DROP DATABASE " + NameTestDB)
	if err != nil {
		return nil, err
	}*/

	name := NameTestDB
	row := conn.QueryRow(ctx, storage.MakeSelectDB(name))
	err = row.Scan(&name)
	if err == pgx.ErrNoRows {
		_, err = conn.Exec(ctx, storage.MakeCreateDB(name))
	}
	if err != nil {
		return nil, err
	}
	connConfig, err = pgx.ParseConfig(fmt.Sprintf(AuthToken, name))
	conn, err = pgx.ConnectConfig(ctx, connConfig)
	st := &Storage{conn: conn}
	err = st.setEnvironment(ctx)
	if err != nil {
		conn.Close(ctx)
		return nil, err
	}
	return st, nil
}

func (s *Storage) setEnvironment(ctx context.Context) (err error) {
	_, err = s.conn.Exec(ctx, storage.ExtensionPostGIS)
	if err != nil {
		return err
	}
	_, err = s.conn.Exec(ctx, storage.ExtensionPostGISTopology)
	if err != nil {
		return err
	}
	_, err = s.conn.Exec(ctx, storage.ExtensionTimescaleDB)
	if err != nil {
		return err
	}
	_, err = s.conn.Exec(ctx, storage.CreateTimeFunction)
	if err != nil {
		return err
	}
	_, err = s.conn.Exec(ctx, CreateEventsTable)
	if err != nil {
		return err
	}
	_, err = s.conn.Exec(ctx, CreateHyperTableEvents)
	if err != nil {
		return err
	}

	_, err = s.conn.Exec(ctx, CreateEventTagsTables1)
	if err != nil {
		return err
	}
	_, err = s.conn.Exec(ctx, CreateTagsTables2)
	if err != nil {
		return err
	}
	_, err = s.conn.Exec(ctx, CreateEventTagsTables2)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteAll(ctx context.Context) (err error) {
	_, err = s.conn.Exec(ctx, "DELETE FROM event_tags1;")
	if err != nil {
		return err
	}
	_, err = s.conn.Exec(ctx, "DELETE FROM event_tags2;")
	if err != nil {
		return err
	}
	_, err = s.conn.Exec(ctx, "DELETE FROM tags2;")
	if err != nil {
		return err
	}
	_, err = s.conn.Exec(ctx, "DELETE FROM events;")
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) insertEvents(ctx context.Context, events []data.Event) (err error) {
	tr, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}
	for _, e := range events {
		var eventId int
		row := tr.QueryRow(ctx, InsertEvent, e.Title, e.Start, e.Finish, e.Center.Lat, e.Center.Lon, e.PostCodes, e.Tags)
		err = row.Scan(&eventId)
		if err != nil {
			return err
		}
		for _, tag := range e.Tags {
			_, err = tr.Exec(ctx, InsertEventTags1, eventId, e.Start, e.Finish, tag)
			if err != nil {
				return err
			}
			var tagId int
			row := tr.QueryRow(ctx, InsertTags2, tag)
			err = row.Scan(&tagId)
			if err != nil {
				return err
			}
			_, err = tr.Exec(ctx, InsertEventTags2, tagId, eventId, e.Start, e.Finish)
			if err != nil {
				return err
			}
		}
	}
	if err = tr.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Storage) FillEvents(ctx context.Context, amount int, conf rand.GenConfig) (err error) {
	events := rand.New().Events(amount, conf)
	delta := 10000

	for ind := 0; ind < amount; ind += delta {
		var slice []data.Event
		if ind+delta < amount {
			slice = events[ind : ind+delta]
		} else {
			slice = events[ind:amount]
		}
		err = s.insertEvents(ctx, slice)
		if err != nil {
			return err
		}
		fmt.Printf("Insert: %v events\n", ind+len(slice))
	}
	return nil
}

func (s *Storage) GetEvents(ctx context.Context, sql string) (ans *Test, err error) {
	begin := time.Now()
	var events []data.Event
	rows, err := s.conn.Query(ctx, sql)
	if err != nil {
		return nil,  err
	}
	for rows.Next() {
		e := data.Event{}
		p := data.Point{}
		err = rows.Scan(&e.Title, &e.Start, &e.Finish, &e.PostCodes, &e.Tags, &p.Lat, &p.Lon)
		if err != nil {
			return nil, err
		}
		e.Center = p
		events = append(events, e)
	}
	return &Test{Events: events, Time: time.Since(begin).Milliseconds()}, nil
}

func main() {
	st, err := NewStorage(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	NY := data.Point{Lat: 40.7, Lon: -74}
	now := int64(time.Now().UTC().Unix())
	fmt.Println(now)
	genConf := rand.GenConfig{Center: NY, DeltaPoint: data.Point{Lat: 0.01, Lon: 0.002}, StartTime: now - 24*3600, FinishTime: now}
	amount := 10

	mode := "test"

	switch mode {
	case "init":
		err := st.FillEvents(context.Background(), amount, genConf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("events successfully was generated")
		break
	case "test":
		tags := []string{"n", "b"}
		start := int64(1581586063)
		finish := int64(1581593063)
		sqls := []string{makeSelectEvents1(tags, start, finish), makeSelectEvents2(tags, start, finish)}

		for i, sql := range sqls {
			test, err := st.GetEvents(context.Background(), sql)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("№%v time: %v size: %v\n", i, test.Time, len(test.Events))
		}

		break
	default:
		break
	}

}
