package db

import (
	"app/models"
	"app/utils"
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var redisServer, _ = miniredis.Run()

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()

	if err != nil {
		panic(err)
	}

	return s
}

func setup() *redis.Client {
	var redisServer = mockRedis()
	var redisClient = redis.NewClient(&redis.Options{
		Addr:        redisServer.Addr(),
		DialTimeout: time.Second,
	})
	return redisClient
}

func teardown() {
	iter := redisClient.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		redisClient.Del(ctx, iter.Val())
	}
	redisServer.Close()
}

type KeyValuePair struct {
	key   string
	value string
}

func insertDataToCache(client *redis.Client, data []KeyValuePair) {
	for _, record := range data {
		client.Set(ctx, record.key, record.value, 0)
	}
}

func sortDataByKey(data []KeyValuePair) []KeyValuePair {
	sort.Slice(data, func(i, j int) bool {
		return data[i].key < data[j].key
	})
	return data
}

func retrieveDataFromCache(client *redis.Client) []KeyValuePair {
	iter := client.Scan(ctx, 0, "*", 0).Iterator()
	data := []KeyValuePair{}
	var key string
	var value string
	for iter.Next(ctx) {
		key = iter.Val()
		value, _ = client.Get(ctx, key).Result()
		data = append(data, KeyValuePair{
			key:   key,
			value: value,
		})
	}
	return data
}

var eventDataAsJsonString = `{
	"id": "event-id-string",
	"name": "My Event",
	"date": "2023-04-20T14:00:00Z",
	"languages": ["English", "French"],
	"videoQuality": ["720p", "1080p"],
	"audioQuality": ["High", "Low"],
	"invitees": ["example1@gmail.com", "example2@gmail.com"],
	"description": "A short description of the event"
}`
var eventDataAsStruct = models.EventData{
	Name:         "My Event",
	Timestamp:    "2023-04-20T14:00:00Z",
	Languages:    []string{"English", "French"},
	VideoQuality: []string{"720p", "1080p"},
	AudioQuality: []string{"High", "Low"},
	Invitees:     []string{"example1@gmail.com", "example2@gmail.com"},
	Description:  "A short description of the event",
}

var eventRespDataAsStruct = models.EventResponseData{
	Id:        "event-id-string",
	EventData: eventDataAsStruct,
}

var GetEventTestCases = []struct {
	description   string
	setupRedis    bool
	innitialCache []KeyValuePair
	submitId      string
	expectedResp  models.EventResponseData
	expectedError error
}{
	{
		description: "Success",
		innitialCache: []KeyValuePair{
			{
				key:   "event-id-string",
				value: eventDataAsJsonString,
			},
		},
		submitId:     "event-id-string",
		expectedResp: eventRespDataAsStruct,
	},
	{
		description:   "Fail - not found",
		submitId:      "non-existent-id",
		expectedError: errors.New("not found"),
	},
}

func TestGetEvent(t *testing.T) {
	for _, testCase := range GetEventTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			setup()
			defer teardown()
			insertDataToCache(redisClient, testCase.innitialCache)
			resp, err := GetEvent(testCase.submitId)
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResp, resp)
		})
	}
}

var CreateEventTestCases = []struct {
	description                     string
	getJsonStringFromStructMockResp string
	getJsonStringFromStructMockErr  error
	innitialCache                   []KeyValuePair
	submitPayload                   models.EventData
	expectedCache                   []KeyValuePair
	expectRespId                    bool
	expectedError                   error
}{
	{
		description: "Success",
		innitialCache: []KeyValuePair{
			{
				key:   "id-control",
				value: "content",
			},
		},
		getJsonStringFromStructMockResp: eventDataAsJsonString,
		submitPayload:                   eventDataAsStruct,
		expectRespId:                    true,
		expectedCache: []KeyValuePair{
			{
				key:   "key-to-replace-with-uuid",
				value: eventDataAsJsonString,
			},
			{
				key:   "id-control",
				value: "content",
			},
		},
	},
	{
		description:                    "Fail - Json build fail",
		submitPayload:                  eventDataAsStruct,
		getJsonStringFromStructMockErr: errors.New("any error"),
		expectedError:                  errors.New("any error"),
	},
}

func TestCreateEvent(t *testing.T) {
	for _, testCase := range CreateEventTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			setup()
			defer teardown()
			utils.GetJsonStringFromStruct = func(data interface{}) (string, error) {
				convertedData := data.(models.EventData)
				convertedData.Id = ""
				testCase.submitPayload.Id = ""
				assert.Equal(t, testCase.submitPayload, convertedData)
				return testCase.getJsonStringFromStructMockResp, testCase.getJsonStringFromStructMockErr
			}
			insertDataToCache(redisClient, testCase.innitialCache)

			respId, err := CreateEvent(testCase.submitPayload)

			if testCase.expectRespId {
				_, uuIderr := uuid.Parse(respId)
				assert.Nil(t, uuIderr)
			} else {
				assert.Empty(t, respId)
			}
			assert.Equal(t, testCase.expectedError, err)
			if len(testCase.expectedCache) > 0 {
				testCase.expectedCache[0].key = respId
				cacheContents := retrieveDataFromCache(redisClient)
				assert.Equal(t,
					sortDataByKey(testCase.expectedCache),
					sortDataByKey(cacheContents),
				)
			}
		})
	}
}

var DeleteEventTestCases = []struct {
	description   string
	innitialCache []KeyValuePair
	expectedCache []KeyValuePair
	submitId      string
	expectedError error
}{
	{
		description: "Success",
		innitialCache: []KeyValuePair{
			{
				key:   "id-to-keep",
				value: "content-1",
			},
			{
				key:   "id-to-delete",
				value: "content-2",
			},
		},
		expectedCache: []KeyValuePair{
			{
				key:   "id-to-keep",
				value: "content-1",
			},
		},
		submitId: "id-to-delete",
	},
	{
		description: "Fail - key does not exist",
		innitialCache: []KeyValuePair{
			{
				key:   "id-1",
				value: "content-1",
			},
			{
				key:   "id-2",
				value: "content-2",
			},
		},
		expectedCache: []KeyValuePair{
			{
				key:   "id-1",
				value: "content-1",
			},
			{
				key:   "id-2",
				value: "content-2",
			},
		},
		submitId:      "non-existent-id",
		expectedError: errors.New("not found"),
	},
}

func TestDeleteEvent(t *testing.T) {
	for _, testCase := range DeleteEventTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			setup()
			defer teardown()
			insertDataToCache(redisClient, testCase.innitialCache)
			err := DeleteEvent(testCase.submitId)
			assert.Equal(t, testCase.expectedError, err)
			cacheContents := retrieveDataFromCache(redisClient)
			assert.Equal(t,
				sortDataByKey(testCase.expectedCache),
				sortDataByKey(cacheContents),
			)
		})
	}
}
