package handler

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	user_model "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/storage_user/impl"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//type StorageUserMemory struct {
//	mx      sync.RWMutex
//	storage map[int]user_model.User
//}
//
//func NewStorageUserMemory() *StorageUserMemory {
//	return &StorageUserMemory{
//		storage: make(map[int]user_model.User),
//	}
//}
//
//func (su *StorageUserMemory) IsUserExists(user user_model.UserInput) (int, error) {
//	su.mx.RLock()
//	defer su.mx.RUnlock()
//
//	for key, val := range su.storage {
//		if val.Name == user.Name && val.Password == user.Password {
//			return key, nil
//		}
//	}
//	return -1, nil
//}
//
//func (su *StorageUserMemory) AddUser(newUser user_model.User) (int, error) {
//	su.mx.Lock()
//	defer su.mx.Unlock()
//
//	for _, val := range su.storage {
//		if val == newUser {
//			return -1, nil
//		}
//	}
//	newId := len(su.storage) + 1
//	su.storage[newId] = newUser
//	return newId, nil
//}
//
//func NewUser (name string, email string, password string) user_model.User {
//	return user_model.User{
//		Name: name,
//		Email: email,
//		Password: password,
//	}
//}
//
//func NewUserInput (name string, password string) user_model.UserInput {
//	return user_model.UserInput{
//		Name: name,
//		Password: password,
//	}
//}

func TestCreateUserSuccessUnit(t *testing.T) {
	var mockStorage, _ = impl.NewStorageUserMemory()

	user1, _ := json.Marshal(user_model.User{"Misha", "Misha@gmail.com", "Misha_1234"})
	user2, _ := json.Marshal(user_model.User{"Glasha", "Glasha@gmail.com", "Glasha_1234"})
	user3, _ := json.Marshal(user_model.User{"Vova", "Vova@gmail.com", "Vova_1234"})

	type args struct {
		str string
	}

	tests := []struct {
		name       string
		args       args
		statusCode int
	}{
		{
			"signup",
			args{string(user1) + "\n"},
			http.StatusOK},
		{
			"signup",
			args{string(user2) + "\n"},
			http.StatusOK},
		{
			"signup",
			args{string(user3) + "\n"},
			http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := echo.New()

			val := validator.New()
			val.RegisterValidation("customPassword", Password)
			router.Validator = &CustomValidator{Validator: val}

			rec, ctx, h := constructRequest("/signup", tt.args.str, router, mockStorage)

			if assert.NoError(t, h.SignUp(ctx)) {
				assert.Equal(t, tt.statusCode, rec.Code)
				assert.Equal(t, tt.args.str, rec.Body.String())
			}
		})
	}
}

func TestCreateUserFailUnit(t *testing.T) {
	var mockStorage, _ = impl.NewStorageUserMemory()
	mockStorage.AddUser(user_model.User{Name: "Misha", Email: "qwerty@gmail.com", Password: "Misha_1234"})

	user1, _ := json.Marshal(user_model.User{"", "Misha@gmail.com", "Misha_1234"})
	user2, _ := json.Marshal(user_model.User{"Glasha","", "Glasha_1234"})
	user3, _ := json.Marshal(user_model.User{"Vova", "Putin@gmail.com", ""})
	user4, _ := json.Marshal(user_model.User{"Vova", "Putin@gmail.com", "Vova1234"})
	user5, _ := json.Marshal(user_model.User{"Vova", "Putin@gmail.com", "Vova_"})
	user6, _ := json.Marshal(user_model.User{"Vova", "Putin@gmail.com", "1234"})
	user7, _ := json.Marshal(user_model.User{"Misha", "qwerty@gmail.com", "Misha_1234"})

	wantedUserResp1, _ := json.Marshal(user_model.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR))
	wantedUserResp2, _ := json.Marshal(user_model.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR))
	wantedUserResp3, _ := json.Marshal(user_model.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR))
	wantedUserResp4, _ := json.Marshal(user_model.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR))
	wantedUserResp5, _ := json.Marshal(user_model.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR))
	wantedUserResp6, _ := json.Marshal(user_model.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR))
	wantedUserResp7, _ := json.Marshal(user_model.NewError(errors.USER_EXISTS_ERROR, errors.USER_EXISTS_DESCR))

	type args struct {
		str string
	}

	tests := []struct {
		name       string
		args       args
		wantedJson string
		statusCode int
	}{
		{
			"signup",
			args{string(user1) + "\n"},
			string(wantedUserResp1)+ "\n",
			http.StatusBadRequest},
		{
			"signup",
			args{string(user2) + "\n"},
			string(wantedUserResp2)+ "\n",
			http.StatusBadRequest},
		{
			"signup",
			args{string(user3) + "\n"},
			string(wantedUserResp3)+ "\n",
			http.StatusBadRequest},
		{
			"signup",
			args{string(user4) + "\n"},
			string(wantedUserResp4)+ "\n",
			http.StatusBadRequest},
		{
			"signup",
			args{string(user5) + "\n"},
			string(wantedUserResp5)+ "\n",
			http.StatusBadRequest},
		{
			"signup",
			args{string(user6) + "\n"},
			string(wantedUserResp6)+ "\n",
			http.StatusBadRequest},
		{
			"signup",
			args{string(user7) + "\n"},
			string(wantedUserResp7)+ "\n",
			http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := echo.New()

			val := validator.New()
			val.RegisterValidation("customPassword", Password)
			router.Validator = &CustomValidator{Validator: val}

			rec, ctx, h := constructRequest("/signup", tt.args.str, router, mockStorage)

			if assert.NoError(t, h.SignUp(ctx)) {
				assert.Equal(t, tt.statusCode, rec.Code)
				assert.Equal(t, tt.wantedJson, rec.Body.String())
			}
		})
	}
}

func TestLoginUserSuccessUnit(t *testing.T) {
	var mockStorage, _ = impl.NewStorageUserMemory()
	mockStorage.AddUser(user_model.User{Name: "Misha", Email: "qwerty@gmail.com", Password: "1234"})
	mockStorage.AddUser(user_model.User{Name: "Glasha", Email: "qwerty@gmail.com", Password: "Glasha123"})
	mockStorage.AddUser(user_model.User{Name: "Vova", Email: "qwerty@gmail.com", Password: "Putin228"})

	user1, _ := json.Marshal(user_model.UserInput{"Misha", "1234"})
	user2, _ := json.Marshal(user_model.UserInput{"Glasha", "Glasha123"})
	user3, _ := json.Marshal(user_model.UserInput{"Vova", "Putin228"})

	type args struct {
		str string
	}

	tests := []struct {
		name       string
		args       args
		statusCode int
	}{
		{
			"auth",
			args{string(user1) + "\n"}, http.StatusOK},
		{
			"auth",
			args{string(user2) + "\n"}, http.StatusOK},
		{
			"auth",
			args{string(user3) + "\n"}, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := echo.New()
			val := validator.New()
			val.RegisterValidation("customPassword", Password)
			router.Validator = &CustomValidator{Validator: val}

			rec, ctx, h := constructRequest("/login", tt.args.str, router, mockStorage)

			if assert.NoError(t, h.Login(ctx)) {
				assert.Equal(t, tt.statusCode, rec.Code)
				assert.Equal(t, tt.args.str, rec.Body.String())
			}
		})
	}
}

func TestLoginUserFailUnit(t *testing.T) {
	var mockStorage, _ = impl.NewStorageUserMemory()

	mockStorage.AddUser(user_model.User{Name: "Misha", Email: "qwerty@gmail.com", Password: "Misha_1234"})

	user1, _ := json.Marshal(user_model.UserInput{"Misha", "Misha_134"})
	user2, _ := json.Marshal(user_model.UserInput{"MishaX", "1234"})
	user3, _ := json.Marshal(user_model.UserInput{"", "Misha_1234"})
	user4, _ := json.Marshal(user_model.UserInput{"Misha", ""})
	user5, _ := json.Marshal(user_model.UserInput{"", ""})

	wantedUserResp1, _ := json.Marshal(user_model.NewError(errors.WRONG_PASSWORD_ERROR, errors.WRONG_PASSWORD_DESCR))
	wantedUserResp2, _ := json.Marshal(user_model.NewError(errors.NO_USER_ERROR, errors.NO_USER_DESCR))
	wantedUserResp3, _ := json.Marshal(user_model.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR))
	wantedUserResp4, _ := json.Marshal(user_model.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR))
	wantedUserResp5, _ := json.Marshal(user_model.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR))

	type args struct {
		str string
	}

	tests := []struct {
		name       string
		args       args
		wantedJson string
		statusCode int
	}{
		{
			"auth",
			args{string(user1)+ "\n"},
			string(wantedUserResp1) + "\n",
			http.StatusUnauthorized},
		{
			"auth",
			args{string(user2)+ "\n"},
			string(wantedUserResp2) + "\n",
			http.StatusUnauthorized},
		{
			"auth",
			args{string(user3)+ "\n"},
			string(wantedUserResp3) + "\n",
			http.StatusBadRequest},
		{
			"auth",
			args{string(user4)+ "\n"},
			string(wantedUserResp4) + "\n",
			http.StatusBadRequest},
		{
			"auth",
			args{string(user5)+ "\n"},
			string(wantedUserResp5) + "\n",
			http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := echo.New()
			val := validator.New()
			router.Validator = &CustomValidator{Validator: val}

			rec, ctx, h := constructRequest("/login", tt.args.str, router, mockStorage)

			if assert.NoError(t, h.Login(ctx)) {
				assert.Equal(t, tt.statusCode, rec.Code)
				assert.Equal(t, tt.wantedJson, rec.Body.String())
			}
		})
	}
}

func TestCreateUserLoginIntegrationSuccess(t *testing.T) {
	var mockStorage, _ = impl.NewStorageUserMemory()

	userSignUp, _ := json.Marshal(user_model.User{"Misha", "Misha@gmail.com", "Misha_1234"})
	userLogin, _ := json.Marshal(user_model.UserInput{"Misha", "Misha_1234"})

	wantedUserSignUpResp, _ := json.Marshal(user_model.User{"Misha", "Misha@gmail.com", "Misha_1234"})
	wantedUserLoginResp, _ := json.Marshal(user_model.UserInput{"Misha", "Misha_1234"})

	type args struct {
		signUp string
		login  string
	}

	tests := []struct {
		name             string
		args             args
		wantedSignupJson string
		wantedLoginJson  string
		signUpStatusCode int
		loginStatusCode  int
	}{
		{
			"signup_login_integration",
			args{string(userSignUp) + "\n",
				string(userLogin) + "\n"},
			string(wantedUserSignUpResp) + "\n",
			string(wantedUserLoginResp) + "\n",
			http.StatusOK,
			http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := echo.New()

			val := validator.New()
			val.RegisterValidation("customPassword", Password)
			router.Validator = &CustomValidator{Validator: val}

			rec, ctx, h := constructRequest("/signup", tt.args.signUp, router, mockStorage)

			if assert.NoError(t, h.SignUp(ctx)) {
				assert.Equal(t, tt.signUpStatusCode, rec.Code)
				assert.Equal(t, tt.wantedSignupJson, rec.Body.String())
			}

			rec, ctx, h = constructRequest("/login", tt.args.login, router, mockStorage)

			if assert.NoError(t, h.Login(ctx)) {
				assert.Equal(t, tt.loginStatusCode, rec.Code)
				assert.Equal(t, tt.wantedLoginJson, rec.Body.String())
			}
		})
	}
}

func TestCreateUserLoginIntegrationFail(t *testing.T) {
	var mockStorage, _ = impl.NewStorageUserMemory()

	userSignUp1, _ := json.Marshal(user_model.User{"Misha", "Misha@gmail.com", "Misha_1234"})
	userLogin1, _ := json.Marshal(user_model.UserInput{"Gosha", "Misha_1234"})

	wantedUserSignUpResp1, _ := json.Marshal(user_model.User{"Misha", "Misha@gmail.com", "Misha_1234"})
	wantedUserLoginResp1, _ := json.Marshal(user_model.NewError(errors.NO_USER_ERROR, errors.NO_USER_DESCR))

	userSignUp2, _ := json.Marshal(user_model.User{"Gosha", ":", "1234"})
	userLogin2, _ := json.Marshal(user_model.UserInput{"Gosha", "1234"})

	wantedUserSignUpResp2, _ := json.Marshal(user_model.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR))
	wantedUserLoginResp2, _ := json.Marshal(user_model.NewError(errors.NO_USER_ERROR, errors.NO_USER_DESCR))

	userSignUp3, _ := json.Marshal(user_model.User{"Sasha", "Sasha@mail.ru", "Sasha_1234"})
	userLogin3, _ := json.Marshal(user_model.UserInput{"Sasha", "Sasha_234"})

	wantedUserSignUpResp3, _ := json.Marshal(user_model.User{"Sasha", "Sasha@mail.ru", "Sasha_1234"})
	wantedUserLoginResp3, _ := json.Marshal(user_model.NewError(errors.WRONG_PASSWORD_ERROR, errors.WRONG_PASSWORD_DESCR))

	type args struct {
		signUp string
		login  string
	}

	tests := []struct {
		name             string
		args             args
		wantedSignupJson string
		wantedLoginJson  string
		signUpStatusCode int
		loginStatusCode  int
	}{
		{
			"signup_login_integration",
			args{string(userSignUp1) + "\n",
				string(userLogin1) + "\n"},
			string(wantedUserSignUpResp1) + "\n",
			string(wantedUserLoginResp1) + "\n",
			http.StatusOK,
			http.StatusUnauthorized},
		{
			"signup_login_integration",
			args{string(userSignUp2) + "\n",
				string(userLogin2) + "\n"},
			string(wantedUserSignUpResp2) + "\n",
			string(wantedUserLoginResp2) + "\n",
			http.StatusBadRequest,
			http.StatusUnauthorized},
		{
			"signup_login_integration",
			args{string(userSignUp3) + "\n",
				string(userLogin3) + "\n"},
			string(wantedUserSignUpResp3) + "\n",
			string(wantedUserLoginResp3) + "\n",
			http.StatusOK,
			http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := echo.New()
			val := validator.New()
			val.RegisterValidation("customPassword", Password)
			router.Validator = &CustomValidator{Validator: val}

			rec, ctx, h := constructRequest("/signup", tt.args.signUp, router, mockStorage)

			if assert.NoError(t, h.SignUp(ctx)) {
				assert.Equal(t, tt.signUpStatusCode, rec.Code)
				assert.Equal(t, tt.wantedSignupJson, rec.Body.String())
			}

			rec, ctx, h = constructRequest("/login", tt.args.login, router, mockStorage)

			if assert.NoError(t, h.Login(ctx)) {
				assert.Equal(t, tt.loginStatusCode, rec.Code)
				assert.Equal(t, tt.wantedLoginJson, rec.Body.String())
			}
		})
	}
}

func constructRequest(target string, login string, router *echo.Echo, mockStorage *impl.StorageUserMemory) (*httptest.ResponseRecorder, echo.Context, *UserHandler) {
	req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(login))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := router.NewContext(req, rec)
	h := &UserHandler{mockStorage}
	return rec, ctx, h
}