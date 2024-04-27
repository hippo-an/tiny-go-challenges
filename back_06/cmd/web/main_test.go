package main

import "testing"

func TestRun(t *testing.T) {
	err := setupAppConfig()

	if err != nil {
		t.Error("failed setupAppConfig()")
	}
}
