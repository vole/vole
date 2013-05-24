package gravatar

import (
	"testing"
)

// TODO(ftrvxmtrx): write real tests

func TestEmailHash(t *testing.T) {
	h := EmailHash("ftrvxmtrx@gmail.com")

	if h != "d96ba36eb0d406aea53f3868cd06fca8" {
		t.Error(h)
	}
}

func TestGetAvatar(t *testing.T) {
	for _, scheme := range []string{"http", "https"} {
		if a, err := GetAvatar(scheme, "d96ba36eb0d406aea53f3868cd06fca8"); err != nil {
			t.Error(a, err)
		}

		if a, err := GetAvatar(scheme, "0"); err != nil {
			t.Error(a, err)
		}

		if a, err := GetAvatar(scheme, "0", DefaultError); err == nil {
			t.Error(a, err)
		}

		if a, err := GetAvatar(scheme, "0.png", DefaultIdentIcon, 256); err != nil {
			t.Error(a, err)
		}

		if a, err := GetAvatar(scheme, "0.png", RatingX, DefaultIdentIcon, 256); err != nil {
			t.Error(a, err)
		}
	}
}

func TestGetAvatarURL(t *testing.T) {
	if url := GetAvatarURL("http", "d96ba36eb0d406aea53f3868cd06fca8"); url == nil {
		t.Error(url)
	}
}

func TestGetProfile(t *testing.T) {
	for _, scheme := range []string{"http", "https"} {
		if e, err := GetProfile(scheme, "d96ba36eb0d406aea53f3868cd06fca8"); err != nil {
			t.Error(e, err)
		}

		if e, err := GetProfile(scheme, "0"); err == nil {
			t.Error(e, err)
		}
	}
}

func TestSetProfile(t *testing.T) {
	if url := GetAvatarURL("http", "d96ba36eb0d406aea53f3868cd06fca8"); url == nil {
		t.Error(url)
	} else {
		values := SetAvatarURLOptions(url, 123, DefaultRetro).Query()

		if values.Get("s") != "123" {
			t.Error(values.Get("s"))
		}

		if values.Get("d") != DefaultRetro {
			t.Error(values.Get("d"))
		}

		values = SetAvatarURLOptions(url).Query()

		if values.Get("s") != "" {
			t.Error(values.Get("s"))
		}

		if values.Get("d") != "" {
			t.Error(values.Get("d"))
		}
	}
}
