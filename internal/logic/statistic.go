package logic

import (
	"Testovoe1/internal/entity"
	"context"
	"strings"
	"time"
)

type Logic struct {
}

func CreateStatistics(ctx context.Context, comments []*entity.Comment) []*entity.Statistic {
	statistics := []*entity.Statistic{}
	for _, comment := range comments {
		words := strings.Fields(comment.Body)
		for _, word := range words {
			if len(statistics) == 0 {
				statistics = append(statistics, &entity.Statistic{
					Postid: comment.PostId,
					Word:   word,
					Count:  1,
					Time:   time.Now(),
				})
			}
			flag := false
			for _, statistic := range statistics {
				if statistic.Word == word && statistic.Postid == comment.PostId {
					statistic.Count++
					flag = true
				}
			}
			if !flag {
				statistics = append(statistics, &entity.Statistic{
					Postid: comment.PostId,
					Word:   word,
					Count:  1,
					Time:   time.Now(),
				})
			}
		}
	}
	return statistics
}
