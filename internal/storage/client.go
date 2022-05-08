package storage

import (
	"fmt"
	"time"

	"github.com/aburluka/sgtask/internal/logger"
	"github.com/aburluka/sgtask/pkg/model/event"
	"github.com/roistat/go-clickhouse"

	"go.uber.org/zap"
)

const BufSize = 500

type Client struct {
	logger *logger.Logger
	Conn   *clickhouse.Conn

	eventCh chan []*event.Event
}

func New(logger *logger.Logger) (*Client, error) {
	url := fmt.Sprintf("http://127.0.0.1:8123")
	logger.With(zap.String("addr", url)).Info("trying to establish connect with clickouse")
	conn := clickhouse.NewConn(url, clickhouse.NewHttpTransport())

	err := conn.Ping()
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to ping clickhouse")
		return nil, err
	}

	c := &Client{
		logger:  logger,
		Conn:    conn,
		eventCh: make(chan []*event.Event, BufSize*3),
	}
	go c.loop()
	return c, nil
}

func (c *Client) Close() {
}

func (c *Client) loop() {
	t := time.NewTicker(5 * time.Second)
	buf := make([]*event.Event, 0, BufSize*2)
	c.logger.Info("entering storage client work loop")
	for {
		select {
		case e := <-c.eventCh:
			buf = append(buf, e...)
			if len(buf) < BufSize {
				continue
			}
		case <-t.C:
			if len(buf) == 0 {
				c.logger.Info("no events, sleeping")
				continue
			}
		}

		err := c.storeEvents(buf)
		if err != nil {
			c.logger.With(zap.Error(err)).Error("failed to store data in clickhouse")
		}
		c.logger.With(zap.Int("event_number", len(buf))).Info("successfully stored events in clickhouse")
		buf = buf[:0]
	}
}

var columns clickhouse.Columns = clickhouse.Columns{
	"client_time",
	"device_id",
	"device_os",
	"session",
	"sequence",
	"event",
	"param_int",
	"param_str",
	"ip",
	"server_time",
}

func (c *Client) storeEvents(events []*event.Event) error {
	var rows clickhouse.Rows
	for i := range events {
		rows = append(rows,
			clickhouse.Row{
				events[i].ClientTime,
				events[i].DeviceID,
				events[i].DeviceOs,
				events[i].Session,
				events[i].Sequence,
				events[i].Event,
				events[i].ParamInt,
				events[i].ParamStr,
				events[i].IP,
				events[i].ServerTime,
			},
		)
	}

	query, err := clickhouse.BuildMultiInsert("sg.events", columns, rows)
	if err != nil {
		c.logger.With(zap.Error(err), zap.String("table", "sg.events")).Error("clickhouse.BuildMultiInsert failed")
		return err
	}

	err = query.Exec(c.Conn)
	if err != nil {
		c.logger.With(zap.Error(err)).Error("query.Exec failed")
		return err
	}

	return nil
}

func (c *Client) PutEvents(events []*event.Event) {
	c.eventCh <- events
}
