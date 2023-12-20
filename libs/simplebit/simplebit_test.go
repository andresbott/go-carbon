package simplebit

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGet(t *testing.T) {

	tcs := []struct {
		name   string
		in     byte
		expect []bool
	}{
		{
			name: "number 156",
			in:   156, // 10011100

			expect: []bool{
				false,
				false,
				true,
				true,
				true,
				false,
				false,
				true,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(fmt.Sprintf("bit: %08b", tc.in))
			got := []bool{}
			for i := 0; i < 8; i++ {
				value := Get(tc.in, i)
				got = append(got, value)
			}
			if diff := cmp.Diff(got, tc.expect); diff != "" {
				t.Errorf("unexpected value (-got +want)\n%s", diff)
			}
		})
	}
}
func TestSet(t *testing.T) {

	tcs := []struct {
		name     string
		in       byte
		position int
		value    bool
		expect   []bool
	}{
		{
			name:     "number 156",
			in:       156, //  10011100
			position: 5,   // 00X00000
			value:    true,

			expect: []bool{
				false,
				false,
				true,
				true,
				true,
				true,
				false,
				true,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			newByte := Set(tc.in, tc.position)

			t.Log(fmt.Sprintf("bit: %08b", newByte))
			got := []bool{}
			for i := 0; i < 8; i++ {
				value := Get(newByte, i)
				got = append(got, value)
			}
			if diff := cmp.Diff(got, tc.expect); diff != "" {
				t.Errorf("unexpected value (-got +want)\n%s", diff)
			}
		})
	}
}

func TestClear(t *testing.T) {

	tcs := []struct {
		name     string
		in       byte
		position int
		value    bool
		expect   []bool
	}{
		{
			name:     "number 156",
			in:       156, //  10011100
			position: 2,   //  00000X00
			value:    true,

			expect: []bool{
				false,
				false,
				false,
				true,
				true,
				false,
				false,
				true,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			newByte := Clear(tc.in, tc.position)

			t.Log(fmt.Sprintf("bit: %08b", newByte))
			got := []bool{}
			for i := 0; i < 8; i++ {
				value := Get(newByte, i)
				got = append(got, value)
			}
			if diff := cmp.Diff(got, tc.expect); diff != "" {
				t.Errorf("unexpected value (-got +want)\n%s", diff)
			}
		})
	}
}
