package zeroGorm_test

import (
	"bytes"
	"git.andresbott.com/Golang/carbon/libs/log/zeroGorm"
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
	"testing"
)

type SampleData struct {
	gorm.Model
	Nanme string
}

func TestGormLog(t *testing.T) {

	t.Run("expect error", func(t *testing.T) {

		var buf bytes.Buffer
		l := zerolog.New(&buf).With().Timestamp().Logger().Level(zerolog.InfoLevel)
		zl := zeroGorm.New(&l, zeroGorm.Cfg{})

		db, _ := gorm.Open(sqlite.Open("file::memory:?cache=bla"), &gorm.Config{
			Logger: zl,
		})
		_ = db

		output := buf.String()
		expect := "\"level\":\"error\""
		if output == "" || !strings.Contains(output, expect) {
			t.Errorf("expected output is not of type error; got %s", output)
		}
		expect = "failed to initialize database, got error no such cache mode: bla"
		if output == "" || !strings.Contains(output, expect) {
			t.Errorf("expected output is empty or does not contain expected string; got %s", output)
		}

	})

	t.Run("expect Error no such table", func(t *testing.T) {
		var buf bytes.Buffer
		l := zerolog.New(&buf).With().Timestamp().Logger().Level(zerolog.InfoLevel)
		zl := zeroGorm.New(&l, zeroGorm.Cfg{})
		_ = zl

		db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
			Logger: zl,
		})

		data := SampleData{
			Model: gorm.Model{},
			Nanme: "bla",
		}
		db.Create(&data)

		output := buf.String()
		expect := "\"level\":\"error\""
		if output == "" || !strings.Contains(output, expect) {
			t.Errorf("expected output is not of type error; got %s", output)
		}
		expect = "\"message\":\"no such table: sample_data\""
		if output == "" || !strings.Contains(output, expect) {
			t.Errorf("expected output is empty or does not contain expected string; got %s", output)
		}
	})
}
