package config

import "testing"

func TestConfig(t *testing.T) {
	SetFilePath("testdata/config.yaml")
	config := Get()
	if config.SlackToken != "YOUR_BOT_TOKEN_HERE" {
		t.Fatal("TestConfig: can't read config file")
	}
}
