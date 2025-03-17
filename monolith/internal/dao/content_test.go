package dao

import (
	"imooc-content-system/internal/model"

	"gorm.io/gorm"

	"testing"

	"gorm.io/driver/mysql"
)

func connDB() *gorm.DB {
	mysqlDB, err := gorm.Open(mysql.Open("root:root123@tcp(localhost:3306)/cms_content?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	db, err := mysqlDB.DB()
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(2)
	return mysqlDB
}

func TestContentDao_Create(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		detail model.ContentDetail
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test create content",
			fields: fields{
				db: connDB(),
			},
			args: args{
				detail: model.ContentDetail{
					ID:          1,
					Title:       "test title",
					Description: "test description",
					Author:      "test author",
					VideoURL:    "test video url",
					Thumbnail:   "test thumbnail",
					Category:    "test category",
				},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ContentDao{
				db: tt.fields.db,
			}
			got, err := c.Create(tt.args.detail)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}
