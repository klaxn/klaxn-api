package integration

import (
	"os"
	"testing"

	"github.com/klaxn/klaxn-api/pkg/client"
)

func setup(t *testing.T) *client.Client {
	if os.Getenv("TEST_ACC") != "true" {
		t.SkipNow()
	}

	c, err := client.New(getEnvOrDefault("KLAXN_URL", "http://localhost:8080"))
	if err != nil {
		t.Fatal(err)
	}

	services, err := c.GetServices()
	if err != nil {
		t.Fatal(err)
	}

	for _, service := range services {
		err := c.DeleteService(service.ID)
		if err != nil {
			t.Fatal(err)
		}
	}

	escalations, err := c.GetEscalations()
	if err != nil {
		t.Fatal(err)
	}

	for _, escalation := range escalations {
		err = c.DeleteEscalation(escalation.ID)
		if err != nil {
			t.Fatal(err)
		}
	}

	return c
}

func getEnvOrDefault(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}
