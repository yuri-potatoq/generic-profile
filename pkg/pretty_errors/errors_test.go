package pretty_errors

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestPrettyError_Error(t *testing.T) {
	for i, tc := range []struct {
		err         error
		expectedStr string
	}{
		{
			err: PrettyError{
				Details: map[string]any{
					"userName": "test",
					"userId":   1,
					"userAttributes": map[string]string{
						"newsletter": "true",
					},
				},
				Message:  "generic error",
				Path:     "user.service.UpdateUser",
				InnerErr: os.ErrClosed,
			},
			expectedStr: "[generic error] from [user.service.UpdateUser] with []",
		},
	} {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			repr := tc.err.Error()
			if strings.Compare(repr, tc.expectedStr) != 0 {
				t.Errorf("not expected representations => expected={%s},given={%s}", tc.expectedStr, repr)
			}
		})
	}
}
