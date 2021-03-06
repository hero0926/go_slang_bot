package util

import (
	"context"
	"log"
	"time"

	"slackbot/envsetting"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	Pushevent        = "PushEvent"
	PullRequestEvent = "PullRequestEvent"
)

// 유저의 깃허브 커밋 여부를 받아오기
func GetGitCommit(id string) (b bool, commit int) {

	// 현재시간 구하기
	koryear, kormon, kordate := time.Now().Date()
	loc, _ := time.LoadLocation("Asia/Seoul")

	var commitArray []string

	// github 연결
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: envsetting.Gitconfig()},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// 유저 이벤트 받기
	events, _, err := client.Activity.ListEventsPerformedByUser(context.Background(), id, true, nil)
	if err != nil {
		log.Println("[ERROR]", err)
		return false, 1
	}

	for _, v := range events {

		utctime := v.GetCreatedAt()
		year, month, day := utctime.In(loc).Date()

		// 오늘 한 커밋만 고르기
		if ((v.GetType() == Pushevent) || (v.GetType() == PullRequestEvent)) && ((koryear == year) && (kormon == month) && (kordate == day)) {
			commitID := v.GetID()
			commitArray = append(commitArray, commitID)
		}
	}

	// 커밋이 없으면 false, 있으면 true
	if len(commitArray) == 0 {
		return false, 0
	} else {

		return true, len(commitArray)
	}

}
