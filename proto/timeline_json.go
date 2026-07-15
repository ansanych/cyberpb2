package proto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// MarshalJSON implements json.Marshaler for Timeline.
// Сериализует dataValues как массив массивов [[1,2,3], [4,5,6]]
// вместо [{"values": [1,2,3]}, {"values": [4,5,6]}].
// Пустые поля (connections, engines, jobs, errors, finished=false) не выводятся.
func (tl *Timeline) MarshalJSON() ([]byte, error) {
	if tl == nil {
		return []byte("null"), nil
	}

	var buf bytes.Buffer
	buf.WriteByte('{')

	// sn
	if tl.Sn != "" {
		writeStringJSON(&buf, "sn", tl.Sn)
	}

	// date
	if tl.Date != "" {
		writeComma(&buf, buf.Len() > 2)
		writeStringJSON(&buf, "date", tl.Date)
	}

	// tz
	if tl.Tz != 0 {
		writeComma(&buf, buf.Len() > 2)
		fmt.Fprintf(&buf, `"tz":%d`, tl.Tz)
	}

	// startTimeUnix
	if tl.StartTimeUnix != 0 {
		writeComma(&buf, buf.Len() > 2)
		fmt.Fprintf(&buf, `"startTimeUnix":%d`, tl.StartTimeUnix)
	}

	// params
	if len(tl.Params) > 0 {
		writeComma(&buf, buf.Len() > 2)
		buf.WriteString(`"params":`)
		writeStringSliceJSON(&buf, tl.Params)
	}

	// connections
	writeEventBlocksJSON(&buf, "connections", tl.Connections, buf.Len() > 2)

	// engines
	writeEventBlocksJSON(&buf, "engines", tl.Engines, buf.Len() > 2)

	// jobs
	writeEventBlocksJSON(&buf, "jobs", tl.Jobs, buf.Len() > 2)

	// errors
	writeEventBlocksJSON(&buf, "errors", tl.Errors, buf.Len() > 2)

	// data
	writeConnectionDataBlocksJSON(&buf, "data", tl.Data, buf.Len() > 2)

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// writeComma пишет запятую перед полем, если нужно.
func writeComma(buf *bytes.Buffer, needComma bool) {
	if needComma {
		buf.WriteByte(',')
	}
}

// writeStringJSON пишет "key":"value" с JSON-экранированием.
func writeStringJSON(buf *bytes.Buffer, key, val string) {
	fmt.Fprintf(buf, `"%s":`, key)
	b, _ := json.Marshal(val)
	buf.Write(b)
}

// writeEventBlocksJSON пишет "name":[...] для EventBlock, если массив непустой.
func writeEventBlocksJSON(buf *bytes.Buffer, name string, blocks []*EventBlock, needComma bool) {
	if len(blocks) == 0 {
		return
	}
	writeComma(buf, needComma)
	fmt.Fprintf(buf, `"%s":[`, name)
	for i, b := range blocks {
		if i > 0 {
			buf.WriteByte(',')
		}
		writeEventBlockJSON(buf, b)
	}
	buf.WriteByte(']')
}

// writeEventBlockJSON пишет {"start":...,"end":...,"value":...}.
func writeEventBlockJSON(buf *bytes.Buffer, b *EventBlock) {
	buf.WriteByte('{')
	fmt.Fprintf(buf, `"start":%d,"end":%d,"value":%s`, b.Start, b.End, formatFloat(float64(b.Value)))
	buf.WriteByte('}')
}

// writeConnectionDataBlocksJSON пишет "name":[...] для ConnectionDataBlock, если массив непустой.
func writeConnectionDataBlocksJSON(buf *bytes.Buffer, name string, blocks []*ConnectionDataBlock, needComma bool) {
	if len(blocks) == 0 {
		return
	}
	writeComma(buf, needComma)
	fmt.Fprintf(buf, `"%s":[`, name)
	for i, b := range blocks {
		if i > 0 {
			buf.WriteByte(',')
		}
		writeConnectionDataBlockJSON(buf, b)
	}
	buf.WriteByte(']')
}

// writeConnectionDataBlockJSON пишет ConnectionDataBlock с dataValues как массив массивов.
func writeConnectionDataBlockJSON(buf *bytes.Buffer, b *ConnectionDataBlock) {
	buf.WriteByte('{')

	// id
	fmt.Fprintf(buf, `"id":%d`, b.Id)

	// finished (только true)
	if b.Finished {
		buf.WriteString(`,"finished":true`)
	}

	// dataValues — как массив массивов
	if len(b.DataValues) > 0 {
		buf.WriteString(`,"dataValues":[`)
		for j, row := range b.DataValues {
			if j > 0 {
				buf.WriteByte(',')
			}
			buf.WriteByte('[')
			for k, v := range row.Values {
				if k > 0 {
					buf.WriteByte(',')
				}
				buf.WriteString(formatFloat(float64(v)))
			}
			buf.WriteByte(']')
		}
		buf.WriteByte(']')
	}

	buf.WriteByte('}')
}

// formatFloat форматирует float64 в JSON-совместимую строку
// без лишних знаков (избегает проблемы 4340.60009765625 вместо 4340.6).
func formatFloat(v float64) string {
	if math.IsInf(v, 0) || math.IsNaN(v) {
		b, _ := json.Marshal(v)
		return string(b)
	}
	// strconv.FormatFloat с 'g' и точностью -1 даёт кратчайшее точное представление
	return strconv.FormatFloat(v, 'g', -1, 32)
}

// writeStringSliceJSON пишет ["a","b","c"] с JSON-экранированием.
func writeStringSliceJSON(buf *bytes.Buffer, strs []string) {
	buf.WriteByte('[')
	for i, s := range strs {
		if i > 0 {
			buf.WriteByte(',')
		}
		b, _ := json.Marshal(s)
		buf.Write(b)
	}
	buf.WriteByte(']')
}
