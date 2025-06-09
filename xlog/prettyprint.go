package xlog

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap/zapcore"
	"reflect"
)

var _prettyLevelToColor = map[zapcore.Level]Color{
	zapcore.DebugLevel:  Magenta,
	zapcore.InfoLevel:   Blue,
	zapcore.WarnLevel:   Yellow,
	zapcore.ErrorLevel:  Red,
	zapcore.DPanicLevel: Red,
	zapcore.PanicLevel:  Red,
	zapcore.FatalLevel:  Red,
}

func printJSON(clr *Color, in any) (err error) {
	rawJSONStr, ok := in.(string)
	if !ok {
		raw, pErr := json.Marshal(in)
		if pErr != nil {
			return pErr
		}
		rawJSONStr = string(raw)
	}
	k := reflect.ValueOf(in).Kind()
	var rawJSON []byte
	if k == reflect.Slice {
		var t []any
		if err = json.Unmarshal([]byte(rawJSONStr), &t); err != nil {
			return err
		}
		rawJSON, err = json.MarshalIndent(&t, "", "\t")
		if err != nil {
			return err
		}
	} else {
		var t map[string]any
		if err = json.Unmarshal([]byte(rawJSONStr), &t); err != nil {
			return err
		}

		rawJSON, err = json.MarshalIndent(&t, "", "\t")
		if err != nil {
			return err
		}
	}
	if clr != nil {
		fmt.Println(clr.Add(string(rawJSON)))
	} else {
		fmt.Println(string(rawJSON))
	}
	return nil
}

const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Color represents a text color.
type Color uint8

// Add adds the coloring to the given string.
func (c Color) Add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}
