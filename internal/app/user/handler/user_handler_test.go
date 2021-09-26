package handler

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/storage_user"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

type StorageUserMemory struct {
	mx      sync.RWMutex
	storage map[int]storage_user.User
}

func NewStorageUserMemory() *StorageUserMemory {
	return &StorageUserMemory{
		storage: make(map[int]storage_user.User),
	}
}

func (su *StorageUserMemory) IsUserExists(user storage_user.UserInput) (int, error) {
	su.mx.RLock()
	defer su.mx.RUnlock()

	for key, val := range su.storage {
		if val.Name == user.Name && val.Password == user.Password {
			return key, nil
		}
	}
	return -1, nil
}

func (su *StorageUserMemory) AddUser(newUser storage_user.User) (int, error) {
	su.mx.Lock()
	defer su.mx.Unlock()

	for _, val := range su.storage {
		if val == newUser {
			return -1, nil
		}
	}
	newId := len(su.storage) + 1
	su.storage[newId] = newUser
	return newId, nil
}

func TestCreateUserSuccessUnit(t *testing.T) {
	var mockStorage = NewStorageUserMemory()

	type args struct {
		str string
	}

	tests := []struct {
		name string
		args args
		statusCode int
	} {
		{
			"signup",
			args{`{"username":"Misha","email":"Misha@gmail.com","password":"1234"}` + "\n"}, http.StatusOK},
		{
			"signup",
			args{`{"username":"Glasha","email":"Glasha@gmail.com","password":"Glasha1234"}` + "\n"}, http.StatusOK},
		{
			"signup",
			args{`{"username":"Vova","email":"Putin@gmail.com","password":"Putin228"}` + "\n"}, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(tt.args.str))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := &UserHandler{mockStorage}

			if assert.NoError(t, h.SignUp(c)) {
				assert.Equal(t, tt.statusCode, rec.Code)
				assert.Equal(t, tt.args.str, rec.Body.String())
			}
		})
	}
}


func TestCreateUserFailUnit(t *testing.T) {
	var mockStorage = NewStorageUserMemory()
	mockStorage.AddUser(storage_user.User{Name: "Misha", Email: "qwerty@gmail.com", Password: "1234"})

	type args struct {
		str string
	}

	tests := []struct {
		name string
		args args
		wantedJson string
		statusCode int
	} {
		{
			"signup",
			args{`{"username":"","email":"Misha@gmail.com","password":"1234"}` + "\n"},
			"",
			http.StatusUnauthorized},
		{
			"signup",
			args{`{"username":"Glasha","email":"","password":"Glasha1234"}` + "\n"},
			"",
			http.StatusUnauthorized},
		{
			"signup",
			args{`{"username":"Vova","email":"Putin@gmail.com","password":""}` + "\n"},
			"",
			http.StatusUnauthorized},
		{
			"signup",
			args{`{"username":"Misha","email":"qwerty@gmail.com","password":"1234"}` + "\n"},
			`{"username":"Misha","email":"qwerty@gmail.com","password":"1234"}` + "\n",
			http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(tt.args.str))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := &UserHandler{mockStorage}

			if assert.NoError(t, h.SignUp(c)) {
				assert.Equal(t, tt.statusCode, rec.Code)
				assert.Equal(t, tt.wantedJson, rec.Body.String())
			}
		})
	}
}


func TestLoginUserSuccessUnit(t *testing.T) {
	var mockStorage = NewStorageUserMemory()
	mockStorage.AddUser(storage_user.User{Name: "Misha", Email: "qwerty@gmail.com", Password: "1234"})
	mockStorage.AddUser(storage_user.User{Name: "Glasha", Email: "qwerty@gmail.com", Password: "Glasha123"})
	mockStorage.AddUser(storage_user.User{Name: "Vova", Email: "qwerty@gmail.com", Password: "Putin228"})

	type args struct {
		str string
	}

	tests := []struct {
		name string
		args args
		statusCode int
	} {
		{
			"auth",
			args{`{"username":"Misha","password":"1234"}` + "\n"}, http.StatusOK},
		{
			"auth",
			args{`{"username":"Glasha","password":"Glasha123"}` + "\n"}, http.StatusOK},
		{
			"auth",
			args{`{"username":"Vova","password":"Putin228"}` + "\n"}, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.args.str))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := &UserHandler{mockStorage}

			if assert.NoError(t, h.Login(c)) {
				assert.Equal(t, tt.statusCode, rec.Code)
				assert.Equal(t, tt.args.str, rec.Body.String())
			}
		})
	}
}



func TestLoginUserFailUnit(t *testing.T) {
	var mockStorage = NewStorageUserMemory()

	mockStorage.AddUser(storage_user.User{Name: "Misha", Email: "qwerty@gmail.com", Password: "1234"})
	mockStorage.AddUser(storage_user.User{Name: "Glasha", Email: "qwerty@gmail.com", Password: "Glasha123"})
	mockStorage.AddUser(storage_user.User{Name: "Vova", Email: "qwerty@gmail.com", Password: "Putin"})

	type args struct {
		str string
	}

	tests := []struct {
		name string
		args args
		statusCode int
	} {
		{
			"auth",
			args{`{"username":"Misha","password":"134"}` + "\n"}, http.StatusBadRequest},
		{
			"auth",
			args{`{"username":"MishaX","password":"1234"}` + "\n"}, http.StatusBadRequest},
		{
			"auth",
			args{`{"username":"","password":"1234"}` + "\n"}, http.StatusBadRequest},
		{
			"auth",
			args{`{"username":"","password":""}` + "\n"}, http.StatusBadRequest},
		{
			"auth",
			args{`{"":" ","":""}` + "\n"}, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.args.str))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := &UserHandler{mockStorage}

			if assert.NoError(t, h.Login(c)) {
				assert.Equal(t, tt.statusCode, rec.Code)
				assert.NotEqual(t, tt.args.str, rec.Body.String())
			}
		})
	}
}
