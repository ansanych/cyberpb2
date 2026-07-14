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
// Пустые поля (engine, job, errors, finished=false) не выводятся.
func (tl *Timeline) MarshalJSON() ([]byte, error) {
	if tl == nil {
		return []byte("null"), nil
	}

	var buf bytes.Buffer
	buf.WriteByte('{')

	// params (если есть)
	if len(tl.Params) > 0 {
		buf.WriteString(`"params":`)
		writeStringSliceJSON(&buf, tl.Params)
	}

	// data
	if len(tl.Data) > 0 {
		if len(tl.Params) > 0 {
			buf.WriteByte(',')
		}
		writeConnectionsJSON(&buf, "data", tl.Data)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func writeConnectionsJSON(buf *bytes.Buffer, name string, conns []*TelemetryConnection) {
	fmt.Fprintf(buf, `"%s":[`, name)
	for i, c := range conns {
		if i > 0 {
			buf.WriteByte(',')
		}
		writeConnectionJSON(buf, c)
	}
	buf.WriteByte(']')
}

func writeConnectionJSON(buf *bytes.Buffer, c *TelemetryConnection) {
	buf.WriteByte('{')

	// id
	fmt.Fprintf(buf, `"id":%d`, c.Id)

	// start
	if c.Start != 0 {
		fmt.Fprintf(buf, `,"start":%d`, c.Start)
	}

	// end
	if c.End != 0 {
		fmt.Fprintf(buf, `,"end":%d`, c.End)
	}

	// finished (только true)
	if c.Finished {
		buf.WriteString(`,"finished":true`)
	}

	// engine
	writeBlocksJSON(buf, "engine", c.Engine)

	// job
	writeBlocksJSON(buf, "job", c.Job)

	// errors
	writeBlocksJSON(buf, "errors", c.Errors)

	// dataValues — как массив массивов
	if len(c.DataValues) > 0 {
		buf.WriteString(`,"dataValues":[`)
		for j, row := range c.DataValues {
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

func writeBlocksJSON(buf *bytes.Buffer, name string, blocks []*TelemetryBlock) {
	if len(blocks) == 0 {
		return
	}
	fmt.Fprintf(buf, `,"%s":[`, name)
	for i, b := range blocks {
		if i > 0 {
			buf.WriteByte(',')
		}
		fmt.Fprintf(buf, `{"start":%d,"end":%d,"value":%s}`, b.Start, b.End, formatFloat(float64(b.Value)))
	}
	buf.WriteByte(']')
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

func writeStringSliceJSON(buf *bytes.Buffer, strs []string) {
	buf.WriteByte('[')
	for i, s := range strs {
		if i > 0 {
			buf.WriteByte(',')
		}
		// Экранирование для JSON
		b, _ := json.Marshal(s)
		buf.Write(b)
	}
	buf.WriteByte(']')
}
