package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"project/internal/config"
	"project/internal/database"
	"project/internal/handlers"
	"project/internal/models"
	"testing"

	"go.mongodb.org/mongo-driver/bson"

	// "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	cfg := config.LoadConfig()
	err := database.InitializeMongo(cfg.Mongo)
	if err != nil {
		panic("Failed to initialize test database: " + err.Error())
	}
	m.Run()
	database.CloseMongo()
}

func setupTest(t *testing.T) {
	// err := database.GetMongoDB().Collection("users").CreateIndex(context.Background(), bson.M{"email": 1}, options.Index().SetUnique(true))
	// if err != nil {
	// 	t.Fatalf("Failed to migrate tables: %v", err)
	// }
}

func teardownTest(t *testing.T) {
	err := database.GetMongoDB().Collection("users").Drop(context.Background())
	if err != nil {
		t.Fatalf("Failed to clean up database: %v", err)
	}
}

func parseErrorResponse(body []byte) (models.ErrorResponse, error) {
	var errorResp models.ErrorResponse
	err := json.Unmarshal(body, &errorResp)
	return errorResp, err
}

func parseSuccessResponse(body []byte, data interface{}) (models.SuccessResponse, error) {
	var successResp models.SuccessResponse
	err := json.Unmarshal(body, &successResp)
	if err == nil && data != nil {
		// تحويل Data إلى النوع المطلوب
		jsonData, _ := json.Marshal(successResp.Data)
		json.Unmarshal(jsonData, data)
	}
	return successResp, err
}

func TestCreateUser_InvalidData(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)
	
	// 💡 اختبار 1: بيانات ناقصة (بدون email)
	invalidUser := map[string]interface{}{
		"name":     "مستخدم بدون email",
		"password": "password123",
		// missing email intentionally
	}
	
	userJSON, _ := json.Marshal(invalidUser)
	
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	assert.NoError(t, err)
	
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	
	handler := handlers.NewUserHandler()
	router := mux.NewRouter()
	router.HandleFunc("/users", handler.CreateUser).Methods("POST")
	
	router.ServeHTTP(rr, req)
	
	// ✅ التحقق من status code
	assert.Equal(t, http.StatusBadRequest, rr.Code, "يجب أن يرجع 400 للبيانات الناقصة")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	
	// ✅ التحقق من response كـ JSON
	errorResp, err := parseErrorResponse(rr.Body.Bytes())
	assert.NoError(t, err)
	assert.False(t, errorResp.Success)
	assert.Contains(t, errorResp.Message, "Email is required", "يجب أن تحتوي الرسالة على خطأ email")
	
	// 💡 اختبار 2: email غير صالح
	invalidEmailUser := map[string]interface{}{
		"name":     "مستخدم",
		"email":    "invalid-email", // email غير صالح
		"password": "password123",   // ✅ تأكد من إضافة password
	}
	
	userJSON2, _ := json.Marshal(invalidEmailUser)
	
	req2, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON2))
	req2.Header.Set("Content-Type", "application/json")
	rr2 := httptest.NewRecorder()
	
	router.ServeHTTP(rr2, req2)
	
	// ✅ الآن يجب أن يرجع 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, rr2.Code, "يجب أن يرجع 400 لـ email غير صالح")
	assert.Equal(t, "application/json", rr2.Header().Get("Content-Type"))
	
	errorResp2, err := parseErrorResponse(rr2.Body.Bytes())
	assert.NoError(t, err)
	assert.False(t, errorResp2.Success)
	assert.Contains(t, errorResp2.Message, "Invalid email format", "يجب أن تحتوي الرسالة على خطأ تنسيق email")
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)
	
	// إنشاء مستخدم أول
	firstUser := models.User{
		Name:     "المستخدم الأول", 
		Email:    "duplicate@example.com", 
		Password: "password123", // ✅ إضافة password
	}
	_, err := database.GetMongoDB().Collection("users").InsertOne(context.Background(), firstUser)
	assert.NoError(t, err)
	
	// محاولة إنشاء مستخدم ثاني بنفس Email
	duplicateUser := map[string]interface{}{
		"name":     "المستخدم الثاني",
		"email":    "duplicate@example.com", // نفس Email
		"password": "password456",           // ✅ إضافة password
	}
	
	userJSON, _ := json.Marshal(duplicateUser)
	
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	
	handler := handlers.NewUserHandler()
	router := mux.NewRouter()
	router.HandleFunc("/users", handler.CreateUser).Methods("POST")
	
	router.ServeHTTP(rr, req)
	
	// ✅ الآن يجب أن يرجع 400 Bad Request بدلاً من 500
	assert.Equal(t, http.StatusBadRequest, rr.Code, "يجب أن يرجع 400 للـ email المكرر")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	
	errorResp, err := parseErrorResponse(rr.Body.Bytes())
	assert.NoError(t, err)
	assert.False(t, errorResp.Success)
	assert.Contains(t, errorResp.Message, "already exists", "يجب أن تحتوي الرسالة على خطأ email مكرر")
}

func TestCreateUser_ValidData(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)
	
	newUser := map[string]interface{}{
		"name":     "محمد علي",
		"email":    "mohamed@example.com", 
		"password": "newpassword123", // ✅ تأكد من إضافة password
	}
	
	userJSON, err := json.Marshal(newUser)
	assert.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	assert.NoError(t, err)
	
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	
	handler := handlers.NewUserHandler()
	router := mux.NewRouter()
	router.HandleFunc("/users", handler.CreateUser).Methods("POST")
	
	router.ServeHTTP(rr, req)
	
	// ✅ يجب أن يرجع 201 Created للبيانات الصحيحة
	assert.Equal(t, http.StatusCreated, rr.Code, "يجب أن يرجع 201 للبيانات الصحيحة")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	
	var createdUser models.UserResponse
	successResp, err := parseSuccessResponse(rr.Body.Bytes(), &createdUser)
	assert.NoError(t, err)
	
	assert.True(t, successResp.Success, "يجب أن يكون success true")
	assert.Contains(t, successResp.Message, "created successfully", "يجب أن تحتوي الرسالة على created successfully")
	assert.Equal(t, "محمد علي", createdUser.Name, "يجب أن يكون الاسم مطابقاً")
	assert.Equal(t, "mohamed@example.com", createdUser.Email, "يجب أن يكون email مطابقاً")
	assert.NotZero(t, createdUser.ID, "يجب أن يكون هناك ID")
	
	// التحقق من وجود المستخدم في database
	var userInDB models.User
	err = database.GetMongoDB().Collection("users").FindOne(context.Background(), bson.M{"_id": createdUser.ID}).Decode(&userInDB)
	assert.NoError(t, err, "يجب أن يوجد المستخدم في database")
	assert.Equal(t, "محمد علي", userInDB.Name, "يجب أن يكون الاسم في database مطابقاً")
	assert.Equal(t, "mohamed@example.com", userInDB.Email, "يجب أن يكون email في database مطابقاً")
}

// باقي الاختبارات بدون تغيير...

func TestGetUserByID(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)
	
	user := models.User{
		Name:     "مستخدم للاختبار", 
		Email:    "testuser@example.com", 
		Password: "testpassword",
	}
	_, err := database.GetMongoDB().Collection("users").InsertOne(context.Background(), user)
	assert.NoError(t, err)

	// ✅ استخدام fmt.Sprintf للتحويل الصحيح
	req, err := http.NewRequest("GET", fmt.Sprintf("/users/%d", user.ID), nil)
	assert.NoError(t, err)
	
	rr := httptest.NewRecorder()
	handler := handlers.NewUserHandler()
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
	
	router.ServeHTTP(rr, req)
	
	assert.Equal(t, http.StatusOK, rr.Code, "يجب أن يرجع 200 للمستخدم الموجود")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	
	var responseUser models.UserResponse
	successResp, err := parseSuccessResponse(rr.Body.Bytes(), &responseUser)
	assert.NoError(t, err)
	
	assert.True(t, successResp.Success)
	assert.Equal(t, user.Name, responseUser.Name)
	assert.Equal(t, user.Email, responseUser.Email)
}

func TestGetUserByID_NotFound(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)
	
	// ✅ استخدام fmt.Sprintf للتحويل الصحيح
	req, err := http.NewRequest("GET", fmt.Sprintf("/users/%d", 999), nil) // ID غير موجود
	assert.NoError(t, err)
	
	rr := httptest.NewRecorder()
	handler := handlers.NewUserHandler()
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
	
	router.ServeHTTP(rr, req)
	
	// ✅ يجب أن يرجع 404 Not Found
	assert.Equal(t, http.StatusNotFound, rr.Code, "يجب أن يرجع 404 للمستخدم غير الموجود")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	
	errorResp, err := parseErrorResponse(rr.Body.Bytes())
	assert.NoError(t, err)
	assert.False(t, errorResp.Success)
	assert.Contains(t, errorResp.Message, "not found", "يجب أن تحتوي الرسالة على خطأ غير موجود")
}