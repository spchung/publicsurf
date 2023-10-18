package service

import (
	"log"
	"public-surf/internal/domain/entity"
	"public-surf/internal/domain/repository"
	"public-surf/pkg/config"
	"public-surf/pkg/database"
	"reflect"
	"testing"
)

func TestUserService_GetUser(t *testing.T) {
	config.GetConfig()
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	type fields struct {
		userRepo  repository.IUserRepository
		photoRepo repository.IPhotoRepository
	}
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.UserViewModel
		wantErr bool
	}{
		{
			name: "Test 1",
			fields: fields{
				userRepo:  repository.NewUserRepository(db),
				photoRepo: repository.NewPhotoRepository(db),
			},
			args: args{
				id: 1,
			},
			want: entity.UserViewModel{
				ID:        1,
				Email:     "chungstephen96@gmail,com",
				FirstName: "Stephen",
				LastName:  "Chung",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserService{
				userRepo:  tt.fields.userRepo,
				photoRepo: tt.fields.photoRepo,
			}
			got, err := s.GetUser(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
