package weberrors

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type ValidationCheckStruct struct {
	Field_a string `json:"field_a" binding:"required"`
}

func TestJsonAppErrorReporter(t *testing.T) {
	r := gin.New()
	r.Use(JSONAppErrorReporter())

	t.Run("no errors", func(t *testing.T) {
		emptyResponse := func(c *gin.Context) {}
		expectedError := AppError{}
		r.GET("/no-error", emptyResponse)

		w := httptest.NewRecorder()
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/no-error", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		decodedError := AppError{}
		json.NewDecoder(w.Body).Decode(&decodedError)
		assert.Equal(t, expectedError, decodedError)
	})

	t.Run("case view ValidationErrors", func(t *testing.T) {
		errorResponse := func(c *gin.Context) {
			fmt.Println("JERE")
			bindEntity := ValidationCheckStruct{}
			if err := c.ShouldBindJSON(&bindEntity); err != nil {
				fmt.Println("JERE2")
				c.Error(err)
				c.Abort()
			}
		}

		entity := ValidationCheckStruct{}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(entity)
		expectedError := ParseAppError(
			ValidationError.ChangeDesc("Field `field_a` is required."),
		)

		r.POST("/validation-error", errorResponse)
		w := httptest.NewRecorder()
		req, _ := http.NewRequestWithContext(context.Background(), "POST", "/validation-error", &buf)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		decodedError := AppError{}
		json.NewDecoder(w.Body).Decode(&decodedError)
		assert.Equal(t, expectedError, decodedError)
	})

	t.Run("case AppErrorWithCode (NotFound)", func(t *testing.T) {
		errorResponse := func(c *gin.Context) {
			c.Error(&NotFound)
		}
		expectedError := ParseAppError(&NotFound)
		r.GET("/not-found", errorResponse)

		w := httptest.NewRecorder()
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/not-found", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		decodedError := AppError{}
		json.NewDecoder(w.Body).Decode(&decodedError)
		assert.Equal(t, expectedError, decodedError)
	})

	t.Run("case default", func(t *testing.T) {
		errorResponse := func(c *gin.Context) {
			c.Error(errors.New("Unexpected error."))
		}
		expectedError := ParseAppError(InternalError.ChangeDesc("Internal Server Error"))
		r.GET("/default-error", errorResponse)

		w := httptest.NewRecorder()
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/default-error", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		decodedError := AppError{}
		json.NewDecoder(w.Body).Decode(&decodedError)
		assert.Equal(t, expectedError, decodedError)
	})

}
