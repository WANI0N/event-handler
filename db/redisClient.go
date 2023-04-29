package db

import (
	"app/structs"
	"app/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// redis is not a backend database and I wouldn't recommend it for storing long term or table/dict data,
// but since this repo has only 1 endpoint and any db selected would need to have a database layer
// I decided to use it since I can demonstrate the layer and unit tests with it and it's easy to setup

var redisEndpoint string = os.Getenv("REDIS_ENDPOINT")
var redisPassword string = os.Getenv("REDIS_PASSWORD")
var ctx = context.Background()
var redisClient = Init()

func Init() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:        redisEndpoint,
		Password:    redisPassword,
		DialTimeout: time.Second,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Logger.Error().Msg(fmt.Sprintf(
			"error connecting to redis cache, err: '%v'", err))
	} else {
		log.Logger.Info().Msg("redis connection successful")
	}
	return rdb
}

var GetEvent = func(id string) (structs.EventData, error) {
	result, err := redisClient.Get(ctx, id).Result()
	if err != nil {
		if err == redis.Nil {
			return structs.EventData{}, errors.New("not found")
		}
		return structs.EventData{}, errors.New(
			"redis connection error",
		)
	}
	var eventData structs.EventData
	jsonParseErr := json.Unmarshal([]byte(result), &eventData)
	if jsonParseErr != nil {
		return structs.EventData{}, jsonParseErr
	}
	return eventData, nil
}

var CreateEvent = func(payload structs.EventData) (string, error) {
	eventId := uuid.NewString()
	payload.Id = eventId
	dataAsJsonString, convertErr := utils.GetJsonStringFromStruct(payload)
	if convertErr != nil {
		log.Logger.Error().Msgf("error converting data to json: %v", convertErr)
		return "", convertErr
	}
	err := redisClient.Set(ctx, eventId, dataAsJsonString, 0).Err()
	if err != nil {
		log.Logger.Error().Msgf("error on setting data to redis: %v", err)
		return "", err
	}
	return eventId, nil
}

var DeleteEvent = func(id string) error {
	result, err := redisClient.Del(ctx, id).Result()
	if err != nil {
		return err
	}
	if result == 0 {
		return errors.New("not found")
	}
	return nil
}
