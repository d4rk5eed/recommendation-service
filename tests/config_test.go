package tests

import (
	"reflect"
	"testing"

	"recommedation-service/pkg/models"
)

func TestReadConfig(t *testing.T) {
	conf, err := models.ReadConfig("../config/test.yaml")
	if err != nil {
		t.Errorf("Error reading config file: %v", err)
	}
	expected := models.ServiceConfig{
		"openai": {
			"url":   "https://api.openai.com/v1/",
			"name":  "openai",
			"key":   "test-key-123",
			"model": "gpt-3.5-turbo",
		},
		"mock-llm": {
			"url":   "http://localhost:8081/v1/",
			"name":  "mock-llm",
			"key":   "test-key",
			"model": "mock",
		},
	}

	if !reflect.DeepEqual(conf, expected) {
		t.Errorf("Error comparing: %v, %v", expected, conf)
	}
}
