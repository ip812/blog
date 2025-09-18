package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/godruoyi/go-snowflake"
)

const (
	CookieKey        = "IP812_BLOG_USERNAME"
	DefaultAvatarURL = "https://avatars.githubusercontent.com/u/2878733?v=4"
)

func generateUsername() string {
	return "user_" + strconv.FormatUint(snowflake.ID(), 10)
}

func getAvatarURL(username string) string {
	parts := strings.Split(username, "_")

	if len(parts) != 2 {
		return DefaultAvatarURL
	}

	id, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return DefaultAvatarURL
	}

	return fmt.Sprintf("https://robohash.org/%d?set=set4", id)
}
