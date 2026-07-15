package proto

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTimeline_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		timeline *Timeline
		want     string
	}{
		{
			name:     "nil timeline",
			timeline: nil,
			want:     "null",
		},
		{
			name:     "empty timeline",
			timeline: &Timeline{},
			want:     `{}`,
		},
		{
			name: "only sn and date",
			timeline: &Timeline{
				Sn:   "ABC123",
				Date: "2024-01-15",
			},
			want: `{"sn":"ABC123","date":"2024-01-15"}`,
		},
		{
			name: "with params",
			timeline: &Timeline{
				Sn:     "ABC123",
				Params: []string{"rpm", "temp"},
			},
			want: `{"sn":"ABC123","params":["rpm","temp"]}`,
		},
		{
			name: "with connections",
			timeline: &Timeline{
				Sn: "ABC123",
				Connections: []*EventBlock{
					{Start: 1000, End: 2000, Value: 1.5},
					{Start: 2000, End: 3000, Value: 2.5},
				},
			},
			want: `{"sn":"ABC123","connections":[{"start":1000,"end":2000,"value":1.5},{"start":2000,"end":3000,"value":2.5}]}`,
		},
		{
			name: "with engines and jobs",
			timeline: &Timeline{
				Sn: "ABC123",
				Engines: []*EventBlock{
					{Start: 1000, End: 1500, Value: 800},
				},
				Jobs: []*EventBlock{
					{Start: 1000, End: 2000, Value: 1},
				},
			},
			want: `{"sn":"ABC123","engines":[{"start":1000,"end":1500,"value":800}],"jobs":[{"start":1000,"end":2000,"value":1}]}`,
		},
		{
			name: "with data and dataValues as array of arrays",
			timeline: &Timeline{
				Sn: "ABC123",
				Data: []*ConnectionDataBlock{
					{
						Id: 1,
						DataValues: []*DataValuesRow{
							{Values: []float32{1.0, 2.0, 3.0}},
							{Values: []float32{4.0, 5.0, 6.0}},
						},
					},
				},
			},
			want: `{"sn":"ABC123","data":[{"id":1,"dataValues":[[1,2,3],[4,5,6]]}]}`,
		},
		{
			name: "with finished=true in data",
			timeline: &Timeline{
				Sn: "ABC123",
				Data: []*ConnectionDataBlock{
					{
						Id:       1,
						Finished: true,
					},
				},
			},
			want: `{"sn":"ABC123","data":[{"id":1,"finished":true}]}`,
		},
		{
			name: "full timeline",
			timeline: &Timeline{
				Sn:            "ABC123",
				Date:          "2024-01-15",
				Tz:            5,
				StartTimeUnix: 1700000000,
				Params:        []string{"rpm", "temp"},
				Connections: []*EventBlock{
					{Start: 1000, End: 2000, Value: 1},
				},
				Engines: []*EventBlock{
					{Start: 1000, End: 1500, Value: 800},
				},
				Jobs: []*EventBlock{
					{Start: 1000, End: 2000, Value: 1},
				},
				Errors: []*EventBlock{
					{Start: 1500, End: 1600, Value: 1},
				},
				Data: []*ConnectionDataBlock{
					{
						Id:       1,
						Finished: true,
						DataValues: []*DataValuesRow{
							{Values: []float32{1.0, 2.0}},
							{Values: []float32{3.0, 4.0}},
						},
					},
				},
			},
			want: `{"sn":"ABC123","date":"2024-01-15","tz":5,"startTimeUnix":1700000000,"params":["rpm","temp"],"connections":[{"start":1000,"end":2000,"value":1}],"engines":[{"start":1000,"end":1500,"value":800}],"jobs":[{"start":1000,"end":2000,"value":1}],"errors":[{"start":1500,"end":1600,"value":1}],"data":[{"id":1,"finished":true,"dataValues":[[1,2],[3,4]]}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.timeline)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("json.Marshal() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestTimeline_MarshalJSON_EmptyFields(t *testing.T) {
	// Проверяем, что пустые поля не выводятся
	timeline := &Timeline{
		Sn:          "ABC123",
		Connections: []*EventBlock{},
		Engines:     nil,
		Jobs:        []*EventBlock{},
		Errors:      nil,
		Data:        []*ConnectionDataBlock{},
	}

	got, err := json.Marshal(timeline)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	want := `{"sn":"ABC123"}`
	if string(got) != want {
		t.Errorf("json.Marshal() = %v, want %v", string(got), want)
	}
}

func TestTimeline_MarshalJSON_DataValuesEmpty(t *testing.T) {
	// Проверяем, что пустой dataValues не выводится
	timeline := &Timeline{
		Sn: "ABC123",
		Data: []*ConnectionDataBlock{
			{
				Id:         1,
				DataValues: []*DataValuesRow{},
			},
		},
	}

	got, err := json.Marshal(timeline)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	want := `{"sn":"ABC123","data":[{"id":1}]}`
	if string(got) != want {
		t.Errorf("json.Marshal() = %v, want %v", string(got), want)
	}
}

func TestTimeline_MarshalJSON_DataValuesWithEmptyRow(t *testing.T) {
	// Проверяем, что строка с пустым Values выводится как []
	timeline := &Timeline{
		Sn: "ABC123",
		Data: []*ConnectionDataBlock{
			{
				Id: 1,
				DataValues: []*DataValuesRow{
					{Values: []float32{}},
					{Values: []float32{1.0}},
				},
			},
		},
	}

	got, err := json.Marshal(timeline)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	want := `{"sn":"ABC123","data":[{"id":1,"dataValues":[[],[1]]}]}`
	if string(got) != want {
		t.Errorf("json.Marshal() = %v, want %v", string(got), want)
	}
}

func TestTimeline_MarshalJSON_FloatPrecision(t *testing.T) {
	// Проверяем, что float32 4340.6 не превращается в 4340.60009765625
	timeline := &Timeline{
		Sn: "ABC123",
		Data: []*ConnectionDataBlock{
			{
				Id: 1,
				DataValues: []*DataValuesRow{
					{Values: []float32{4340.6}},
				},
			},
		},
		Connections: []*EventBlock{
			{Start: 1000, End: 1500, Value: 4340.6},
		},
	}

	got, err := json.Marshal(timeline)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	gotStr := string(got)

	// Должно содержать 4340.6, а не 4340.60009765625
	if !strings.Contains(gotStr, "4340.6") {
		t.Errorf("expected 4340.6 in result, got: %s", gotStr)
	}

	// Не должно содержать 4340.60009765625
	if strings.Contains(gotStr, "4340.60009765625") {
		t.Errorf("unexpected precision loss: got %s", gotStr)
	}
}

func TestTimeline_MarshalJSON_EventBlockValue(t *testing.T) {
	// Проверяем, что EventBlock с value=1 выводится как 1, а не 1.0
	timeline := &Timeline{
		Sn: "ABC123",
		Connections: []*EventBlock{
			{Start: 1000, End: 2000, Value: 1},
		},
	}

	got, err := json.Marshal(timeline)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	want := `{"sn":"ABC123","connections":[{"start":1000,"end":2000,"value":1}]}`
	if string(got) != want {
		t.Errorf("json.Marshal() = %v, want %v", string(got), want)
	}
}

func TestTimeline_MarshalJSON_EventBlockWithID(t *testing.T) {
	// Проверяем, что EventBlock с id выводит id первым
	timeline := &Timeline{
		Sn: "ABC123",
		Connections: []*EventBlock{
			{Id: 42, Start: 1000, End: 2000, Value: 1.5},
		},
	}

	got, err := json.Marshal(timeline)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	want := `{"sn":"ABC123","connections":[{"id":42,"start":1000,"end":2000,"value":1.5}]}`
	if string(got) != want {
		t.Errorf("json.Marshal() = %v, want %v", string(got), want)
	}
}

func TestTimeline_MarshalJSON_EventBlockWithoutID(t *testing.T) {
	// Проверяем, что EventBlock без id не выводит id
	timeline := &Timeline{
		Sn: "ABC123",
		Connections: []*EventBlock{
			{Id: 0, Start: 1000, End: 2000, Value: 1.5},
		},
	}

	got, err := json.Marshal(timeline)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	want := `{"sn":"ABC123","connections":[{"start":1000,"end":2000,"value":1.5}]}`
	if string(got) != want {
		t.Errorf("json.Marshal() = %v, want %v", string(got), want)
	}
}
