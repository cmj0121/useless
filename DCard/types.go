/* Copyright (C) 2020-2020 cmj. All right reserved. */
package dcard

import (
	"fmt"
)

type MediaMeta struct {
	Id            string   `json:"id"`
	Url           string   `json:"url"`
	NormalizedUrl string   `json:"normalizedUrl"`
	Thumbnail     string   `json:"thumbnail"`
	Type          string   `json:"type"`
	Tags          []string `json:"tags"`
	CreatedAt     string   `json:"createdAt"`
	UpdatedAt     string   `json:"updatedAt"`
	Width         int      `json:"width"`
	Height        int      `json:"height"`
}

type Media struct {
	Url string `json="url"`
}

type DCardPost struct {
	Id                  int64       `json:"id"`
	Title               string      `json:"title"`
	Excerpt             string      `json:"excerpt"`
	AnonymousSchool     bool        `json:"anonymousSchool"`
	AnonymousDepartment bool        `json:"anonymousDepartment"`
	Pinned              bool        `json:"pinned"`
	ForumId             string      `json:"forumId"`
	ReplyId             int64       `json:"replyId"`
	CreatedAt           string      `json:"createdAt"`
	UpdatedAt           string      `json:"updatedAt"`
	CommentCount        int         `json:"commentCount"`
	LikeCount           int         `json:"likeCount"`
	WithNickname        bool        `json:"withNickname"`
	Tags                []string    `json:"tags"`
	Topics              []string    `json:"topics"`
	ForumName           string      `json:"forumName"`
	ForumAlias          string      `json:"forumAlias"`
	Gender              string      `json:"gender"`
	ReplyTitle          string      `json:"replyTitle"`
	ReportReason        string      `json:"reportReason"`
	MediaMeta           []MediaMeta `json:"mediaMeta"`
	Hidden              bool        `json:"hidden"`
	CustomStyle         string      `json:"customStyle"`
	IsSuspiciousAccount bool        `json:"isSuspiciousAccount"`
	Layout              string      `json:"layout"`
	WithImages          bool        `json:"withImages"`
	WithVideos          bool        `json:"withVideos"`
	Media               []Media     `json:"media"`
	ReportReasonText    string      `json:"reportReasonText"`
	PostAvatar          string      `json:"postAvatar"`
}

func (p DCardPost) String() (out string) {
	out = fmt.Sprintf("[%d@%s] %s", p.Id, p.ForumName, p.Title)
	return
}

type DCardBoard struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

func (b DCardBoard) String() (out string) {
	out = fmt.Sprintf("%s (%s)", b.Alias, b.Name)
	return
}
