package main

import (
	"reflect"
	"testing"
)

func TestParseResult(t *testing.T) {
	type args struct {
		result string
	}
	tests := []struct {
		name       string
		args       args
		wantRecord MonitorRecord
		wantErr    bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRecord, err := ParseResult(tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRecord, tt.wantRecord) {
				t.Errorf("ParseResult() = %v, want %v", gotRecord, tt.wantRecord)
			}
		})
	}
}

func TestWriteRecordToInfluxDB(t *testing.T) {
	type args struct {
		address     string
		db          string
		measurement string
		record      MonitorRecord
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteRecordToInfluxDB(tt.args.address, tt.args.db, tt.args.measurement, tt.args.record)
		})
	}
}
