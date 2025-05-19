package JWTParser

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestParseUnverified(t *testing.T) {
	tests := []struct {
		TestName      string
		JWToken       string
		isJWTValid    bool
		errorExpected bool
		errorContains string
	}{
		{
			TestName:      "PositiveJWTParse",
			JWToken:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2NvdW50SWQiOiI0NDUiLCJDb21wYW55SWQiOiIyNzYiLCJFbWFpbCI6IjFwYXZlbC5zaGx5a292QGl0bHRlYW0udGVzdCIsIlBob25lIjoiKzcgODIxIDczOC0wMS02MCIsIkNvbXBhbnlDb2RlIjoiIiwiQ29tcGFueU5hbWUiOiJUcmFuc2Zvcm1lcnMgY29ycG9yYXRpb24iLCJDb21wYW55TG9jYWxlIjoicnUiLCJSb2xlIjoiT3duZXIiLCJJc0JvdCI6IkZhbHNlIiwibmJmIjoxNzQ3NTUwMzQ3LCJleHAiOjE3NDc1NTA2NDcsImlzcyI6IlVudWtTZXJ2ZXIiLCJhdWQiOiJVbnVrQ2xpZW50In0.TKmfSVV9WQyW1hQjnsXSnsZRM4AwoYbSonqdsezvOWk",
			isJWTValid:    true,
			errorExpected: false,
			errorContains: "",
		},
		{
			TestName:      "NegativeJWTParse",
			JWToken:       "invalid",
			isJWTValid:    false,
			errorExpected: true,
			errorContains: "token contains an invalid number of segments",
		},
	}
	var logOutput bytes.Buffer
	log := slog.New(slog.NewTextHandler(&logOutput, nil))
	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			jwtClaims, err := ParseUnverified(tt.JWToken, log)
			if tt.errorExpected {
				if err == nil {
					t.Errorf("ParseUnverified(%s) should have failed, but there is no error", tt.JWToken)
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Error not contains expected text %v", tt.errorContains)
				}
				return
			}
			if err != nil {
				t.Error("ParseUnverified failed: " + err.Error())
			}
			if jwtClaims == nil {
				t.Error("ParseUnverified() returned nil")
			}

		})
	}
}

func TestGetExpirationTime(t *testing.T) {
	mskLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		t.Fatalf("Failed to load MSK location: %v", err)
	}
	tests := []struct {
		TestName               string
		JWToken                string
		isJWTValid             bool
		expectedExpirationTime time.Time
		errorExpected          bool
		errorContains          string
	}{
		{
			TestName:               "PositiveGetExpirationTime",
			JWToken:                "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2NvdW50SWQiOiI0NDUiLCJDb21wYW55SWQiOiIyNzYiLCJFbWFpbCI6IjFwYXZlbC5zaGx5a292QGl0bHRlYW0udGVzdCIsIlBob25lIjoiKzcgODIxIDczOC0wMS02MCIsIkNvbXBhbnlDb2RlIjoiIiwiQ29tcGFueU5hbWUiOiJUcmFuc2Zvcm1lcnMgY29ycG9yYXRpb24iLCJDb21wYW55TG9jYWxlIjoicnUiLCJSb2xlIjoiT3duZXIiLCJJc0JvdCI6IkZhbHNlIiwibmJmIjoxNzQ3NTUwMzQ3LCJleHAiOjE3NDc1NTA2NDcsImlzcyI6IlVudWtTZXJ2ZXIiLCJhdWQiOiJVbnVrQ2xpZW50In0.TKmfSVV9WQyW1hQjnsXSnsZRM4AwoYbSonqdsezvOWk",
			isJWTValid:             true,
			expectedExpirationTime: time.Date(2025, 5, 18, 9, 44, 7, 0, mskLocation),
			errorExpected:          false,
			errorContains:          "",
		},
		{
			TestName:               "NegativeGetExpirationTime",
			JWToken:                "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2NvdW50SWQiOiI0NDUiLCJDb21wYW55SWQiOiIyNzYiLCJDb21wYW55Q29kZSI6IiIsIkNvbXBhbnlMb2NhbGUiOiJydSIsIkNvbXBhbnlOYW1lIjoiVHJhbnNmb3JtZXJzIGNvcnBvcmF0aW9uIiwiRW1haWwiOiIxcGF2ZWwuc2hseWtvdkBpdGx0ZWFtLnRlc3QiLCJJc0JvdCI6IkZhbHNlIiwgIlBob25lIjoiKzcgODIxIDczOC0wMS02MCIsIlJvbGUiOiJPd25lciIsImF1ZCI6IlVudWtDbGllbnQiLCJleHAiOjAsImV4cCI6MCwiaXNzIjoiVW51a1NlcnZlciIsIm5iZiI6MTY1NDA2ODAwMH0.5AhdZ6W5Q1n8q7X2Y9vVbVcVdVeVfVgVhViVjVkVlVm",
			isJWTValid:             false,
			expectedExpirationTime: time.Date(2025, 5, 18, 9, 44, 7, 0, mskLocation),
			errorExpected:          true,
			errorContains:          "token has no expiration claim",
		},
	}
	var logOutput bytes.Buffer
	log := slog.New(slog.NewTextHandler(&logOutput, nil))
	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			jwtClaims, err := ParseUnverified(tt.JWToken, log)
			if err != nil {
				t.Error("ParseUnverified failed: " + err.Error())
			}
			expTime, err := GetExpirationTime(jwtClaims, log)
			if tt.errorExpected {
				if err == nil {
					t.Errorf("GetExpirationTime(%s) should have failed, but there is no error", tt.JWToken)
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Error not contains expected text %v", tt.errorContains)
				}
				return
			}
			if err != nil {
				t.Error("GetExpirationTime failed: " + err.Error())
			}
			if jwtClaims == nil {
				t.Error("GetExpirationTime returned nil")
			}
			if !expTime.Equal(tt.expectedExpirationTime) {
				t.Errorf("GetExpirationTime returned %v, want %v", expTime, tt.expectedExpirationTime)
			}
		})
	}
}
