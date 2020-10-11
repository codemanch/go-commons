package logging

import (
	"reflect"
	"testing"
)

// TestGetLogger --> Testing Logger object creation
func TestGetLogger(t *testing.T) {
	tests := []struct {
		name string
		want *Logger
	}{
		{
			name: "Logger",
			want: GetLogger(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLogger(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogger_IsEnabled(t *testing.T) {
	type fields struct {
		sev             Severity
		pkgName         string
		errorEnabled    bool
		warnEnabled     bool
		infoEnabled     bool
		debugEnabled    bool
		traceEnabled    bool
		includeFunction bool
		includeLine     bool
	}
	type args struct {
		sev Severity
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{
			name: "WarnTest_true",
			fields: fields{
				sev : WarnLvl,
				warnEnabled: true,
			},
			args: args{
				sev: WarnLvl,
			},
			want: true,
		},
		{
			name: "WarnTest_Fail",
			fields: fields{
				sev : InfoLvl,
				infoEnabled: true,
			},
			args: args{
				sev: WarnLvl,
			},
			want: false,
		},
		{
			name: "ErrorTest",
			fields: fields{
				sev : ErrLvl,
				errorEnabled: true,
			},
			args: args{
				sev: ErrLvl,
			},
			want: true,
		},
		{
			name: "ErrorTest_Fail",
			fields: fields{
				sev : TraceLvl,
				traceEnabled: true,
			},
			args: args{
				sev: ErrLvl,
			},
			want: false,
		},
		{
			name: "InfoTest",
			fields: fields{
				sev : InfoLvl,
				infoEnabled: true,
			},
			args: args{
				sev: InfoLvl,
			},
			want: true,
		},
		{
			name: "InfoTest_Fail",
			fields: fields{
				sev : TraceLvl,
				traceEnabled: true,
			},
			args: args{
				sev: InfoLvl,
			},
			want: false,
		},
		{
			name: "DebugTest",
			fields: fields{
				sev : DebugLvl,
				warnEnabled: true,
			},
			args: args{
				sev: DebugLvl,
			},
			want: true,
		},
		{
			name: "DebugTest_Fail",
			fields: fields{
				sev : TraceLvl,
				traceEnabled: true,
			},
			args: args{
				sev: DebugLvl,
			},
			want: false,
		},
		{
			name: "TraceTest",
			fields: fields{
				sev : TraceLvl,
				warnEnabled: true,
			},
			args: args{
				sev: TraceLvl,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				sev:             tt.fields.sev,
				pkgName:         tt.fields.pkgName,
				errorEnabled:    tt.fields.errorEnabled,
				warnEnabled:     tt.fields.warnEnabled,
				infoEnabled:     tt.fields.infoEnabled,
				debugEnabled:    tt.fields.debugEnabled,
				traceEnabled:    tt.fields.traceEnabled,
				includeFunction: tt.fields.includeFunction,
				includeLine:     tt.fields.includeLine,
			}
			if got := l.IsEnabled(tt.args.sev); got != tt.want {
				t.Errorf("IsEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}