package gravatar_test

import (
	"bytes"
	"fmt"
	gr "github.com/ftrvxmtrx/gravatar"
	"image"
	_ "image/png"
)

func ExampleGetAvatar() {
	// get avatar image (128x128) using HTTP transport
	emailHash := gr.EmailHash("ftrvxmtrx@gmail.com")
	raw, err := gr.GetAvatar("http", emailHash, 128)

	// get avatar image (32x32) using HTTP transport
	// allow images of any rating level
	raw, err = gr.GetAvatar("http", emailHash, gr.RatingX, 32)

	// get avatar image (default size, png format) with fallback to "retro"
	// generated avatar.
	// use HTTPS transport
	emailHash = "cfcd208495d565ef66e7dff9f98764da.png"
	raw, err = gr.GetAvatar("https", emailHash, gr.DefaultRetro)

	if err == nil {
		var cfg image.Config
		var format string

		rawb := bytes.NewReader(raw)
		cfg, format, err = image.DecodeConfig(rawb)
		fmt.Println(cfg, format)
	}
}

func ExampleGetAvatarURL() {
	// get URL to avatar image of size 256x256
	// fall back to "monster" generated avatar
	emailHash := gr.EmailHash("ftrvxmtrx@gmail.com")
	url := gr.GetAvatarURL("https", emailHash, gr.DefaultMonster, 256)
	fmt.Println(url.String())
}

func ExampleGetProfile() {
	// get profile using HTTPS transport
	emailHash := gr.EmailHash("ftrvxmtrx@gmail.com")
	profile, err := gr.GetProfile("https", emailHash)

	if err == nil {
		fmt.Println(profile.PreferredUsername)
		fmt.Println(profile.ProfileUrl)
	}
}

func ExampleSetAvatarURLOptions() {
	// get URL to avatar image of default size
	emailHash := gr.EmailHash("ftrvxmtrx@gmail.com")
	url := gr.GetAvatarURL("https", emailHash)
	fmt.Printf("default URL: %s", url.String())
	// set size to 256x256
	// fall back to "monster" generated avatar
	gr.SetAvatarURLOptions(url, gr.DefaultMonster, 256)
	fmt.Printf("modified URL: %s", url.String())
	// reset back to the default one
	gr.SetAvatarURLOptions(url)
	fmt.Printf("URL after reset: %s", url.String())
}
