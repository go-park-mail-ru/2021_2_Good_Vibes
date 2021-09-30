package impl

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	userModel "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"testing"
)

func storageInit() map[string]userModel.UserStorage {
	storage := make(map[string]userModel.UserStorage)
	password1, _ := bcrypt.GenerateFromPassword([]byte("Misha_1234"), bcrypt.DefaultCost)
	password2, _ := bcrypt.GenerateFromPassword([]byte("Sasha_1234"), bcrypt.DefaultCost)
	password3, _ := bcrypt.GenerateFromPassword([]byte("Gosha_1234"), bcrypt.DefaultCost)

	user1 := userModel.NewUserStorage(1, "Misha", "Misha@mail.ru", string(password1))
	user2 := userModel.NewUserStorage(2, "Sasha", "Sasha@mail.ru", string(password2))
	user3 := userModel.NewUserStorage(3, "Gosha", "Gosha@mail.ru", string(password3))

	storage[user1.Name] = user1
	storage[user2.Name] = user2
	storage[user3.Name] = user3

	return storage
}

func TestStorageUserMemory_IsUserExistsSuccess(t *testing.T) {
	storage := storageInit()

	userInput1 := userModel.UserInput{"Misha", "Misha_1234"}

	type args struct {
		user userModel.UserInput
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"success",
			args{userInput1},
			1,
			false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			su := &StorageUserMemory{
				mx:     sync.RWMutex{},
				storage: storage,
			}
			got, err := su.IsUserExists(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsUserExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsUserExists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorageUserMemory_IsUserExistsFail(t *testing.T) {
	storage := storageInit()

	userInput1 := userModel.UserInput{"Lena", "Lena_1234"}
	userInput2 := userModel.UserInput{"Sasha", "Sasha_123"}

	type args struct {
		user userModel.UserInput
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"fail",
			args{userInput1},
			errors.NO_USER_ERROR,
			false},
		{"fail",
			args{userInput2},
			errors.WRONG_PASSWORD_ERROR,
			false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			su := &StorageUserMemory{
				mx:     sync.RWMutex{},
				storage: storage,
			}
			got, err := su.IsUserExists(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsUserExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsUserExists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorageUserMemory_AddUserSuccess(t *testing.T) {
	storage := storageInit()
	userInput1 := userModel.User{"Alla", "Alla@gmail.com", "Alla_1234"}

	type args struct {
		user userModel.User
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"success",
			args{userInput1},
			4,
			false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			su := &StorageUserMemory{
				mx:     sync.RWMutex{},
				storage: storage,
			}
			_, err := su.AddUser(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsUserExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStorageUserMemory_AddUserFail(t *testing.T) {
	storage := storageInit()

	userInput1 := userModel.User{"Misha", "Misha@gmail.com", "Misha_1234"}

	type args struct {
		user userModel.User
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"fail",
			args{userInput1},
			errors.USER_EXISTS_ERROR,
			false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			su := &StorageUserMemory{
				mx:     sync.RWMutex{},
				storage: storage,
			}
			got, err := su.AddUser(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsUserExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsUserExists() got = %v, want %v", got, tt.want)
			}
		})
	}
}
