package gravatar

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GravatarProfile stores profile information associated with an email.
type GravatarProfile struct {
	AboutMe string

	Accounts []struct {
		ShortName string
		Domain    string
		Url       string
		Verified  bool `json:",string"`
		Username  string
		Display   string
	}

	CurrentLocation string

	DisplayName string

	Emails []struct {
		Primary bool `json:",string"`
		Value   string
	}

	Hash string

	Id int `json:",string"`

	Ims []struct {
		Type  string
		Value string
	}

	Name struct {
		Family    string
		Formatted string
		Given     string
	}

	PhoneNumbers []struct {
		Type  string
		Value string
	}

	Photos []struct {
		Type  string
		Value string
	}

	PreferredUsername string

	ProfileBackground struct {
		Color    string
		Position string
		Repeat   string
		Url      string
	}

	ProfileUrl string

	ThumbnailUrl string

	Urls []struct {
		Title string
		Value string
	}
}

type rating string

const gravatarHost = "gravatar.com"

// List of (optional) values for a "default action" option for GetAvatar.
// Each option defines what GetAvatar has to do in case of non-existing
// user.
const (
	// DefaultBlank defaults to a transparent PNG image.
	DefaultBlank = "blank"

	// DefaultError defaults to an error.
	DefaultError = "404"

	// DefaultIdentIcon defaults to a generated geometric pattern.
	DefaultIdentIcon = "identicon"

	// DefaultMonster defaults to a generated 'monster' with different colors
	// and faces.
	DefaultMonster = "monsterid"

	// DefaultMysteryMan defaults to a simple, cartoon-style silhouetted outline
	// of a person.
	DefaultMysteryMan = "mm"

	// DefaultRetro defaults to a generated 8-bit arcade-style pixelated faces.
	DefaultRetro = "retro"

	// DefaultWavatar defaults to a generated faces with differing features and
	// backgrounds.
	DefaultWavatar = "wavatar"
)

// List of (optional) values to specify allowed rating (up to and including
// that).
// If the requested email doesn't have any image of allowed level, a default
// image will be used.
const (
	// RatingG is suitable for display on all websites with any audience type.
	RatingG = rating("g")

	// RatingPG may contain rude gestures, provocatively dressed individuals, the
	// lesser swear words, or mild violence.
	RatingPG = rating("pg")

	// RatingR may contain such things as harsh profanity, intense violence,
	// nudity, or hard drug use.
	RatingR = rating("r")

	// RatingX may contain hardcore sexual imagery or extremely disturbing
	// violence.
	RatingX = rating("x")
)

var client = new(http.Client)

// EmailHash converts an email to lowercase and returns its MD5 hash as hex
// string.
func EmailHash(email string) string {
	m := md5.New()
	io.WriteString(m, strings.ToLower(email))
	return fmt.Sprintf("%x", m.Sum(nil))
}

// GetAvatar does a HTTP(S) request and returns an avatar image.
//
// Optional arguments include Default* (default actions), image size and
// Rating* (rating level, default is RatingG).
//
// Instead of Default* predefined constants you may also use a direct URL to an
// image.
func GetAvatar(scheme, emailHash string, opts ...interface{}) (data []byte, err error) {
	url := GetAvatarURL(scheme, emailHash, opts...)
	err = run(url, get_avatar(&data))
	return
}

// GetAvatarURL returns an URL to avatar image.
//
// Optional arguments include Default* (default actions), image size and
// Rating* (rating level, default is RatingG).
//
// Instead of Default* predefined constants you may also use a direct URL to an
// image.
func GetAvatarURL(scheme, emailHash string, opts ...interface{}) *url.URL {
	url := &url.URL{
		Scheme: scheme,
		Host:   gravatarHost,
		Path:   "/avatar/" + emailHash,
	}

	return SetAvatarURLOptions(url, opts...)
}

// GetProfile does a HTTP(S) request and returns gravatar profile.
func GetProfile(scheme, emailHash string) (g GravatarProfile, err error) {
	url := &url.URL{
		Scheme: scheme,
		Host:   gravatarHost,
		Path:   "/" + emailHash + ".json",
	}

	err = run(url, unmarshal_json(&g))
	return
}

// SetAvatarURLOptions sets options for an URL to avatar image.
//
// Options include Default* (default actions), image size and Rating* (rating
// level, default is RatingG).
//
// Instead of Default* predefined constants you may also use a direct URL to an
// image.
//
// Calling SetAvatarURLOptions(url), e.g. without any options, resets the
// options.
func SetAvatarURLOptions(u *url.URL, opts ...interface{}) *url.URL {
	values := make(url.Values)

	for _, opt := range opts {
		switch o := opt.(type) {
		case int:
			values.Set("s", strconv.Itoa(o))

		case rating:
			values.Set("r", string(o))

		case string:
			values.Set("d", o)
		}
	}

	u.RawQuery = values.Encode()
	return u
}

func get_avatar(dst *[]byte) func([]byte) error {
	return func(data []byte) (err error) {
		*dst = data[:]
		return
	}
}

func run(url *url.URL, f func([]byte) error) (err error) {
	var res *http.Response
	res, err = client.Get(url.String())

	if err == nil {
		var data []byte
		defer res.Body.Close()

		if data, err = ioutil.ReadAll(res.Body); err == nil {
			if res.StatusCode == http.StatusOK {
				err = f(data)
			} else {
				err = errors.New(string(data))
			}
		}
	}

	return
}

func unmarshal_json(g *GravatarProfile) func([]byte) error {
	return func(data []byte) (err error) {
		obj := struct {
			Entries []GravatarProfile `json:"entry"`
		}{[]GravatarProfile{}}

		if err = json.Unmarshal(data, &obj); err == nil && len(obj.Entries) > 0 {
			*g = obj.Entries[0]
		}

		return
	}
}
