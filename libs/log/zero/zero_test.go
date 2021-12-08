package zero_test

import (
	"bytes"
	"git.andresbott.com/Golang/carbon/libs/log/zero"
	"strings"
	"testing"
)

func TestZero(t *testing.T) {

	t.Run("log multiple fields", func(t *testing.T) {
		var buf bytes.Buffer
		z := zero.NewZero(zero.InfoLevel, &buf)

		z.Info("msg", "k", "v", "k2", "v2")

		_ = z
		output := buf.String()
		expects := []string{
			`{"level":"info","k":"v","k2":"v2",`,
			`"message":"msg"}`,
		}

		for _, expect := range expects {
			if output == "" || !strings.Contains(output, expect) {
				t.Errorf("expected output is empty or does not contain expected string; got %s", output)
			}
		}
	})

	t.Run("log multiple odd fields", func(t *testing.T) {
		var buf bytes.Buffer
		z := zero.NewZero(zero.InfoLevel, &buf)

		z.Info("msg", "v1", "v2", "v3")

		_ = z
		output := buf.String()
		expects := []string{
			`{"level":"info","data":"v1,v2,v3",`,
			`"message":"msg"}`,
		}

		for _, expect := range expects {
			if output == "" || !strings.Contains(output, expect) {
				t.Errorf("expected output is empty or does not contain expected string; got %s", output)
			}
		}
	})

}
