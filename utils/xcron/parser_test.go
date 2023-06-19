package xcron

import (
	"reflect"
	"testing"
	"time"
)

func TestParser_Parse(t *testing.T) {
	type fields struct {
		options ParseOption
	}
	type args struct {
		spec string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Schedule
		wantErr bool
	}{
		{args: args{spec: "@daily 3h"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(Second | Minute | Hour | Dom | Month | DowOptional | Descriptor)
			got, err := p.Parse(tt.args.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(time.Now())
			t.Log(got.Next(time.Now()))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
