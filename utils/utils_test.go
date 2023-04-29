package utils

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAppendContextError(t *testing.T) {
	t.Run("Check if error is appended to context", func(t *testing.T) {
		err := errors.New("any error")
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		AppendContextError(ctx, err)
		assert.Equal(t, ctx.Errors[0].Err, err)
	})
}

func TestGetJsonStringFromStruct(t *testing.T) {
	t.Run("Check if struct is converted to valid JSON", func(t *testing.T) {
		var sample = struct {
			FieldA string
			FieldB int
			FieldC []int
		}{}
		sample.FieldA = "a"
		sample.FieldB = 1
		sample.FieldC = []int{2, 3}
		jsonString, err := GetJsonStringFromStruct(sample)
		assert.Nil(t, err)
		assert.Equal(t, jsonString, `{"FieldA":"a","FieldB":1,"FieldC":[2,3]}`)
	})
}
