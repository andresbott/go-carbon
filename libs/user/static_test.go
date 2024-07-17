package user

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestUserFromFile(t *testing.T) {
	tcs := []struct {
		name   string
		user   string
		pw     string
		expect bool
	}{
		{
			name:   "expect login true",
			user:   "demo",
			pw:     "demo",
			expect: true,
		},
		{
			name:   "expect  login false on disabled account",
			user:   "user1",
			pw:     "1234",
			expect: false,
		},
		{
			name:   "expect  login false on wrong password account",
			user:   "user1",
			pw:     "12345",
			expect: false,
		},
	}

	files := map[string]string{
		"yaml":     "testdata/users.yaml",
		"json":     "testdata/users.json",
		"htpasswd": "testdata/users.htpasswd",
	}
	for k, v := range files {
		t.Run(k, func(t *testing.T) {
			file := v
			users, err := FromFile(file)
			if err != nil {
				t.Fatal(err)
			}

			for _, tc := range tcs {
				t.Run(tc.name, func(t *testing.T) {
					got := users.AllowLogin(tc.user, tc.pw)
					if diff := cmp.Diff(got, tc.expect); diff != "" {
						t.Errorf("unexpected value (-got +want)\n%s", diff)
					}
				})
			}
		})
	}

}
