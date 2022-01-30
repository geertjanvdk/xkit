// Copyright (c) 2019, 2021 Geert JM Vanderkelen

package xlog

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/geertjanvdk/xkit/xansi"
)

type TextFormatType int

const (
	TextFullFields TextFormatType = iota + 1
	TextCompat
)

type TextFormat struct {
	FormatType TextFormatType
	TimeFormat string // defaults to  (tf *TextFormat)
}

var (
	styleBlueBold  = xansi.Render{xansi.Blue, xansi.Bold}
	styleRedBold   = xansi.Render{xansi.Red, xansi.Bold}
	styleBlack     = xansi.Render{xansi.Black}
	styleBlackBold = xansi.Render{xansi.Black, xansi.Bold}
	styleGray      = xansi.Render{xansi.BrightBlack}
)

func (tf *TextFormat) fullFields(e *Entry) ([]byte, error) {
	var buf = bytes.Buffer{}

	// built-in fields come first
	writeField(&buf, FieldTime, tf.serializeFieldValue(e.Time))
	writeField(&buf, FieldLevel, tf.serializeFieldValue(e.Level))
	if e.Scope != "" {
		writeField(&buf, FieldScope, tf.serializeFieldValue(e.Scope))
	}
	writeField(&buf, FieldMsg, tf.serializeFieldValue(e.message))
	if e.ErrCode != "" {
		writeField(&buf, FieldErrCode, tf.serializeFieldValue(e.ErrCode))
	}

	tf.writeFields(&buf, e)

	buf.WriteRune('\n')
	return buf.Bytes(), nil
}

func (tf *TextFormat) textCompact(e *Entry) ([]byte, error) {
	var buf = bytes.Buffer{}

	tf.TimeFormat = time.RFC3339

	writeField(&buf, "", styleBlueBold.Sprintf(tf.serializeFieldValue(e.Time)))

	levelStyle := styleBlack
	msgStyle := styleBlack
	switch e.Level {
	case ErrorLevel:
		levelStyle = styleRedBold
		msgStyle = styleBlackBold
	}

	writeField(&buf, "", levelStyle.Sprintf("[%-5s]", strings.ToUpper(tf.serializeFieldValue(e.Level))))
	buf.WriteString(xansi.Reset())
	writeField(&buf, "", msgStyle.Sprintf("%-*s |", 50, e.message))
	buf.WriteString(xansi.Reset())

	if e.Scope != "" {
		writeField(&buf, "scope", styleGray.Sprintf(e.Scope))
	}

	tf.writeFields(&buf, e)

	buf.WriteRune('\n')
	return buf.Bytes(), nil
}

func (tf *TextFormat) writeFields(buf *bytes.Buffer, e *Entry) {
	for k, v := range e.Fields {
		if reservedFields[k] {
			k = "_" + k
		}
		writeField(buf, k, tf.serializeFieldValue(v))
	}
}

func (tf *TextFormat) Format(e *Entry) ([]byte, error) {
	if tf.TimeFormat == "" {
		tf.TimeFormat = time.RFC3339Nano
	}

	switch tf.FormatType {
	case TextCompat:
		return tf.textCompact(e)
	case TextFullFields:
		fallthrough
	default:
		return tf.fullFields(e)
	}
}

func writeField(buf *bytes.Buffer, key, value string) {
	if key == "" {
		buf.WriteString(value + " ")
	} else {
		buf.WriteString(key + "=" + value + " ")
	}
}

func (tf *TextFormat) serializeFieldValue(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return strconv.Quote(v)
	case []byte:
		return strconv.Quote(string(v))
	case Level:
		return levelName[v]
	case time.Time:
		return v.Format(tf.TimeFormat)
	case time.Duration:
		return v.String()
	case bool, *bool:
		return fmt.Sprintf("%v", v)
	default:
		switch v := numTo64(value).(type) {
		case int64:
			return strconv.FormatInt(v, 10)
		case uint64:
			return strconv.FormatUint(v, 10)
		case float64:
			return strconv.FormatFloat(v, 'E', -1, 64)
		}
	}

	// last resort
	return strconv.Quote(fmt.Sprintf("%q", value))
}

func numTo64(n interface{}) interface{} {
	switch n := n.(type) {
	case int:
		return int64(n)
	case uint:
		return uint64(n)
	case int8:
		return int64(n)
	case int16:
		return int64(n)
	case int32:
		return int64(n)
	case int64:
		return int64(n)
	case uint8:
		return uint64(n)
	case uint16:
		return uint64(n)
	case uint32:
		return uint64(n)
	case uint64:
		return uint64(n)
	case float32:
		return float64(n)
	case float64:
		return float64(n)
	default:
		return nil
	}
}
