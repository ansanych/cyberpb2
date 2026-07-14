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
			name: "only params",
			timeline: &Timeline{
				Params: []string{"rpm", "temp"},
			},
			want: `{"params":["rpm","temp"]}`,
		},
		{
			name: "single connection without dataValues",
			timeline: &Timeline{
				Data: []*TelemetryConnection{
					{
						Id:    1,
						Start: 1000,
						End:   2000,
					},
				},
			},
			want: `{"data":[{"id":1,"start":1000,"end":2000}]}`,
		},
		{
			name: "connection with finished=true",
			timeline: &Timeline{
				Data: []*TelemetryConnection{
					{
						Id:       1,
						Start:    1000,
						End:      2000,
						Finished: true,
					},
				},
			},
			want: `{"data":[{"id":1,"start":1000,"end":2000,"finished":true}]}`,
		},
		{
			name: "connection with dataValues as array of arrays",
			timeline: &Timeline{
				Data: []*TelemetryConnection{
					{
						Id:    1,
						Start: 1000,
						End:   2000,
						DataValues: []*DataValuesRow{
							{Values: []float32{1.0, 2.0, 3.0}},
							{Values: []float32{4.0, 5.0, 6.0}},
						},
					},
				},
			},
			want: `{"data":[{"id":1,"start":1000,"end":2000,"dataValues":[[1,2,3],[4,5,6]]}]}`,
		},
		{
			name: "connection with engine blocks",
			timeline: &Timeline{
				Data: []*TelemetryConnection{
					{
						Id:    1,
						Start: 1000,
						End:   2000,
						Engine: []*TelemetryBlock{
							{Start: 1000, End: 1500, Value: 1.5},
							{Start: 1500, End: 2000, Value: 2.5},
						},
					},
				},
			},
			want: `{"data":[{"id":1,"start":1000,"end":2000,"engine":[{"start":1000,"end":1500,"value":1.5},{"start":1500,"end":2000,"value":2.5}]}]}`,
		},
		{
			name: "full timeline with params and data",
			timeline: &Timeline{
				Params: []string{"rpm", "temp"},
				Data: []*TelemetryConnection{
					{
						Id:       1,
						Start:    1000,
						End:      2000,
						Finished: true,
						Engine: []*TelemetryBlock{
							{Start: 1000, End: 1500, Value: 1.5},
						},
						DataValues: []*DataValuesRow{
							{Values: []float32{1.0, 2.0}},
							{Values: []float32{3.0, 4.0}},
						},
					},
				},
			},
			want: `{"params":["rpm","temp"],"data":[{"id":1,"start":1000,"end":2000,"finished":true,"engine":[{"start":1000,"end":1500,"value":1.5}],"dataValues":[[1,2],[3,4]]}]}`,
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
		Data: []*TelemetryConnection{
			{
				Id:     1,
				Start:  1000,
				End:    2000,
				Engine: []*TelemetryBlock{},
				Job:    nil,
				Errors: nil,
			},
		},
	}

	got, err := json.Marshal(timeline)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	want := `{"data":[{"id":1,"start":1000,"end":2000}]}`
	if string(got) != want {
		t.Errorf("json.Marshal() = %v, want %v", string(got), want)
	}
}

func TestTimeline_MarshalJSON_DataValuesEmpty(t *testing.T) {
	// Проверяем, что пустой dataValues не выводится
	timeline := &Timeline{
		Data: []*TelemetryConnection{
			{
				Id:         1,
				Start:      1000,
				End:        2000,
				DataValues: []*DataValuesRow{},
			},
		},
	}

	got, err := json.Marshal(timeline)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	want := `{"data":[{"id":1,"start":1000,"end":2000}]}`
	if string(got) != want {
		t.Errorf("json.Marshal() = %v, want %v", string(got), want)
	}
}

func TestTimeline_MarshalJSON_FloatPrecision(t *testing.T) {
	// Проверяем, что float32 4340.6 не превращается в 4340.60009765625
	timeline := &Timeline{
		Data: []*TelemetryConnection{
			{
				Id:    1,
				Start: 1000,
				End:   2000,
				DataValues: []*DataValuesRow{
					{Values: []float32{4340.6}},
				},
				Engine: []*TelemetryBlock{
					{Start: 1000, End: 1500, Value: 4340.6},
				},
			},
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

func TestTimeline_MarshalJSON_DataValuesWithEmptyRow(t *testing.T) {
	// Проверяем, что строка с пустым Values выводится как []
	timeline := &Timeline{
		Data: []*TelemetryConnection{
			{
				Id:    1,
				Start: 1000,
				End:   2000,
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

	want := `{"data":[{"id":1,"start":1000,"end":2000,"dataValues":[[],[1]]}]}`
	if string(got) != want {
		t.Errorf("json.Marshal() = %v, want %v", string(got), want)
	}
}
