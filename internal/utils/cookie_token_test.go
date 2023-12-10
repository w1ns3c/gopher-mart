package utils

import (
	"net/http"
	"testing"
	"time"
)

func TestCreateJWTcookie(t *testing.T) {
	type args struct {
		userid     string
		secret     string
		lifetime   time.Duration
		cookieName string
	}
	tests := []struct {
		name        string
		args        args
		cookieValue string
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name: "Test 1",
			args: args{
				userid:     "user1",
				secret:     "supersecret",
				lifetime:   time.Hour * 4,
				cookieName: "JWT",
			},
			wantErr: false,
		},
		{
			name: "Test 2 (expired cookie)",
			args: args{
				userid:     "22",
				secret:     "",
				lifetime:   time.Hour * -40,
				cookieName: "JWT",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCookie, err := CreateJWTcookie(tt.args.userid, tt.args.secret, tt.args.lifetime, tt.args.cookieName)
			if err != nil {
				t.Errorf("CreateJWTcookie() can't generate cookie, got error = %v", err)
				return
			}
			if gotCookie.Name != tt.args.cookieName {
				t.Errorf("CreateJWTcookie() got cookie.Name = %v, want %v",
					gotCookie.Name, tt.args.cookieName)
			}
			lifeTime := time.Now().Add(tt.args.lifetime - 5*time.Minute)
			if gotCookie.Expires.Before(lifeTime) {
				t.Errorf("CreateJWTcookie() lifetime = %v, want nearly %v", gotCookie.Expires, lifeTime)
			}

			userID, err := CheckJWTcookie(gotCookie, tt.args.secret)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("CreateJWTcookie() can't generate jwt cookie got = %v", err)
				}

			} else {
				if userID != tt.args.userid {
					t.Errorf("CreateJWTcookie() user in cookie got = %v, want = %v", userID, tt.args.userid)
				}
			}

		})
	}
}

func TestCheckJWTcookie(t *testing.T) {
	type args struct {
		cookie *http.Cookie
		secret string
	}
	tests := []struct {
		name       string
		args       args
		wantUserID string
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "Valid cookie",
			args: args{
				cookie: &http.Cookie{
					Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE0ODk1ODYxMTEsIlVzZXJJRCI6IjYzZjFjMzlmODY3YjI5NDA0NGQ0OGRjZTYwYzdkY2JjIn0.5_ibfJ6AwWnQO0qWzxYoDXLwfGKrLx59aLI48MCvhDU",
				},
				secret: "supersecret",
			},
			wantUserID: "63f1c39f867b294044d48dce60c7dcbc",
			wantErr:    false,
		},
		{
			name: "Expired cookie",
			args: args{
				cookie: &http.Cookie{
					Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE4NjExMSwiVXNlcklEIjoiMTIzMjEzMjExIn0.NGAPdTMreBv5V_mR7bEq1vEMsecpgQBopVEParQBe1g",
				},
				secret: "supersecret",
			},
			wantUserID: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := CheckJWTcookie(tt.args.cookie, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckJWTcookie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("CheckJWTcookie() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
				return
			}
		})
	}
}
