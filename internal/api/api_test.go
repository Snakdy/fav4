package api

import "testing"
import "github.com/stretchr/testify/assert"

func TestNewIconAPI(t *testing.T) {
	api := NewIconAPI()
	assert.NotNil(t, api)
}

func TestIconAPI_parse(t *testing.T) {
	var cases = []struct {
		in  string
		out string
		err bool
	}{
		{
			"",
			"",
			true,
		},
		{
			"http://google.com",
			"https://google.com",
			false,
		},
		{
			"github.com",
			"https://github.com",
			false,
		},
		{
			"foo.bar/foobar",
			"https://foo.bar/foobar",
			false,
		},
		{
			"- --",
			"",
			true,
		},
	}
	api := &IconAPI{}
	for _, tt := range cases {
		t.Run(tt.in, func(t *testing.T) {
			out, err := api.parse(tt.in)
			if tt.err {
				assert.Nil(t, out)
				return
			}
			assert.NotNil(t, out)
			assert.EqualValues(t, tt.out, out.String())
			assert.EqualValues(t, tt.err, err != nil)
		})
	}
}
