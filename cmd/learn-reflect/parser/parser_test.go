package parser

import "testing"

func TestExtractFields(t *testing.T) {
	cfg := struct {
		Version string
		Build   string
		Web     struct {
			APIHost string `conf:"default:test-host"`
			Port    string `conf:"default:test-port"`
		}
		DB struct {
			ConnStr string `conf:"default:test-conn-str"`
			Nested  struct {
				TestField    string `conf:"default:test-field"`
				TestIntField int    `conf:"default:10"`
			}
		}
	}{}

	fields, err := ExtractFields(&cfg)
	if err != nil {
		t.Errorf("should not have failed. error: %v", err)
	}
	for _, v := range fields {
		t.Log(v.Name)
		t.Log(v.Key)
		t.Log("---")
	}
	t.Log(cfg)
}
