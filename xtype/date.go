package xtype

import (
	"encoding/json"
	"github.com/kurzgesagtz/xgo/internal/ptr"
	"strings"
	"time"
)

type Date time.Time

func (d *Date) Time() *time.Time {
	if d == nil {
		return nil
	}
	return ptr.Time(time.Time(*d))
}

func (d *Date) MarshalJSON() ([]byte, error) {
	if d == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(d.Time())
}

func (d *Date) UnmarshalJSON(data []byte) (err error) {
	raw := strings.TrimSuffix(strings.TrimPrefix(string(data), "\""), "\"")
	var t time.Time
	t, err = time.Parse(time.RFC3339, raw)
	if err != nil {
		t, err = time.Parse("2006-01-02", raw)
		if err != nil {
			return err
		}
	}
	*d = (Date)(t)
	return nil
}
