/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
 */

package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var bufPool = newBufPool()
var allowedFields = []string{
	zerolog.TimestampFieldName,
	zerolog.LevelFieldName,
	moduleNameKey,
	loggerNameKey,
	zerolog.CallerFieldName,
	zerolog.MessageFieldName,
}

// Writer for log.
type Writer struct {
	out        io.Writer
	timeFormat string
}

// Write write.
func (w *Writer) Write(in []byte) (int, error) {
	fields := make(map[string]interface{})
	decoder := json.NewDecoder(bytes.NewReader(in))
	decoder.UseNumber()
	if err := decoder.Decode(&fields); err != nil {
		return 0, fmt.Errorf("cannot decode fields: %s", err)
	}

	buf := bufPool.GetBuf()
	defer bufPool.PutBuf(buf)
	w.writeFields(buf, fields)
	buf.WriteByte('\n')
	_, err := buf.WriteTo(w.out)

	return len(in), err
}

func (w *Writer) writeFields(buf io.StringWriter, fields map[string]interface{}) {
	for i, allowedField := range allowedFields {
		field, ok := fields[allowedField]
		if !ok {
			continue
		}
		if i > 0 {
			_, err := buf.WriteString(" | ")
			if err != nil {
				continue
			}
		}

		switch allowedField {
		case zerolog.TimestampFieldName:
			w.writeField(buf, field)
		case zerolog.LevelFieldName:
			_, err := buf.WriteString(strings.ToUpper(field.(string)))
			if err != nil {
				return
			}
		default:
			_, err := buf.WriteString(field.(string))
			if err != nil {
				return
			}
		}
	}
}

func (w *Writer) writeField(buf io.StringWriter, field interface{}) {
	timestamp, err := parseTimestamp(field)
	if err != nil || timestamp == nil {
		_, err := buf.WriteString("<nil>")
		if err != nil {
			return
		}
	} else {
		_, err := buf.WriteString(timestamp.Format(w.timeFormat))
		if err != nil {
			return
		}
	}
}

func parseTimestamp(timeIn interface{}) (*time.Time, error) {
	switch t := timeIn.(type) {
	case string:
		timestamp, err := time.Parse(zerolog.TimeFieldFormat, t)
		if err != nil {
			return nil, err
		}
		return &timestamp, nil
	case json.Number:
		timeInteger, err := t.Int64()
		if err != nil {
			return nil, err
		}
		var sec, nsec int64 = timeInteger, 0
		switch zerolog.TimeFieldFormat {
		case zerolog.TimeFormatUnixMs:
			nsec = int64(time.Duration(timeInteger) * time.Millisecond)
			sec = 0
		case zerolog.TimeFormatUnixMicro:
			nsec = int64(time.Duration(timeInteger) * time.Microsecond)
			sec = 0
		}
		timestamp := time.Unix(sec, nsec)
		return &timestamp, nil
	default:
		return nil, nil
	}
}

type buffersPool struct {
	pool *sync.Pool
}

// newBufPool new buffer pool.
func newBufPool() *buffersPool {
	pool := &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 100))
		},
	}
	return &buffersPool{pool: pool}
}

// GetBuf get buffer.
func (b *buffersPool) GetBuf() *bytes.Buffer {
	return b.pool.Get().(*bytes.Buffer)
}

// PutBuf put buffer.
func (b *buffersPool) PutBuf(buf *bytes.Buffer) {
	buf.Reset()
	b.pool.Put(buf)
}
