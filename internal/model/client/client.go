package client

import "github.com/oleksiy-os/porto-events/internal/model"

//goland:noinspection GoNameStartsWithPackageName
type ClientI interface {
	Publish(*[]model.Event)
}
