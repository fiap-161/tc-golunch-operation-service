package shared

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils_BuildPath(t *testing.T) {
	type args struct {
		path   string
		params []BuildPathParam
	}

	type want struct {
		path string
		err  error
	}

	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "Given a invalid path and params, when BuildPath is called, then it should return error",
			args: args{
				path: "invalid",
				params: []BuildPathParam{
					{
						Key:   "userId",
						Value: "123",
					},
					{
						Key:   "orderId",
						Value: "456",
					},
				},
			},
			want: want{
				path: "",
				err:  errors.New("placeholder not found in path template"),
			},
		},
		{
			name: "Given a valid path and params, when BuildPath is called, then it should return the correct path",
			args: args{
				path: "/api/v1/users/{user_id}/orders/{order_id}",
				params: []BuildPathParam{
					{
						Key:   "user_id",
						Value: "123",
					},
					{
						Key:   "order_id",
						Value: "456",
					},
				},
			},
			want: want{
				path: "/api/v1/users/123/orders/456",
				err:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildPath(tt.args.path, tt.args.params)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.path, got)
		})
	}

}
