package players_test

import (
	"net/mail"
	"testing"

	"github.com/fernandoocampo/players/internal/appkit/unittests"
	"github.com/fernandoocampo/players/internal/players"
	"github.com/stretchr/testify/assert"
)

func TestValidateNewPlayer(t *testing.T) {
	testCases := map[string]struct {
		newPlayer players.NewPlayer
		isError   bool
		want      string
	}{
		"valid_player": {
			newPlayer: players.NewPlayer{
				FirstName: "Fernando",
				LastName:  "Ocampo",
				Nickname:  "belsonnoles",
				Email:     *unittests.NewEmailAddress(t, "belsonnoles@anyemail.com"),
				Password:  "A3BB605190830A01828F4D987A8C26CCBE3D4DAC0FAEF9482FBCB2B3CCB19CB8",
				Country:   "Spain",
			},
		},
		"invalid_email": {
			newPlayer: players.NewPlayer{
				FirstName: "Fernando",
				LastName:  "Ocampo",
				Nickname:  "belsonnoles",
				Password:  "A3BB605190830A01828F4D987A8C26CCBE3D4DAC0FAEF9482FBCB2B3CCB19CB8",
				Country:   "Spain",
			},
			isError: true,
			want:    "email is empty",
		},
		"invalid_empty_first_name": {
			newPlayer: players.NewPlayer{
				LastName: "Ocampo",
				Nickname: "belsonnoles",
				Email:    *unittests.NewEmailAddress(t, "belsonnoles@anyemail.com"),
				Password: "A3BB605190830A01828F4D987A8C26CCBE3D4DAC0FAEF9482FBCB2B3CCB19CB8",
				Country:  "Spain",
			},
			isError: true,
			want:    "first name is empty",
		},
		"invalid_empty_last_name": {
			newPlayer: players.NewPlayer{
				FirstName: "Fernando",
				Nickname:  "belsonnoles",
				Email:     *unittests.NewEmailAddress(t, "belsonnoles@anyemail.com"),
				Password:  "A3BB605190830A01828F4D987A8C26CCBE3D4DAC0FAEF9482FBCB2B3CCB19CB8",
				Country:   "Spain",
			},
			isError: true,
			want:    "last name is empty",
		},
		"invalid_empty_nickname": {
			newPlayer: players.NewPlayer{
				FirstName: "Fernando",
				LastName:  "Ocampo",
				Email:     *unittests.NewEmailAddress(t, "belsonnoles@anyemail.com"),
				Password:  "A3BB605190830A01828F4D987A8C26CCBE3D4DAC0FAEF9482FBCB2B3CCB19CB8",
				Country:   "Spain",
			},
			isError: true,
			want:    "nickname is empty",
		},
		"invalid_empty_country": {
			newPlayer: players.NewPlayer{
				FirstName: "Fernando",
				LastName:  "Ocampo",
				Nickname:  "belsonnoles",
				Email:     *unittests.NewEmailAddress(t, "belsonnoles@anyemail.com"),
				Password:  "A3BB605190830A01828F4D987A8C26CCBE3D4DAC0FAEF9482FBCB2B3CCB19CB8",
			},

			want:    "country is empty",
			isError: true,
		},
		"invalid_empty_password": {
			newPlayer: players.NewPlayer{
				FirstName: "Fernando",
				LastName:  "Ocampo",
				Nickname:  "belsonnoles",
				Email:     *unittests.NewEmailAddress(t, "belsonnoles@anyemail.com"),
				Country:   "Spain",
			},

			want:    "password is empty",
			isError: true,
		},
		"invalid_empty_multiple_fields": {
			newPlayer: players.NewPlayer{
				FirstName: "Fernando",
				LastName:  "Ocampo",
				Email:     *unittests.NewEmailAddress(t, "belsonnoles@anyemail.com"),
				Country:   "Spain",
			},

			want: `nickname is empty
password is empty`,
			isError: true,
		},
	}

	for testName, testData := range testCases {
		t.Run(testName, func(t *testing.T) {
			// When
			err := testData.newPlayer.Validate()

			// Then
			if testData.isError {
				assert.Error(t, err)
				assert.Equal(t, testData.want, err.Error())
			}

			if !testData.isError {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUpdatePlayer(t *testing.T) {
	testCases := map[string]struct {
		newPlayer players.UpdatePlayer
		isError   bool
		want      string
	}{
		"valid_player_update": {
			newPlayer: players.UpdatePlayer{
				FirstName: players.NewString("Fernando"),
				LastName:  players.NewString("Ocampo"),
				Nickname:  players.NewString("belsonnoles"),
				Email:     unittests.NewEmailAddress(t, "belsonnoles@anyemail.com"),
				Password:  players.NewString(passwordFixture()),
				Country:   players.NewString("Spain"),
			},
		},
		"invalid_email_update": {
			newPlayer: players.UpdatePlayer{
				Email: &mail.Address{},
			},
			isError: true,
			want:    "email is empty",
		},
		"invalid_empty_first_name_update": {
			newPlayer: players.UpdatePlayer{
				FirstName: players.NewString(""),
			},
			isError: true,
			want:    "first name is empty",
		},
		"invalid_empty_last_name_update": {
			newPlayer: players.UpdatePlayer{
				LastName: players.NewString(""),
			},
			isError: true,
			want:    "last name is empty",
		},
		"invalid_empty_nickname_update": {
			newPlayer: players.UpdatePlayer{
				Nickname: players.NewString(""),
			},
			isError: true,
			want:    "nickname is empty",
		},
		"invalid_empty_country_update": {
			newPlayer: players.UpdatePlayer{
				Country: players.NewString(""),
			},

			want:    "country is empty",
			isError: true,
		},
		"invalid_empty_password_update": {
			newPlayer: players.UpdatePlayer{
				Password: players.NewString(""),
			},

			want:    "password is empty",
			isError: true,
		},
		"invalid_empty_multiple_fields_update": {
			newPlayer: players.UpdatePlayer{
				Nickname: players.NewString(""),
				Password: players.NewString(""),
			},

			want: `nickname is empty
password is empty`,
			isError: true,
		},
	}

	for testName, testData := range testCases {
		t.Run(testName, func(t *testing.T) {
			// When
			err := testData.newPlayer.Validate()

			// Then
			if testData.isError {
				assert.Error(t, err)
				assert.Equal(t, testData.want, err.Error())
			}

			if !testData.isError {
				assert.NoError(t, err)
			}
		})
	}
}

func passwordFixture() string {
	return "A3BB605190830A01828F4D987A8C26CCBE3D4DAC0FAEF9482FBCB2B3CCB19CB8"
}
