package telegramApi

import (
	"github.com/oleksiy-os/porto-events/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_truncateString(t *testing.T) {
	tests := []struct {
		name string
		args model.Event
		msg  string
		want string
	}{
		{
			name: "dot with space",
			args: model.Event{
				Description: "aliquet lectus proin nibh nisl condimentum id venenatis a condimentum vitae sapien pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas sed tempus urna et pharetra pharetra massa massa ultricies mi quis hendrerit dolor magna eget est lorem ipsum dolor sit amet consectetur adipiscing elit pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas integer eget aliquet nibh praesent tristique magna sit amet purus gravida quis blandit turpis cursus in hac habitasse platea dictumst quisque sagittis purus sit amet volutpat consequat mauris nunc congue nisi vitae suscipit tellus mauris a diam maecenasaliquet lectus proin nibh nisl condimentum id venenatis a condimentum vitae sapien pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas sed tempus urna et pharetra pharetra. massa massa ultricies mi quis hendrerit dolor magna eget est lorem ipsum dolor sit amet consectetur.adipiscing. elit. here the end. some text. text.",
				Url:         "1234",
				Title:       "12",
				LocationMap: "aaa bbb",
			},
			msg:  `<b><a href="%s">%s</a></b> &#10;%s &#10;📍 <a href="%s">%s</a> &#10;🗓 %s &#10;🕒 %s &#10;%s`,
			want: "aliquet lectus proin nibh nisl condimentum id venenatis a condimentum vitae sapien pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas sed tempus urna et pharetra pharetra massa massa ultricies mi quis hendrerit dolor magna eget est lorem ipsum dolor sit amet consectetur adipiscing elit pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas integer eget aliquet nibh praesent tristique magna sit amet purus gravida quis blandit turpis cursus in hac habitasse platea dictumst quisque sagittis purus sit amet volutpat consequat mauris nunc congue nisi vitae suscipit tellus mauris a diam maecenasaliquet lectus proin nibh nisl condimentum id venenatis a condimentum vitae sapien pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas sed tempus urna et pharetra pharetra.",
		},
		{
			name: "shorter than limit",
			args: model.Event{
				Description: "aliquet lectus. proin. the end",
				Url:         "1234",
				Title:       "12",
				LocationMap: "aaa bbb",
			},
			want: "aliquet lectus. proin. the end",
		},
		{
			name: "empty",
			args: model.Event{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, truncateString(&tt.args, &tt.msg))
		})
	}
}
