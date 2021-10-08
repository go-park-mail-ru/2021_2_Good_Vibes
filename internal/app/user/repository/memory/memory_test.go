package memory

//
//import (
//	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
//	models "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
//	"golang.org/x/crypto/bcrypt"
//	"sync"
//	"testing"
//)
//
//func storageInit() map[string]models.UserDataStorage {
//	storage := make(map[string]models.UserDataStorage)
//	password1, _ := bcrypt.GenerateFromPassword([]byte("Misha_1234"), bcrypt.DefaultCost)
//	password2, _ := bcrypt.GenerateFromPassword([]byte("Sasha_1234"), bcrypt.DefaultCost)
//	password3, _ := bcrypt.GenerateFromPassword([]byte("Gosha_1234"), bcrypt.DefaultCost)
//
//	user1 := models.UserDataStorage{1, "Misha", "Misha@mail.ru", string(password1)}
//	user2 := models.UserDataStorage{2, "Sasha", "Sasha@mail.ru", string(password2)}
//	user3 := models.UserDataStorage{3, "Gosha", "Gosha@mail.ru", string(password3)}
//
//	storage[user1.Name] = user1
//	storage[user2.Name] = user2
//	storage[user3.Name] = user3
//
//	return storage
//}
//
//func TestStorageUserMemory_IsUserExistsSuccess(t *testing.T) {
//	storage := storageInit()
//
//	UserDataForInput1 := models.UserDataForInput{"Misha", "Misha_1234"}
//
//	type args struct {
//		user models.UserDataForInput
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    int
//		wantErr bool
//	}{
//		{"success",
//			args{UserDataForInput1},
//			1,
//			false},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			su := &StorageUserMemory{
//				mx:     sync.RWMutex{},
//				storage: storage,
//			}
//			got, err := su.IsUserExists(tt.args.user)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("IsUserExists() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("IsUserExists() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestStorageUserMemory_IsUserExistsFail(t *testing.T) {
//	storage := storageInit()
//
//	UserDataForInput1 := models.UserDataForInput{"Lena", "Lena_1234"}
//	UserDataForInput2 := models.UserDataForInput{"Sasha", "Sasha_123"}
//
//	type args struct {
//		user models.UserDataForInput
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    int
//		wantErr bool
//	}{
//		{"fail",
//			args{UserDataForInput1},
//			errors.NO_USER_ERROR,
//			false},
//		{"fail",
//			args{UserDataForInput2},
//			errors.WRONG_PASSWORD_ERROR,
//			false},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			su := &StorageUserMemory{
//				mx:     sync.RWMutex{},
//				storage: storage,
//			}
//			got, err := su.IsUserExists(tt.args.user)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("IsUserExists() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("IsUserExists() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestStorageUserMemory_AddUserSuccess(t *testing.T) {
//	storage := storageInit()
//	UserDataForInput1 := models.UserDataForReg{"Alla", "Alla@gmail.com", "Alla_1234"}
//
//	type args struct {
//		user models.UserDataForReg
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    int
//		wantErr bool
//	}{
//		{"success",
//			args{UserDataForInput1},
//			4,
//			false},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			su := &StorageUserMemory{
//				mx:     sync.RWMutex{},
//				storage: storage,
//			}
//			_, err := su.AddUser(tt.args.user)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("IsUserExists() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//		})
//	}
//}
//
//func TestStorageUserMemory_AddUserFail(t *testing.T) {
//	storage := storageInit()
//
//	UserDataForInput1 := models.UserDataForReg{"Misha", "Misha@gmail.com", "Misha_1234"}
//
//	type args struct {
//		user models.UserDataForReg
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    int
//		wantErr bool
//	}{
//		{"fail",
//			args{UserDataForInput1},
//			errors.USER_EXISTS_ERROR,
//			false},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			su := &StorageUserMemory{
//				mx:     sync.RWMutex{},
//				storage: storage,
//			}
//			got, err := su.AddUser(tt.args.user)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("IsUserExists() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("IsUserExists() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
