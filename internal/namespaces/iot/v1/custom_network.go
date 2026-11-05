package iot

import (
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
)

var iotNetworkCreateLayout = humanv2.Layout{
	Sections: []humanv2.Section{
		{
			Path:   "Network",
			Title:  "Network",
			Footer: "WARNING: The secret below is displayed only once, we do not keep it. Make sure to store it securely",
		},
		{
			Path:  "Secret",
			Title: "Secret",
		},
	},
}
