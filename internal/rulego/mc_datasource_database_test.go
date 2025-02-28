package rulego

import (
	"reflect"
	"testing"
)

func TestDataSourceDatabaseNode_processSQLAndParams(t *testing.T) {
	type fields struct {
		Config DataSourceDatabaseNodeConfiguration
	}
	type args struct {
		msgData map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   []interface{}
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				Config: DataSourceDatabaseNodeConfiguration{
					DatasourceSql: "select * from kinds limit ${msg.start}",
				},
			},
			args: args{
				msgData: map[string]interface{}{
					"start": 10,
				},
			},
			want:    "select * from kinds limit ?",
			want1:   []interface{}{10},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &DataSourceDatabaseNode{
				Config: tt.fields.Config,
			}
			got, got1, err := n.processSQLAndParams(tt.args.msgData)
			if (err != nil) != tt.wantErr {
				t.Errorf("DataSourceDatabaseNode.processSQLAndParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DataSourceDatabaseNode.processSQLAndParams() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("DataSourceDatabaseNode.processSQLAndParams() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
