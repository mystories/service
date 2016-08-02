package service

import (
	"log"
	"context"
	"model"
	"time"
	"fmt"
)

type Article struct {
	Ctx		*context.Context
}

type ArticleGetArgs struct {
	Id		int		`msgpack:"id"`
}

type ArticleType struct {
	CreatedAt	int64
	UpdateAt	int64
	Title		string
	Source		string
	Content		string
}

type ArticleGetResult struct {
	Error		int
	Msg			string
	Article		ArticleType
}

type ArticleListArgs struct {
	Maxid		int		`msgpack:"maxid"`
	Count		int		`msgpack:"count"`
	NeedContent	bool	`msgpack:"need_content"`
}

type ArticleListResult struct {
	Error		int
	Msg			string
	Articles	[]ArticleType
}

func (a *Article) Get(args *ArticleGetArgs, reply *ArticleGetResult) error {
	if args.Id == 0 {
		(*reply).Error = 1
		(*reply).Msg = "invalid id"
		log.Printf("Article.Get invalid params, id=%d\n", args.Id)
		return nil
	}
	log.Printf("Article.Get id=%d\n", args.Id)
	rows, err := a.Ctx.Db.Query(fmt.Sprintf("SELECT * FROM articles where id=%d", args.Id))
	if err != nil {
		(*reply).Error = 2
		(*reply).Msg = "failed to get article"
		log.Printf("failed to get article, err=%v\n", err)
		return nil
	}
	var article model.Articles
	found := false
	for rows.Next() {
		err := rows.Scan(
			&article.Id,
			&article.CreatedAt,
			&article.UpdateAt,
			&article.Title,
			&article.Source,
			&article.Content,
		)
		if err != nil {
			(*reply).Error = 2
			(*reply).Msg = "failed to get article"
			log.Printf("failed to scan article, err=%v\n", err)
			return nil
		}
		found = true
		break
	}
	if !found {
		(*reply).Error = 3
		(*reply).Msg = "article not found"
		log.Println("article not found")
		return nil
	}
	createdAt, _ := time.Parse("2006-01-02 15:04:05", article.CreatedAt)
	updateAt, _ := time.Parse("2006-01-02 15:04:05", article.UpdateAt)
	(*reply).Article = ArticleType{
		CreatedAt: createdAt.Unix(),
		UpdateAt: updateAt.Unix(),
		Title: article.Title,
		Source: article.Source,
		Content: article.Content,
	}
	return nil
}

func (a *Article) List(args *ArticleListArgs, reply *ArticleListResult) error {
	log.Printf("Article.List maxid=%d count=%d\n", args.Maxid, args.Count)
	sql := "SELECT * FROM articles"
	if args.Maxid != 0 {
		sql += fmt.Sprintf(" where id<%d", args.Maxid)
	}
	if args.Count != 0 {
		sql += fmt.Sprintf(" limit %d", args.Count)
	}
	log.Printf("sql=%s\n", sql)
	rows, err := a.Ctx.Db.Query(sql)
	if err != nil {
		(*reply).Error = 2
		(*reply).Msg = "failed to get article"
		log.Printf("failed to get article, err=%v\n", err)
		return nil
	}
	var article model.Articles
	(*reply).Articles = []ArticleType{}
	for rows.Next() {
		err := rows.Scan(
			&article.Id,
			&article.CreatedAt,
			&article.UpdateAt,
			&article.Title,
			&article.Source,
			&article.Content,
		)
		if err != nil {
			(*reply).Error = 2
			(*reply).Msg = "failed to get articles"
			log.Printf("failed to scan article, err=%v\n", err)
			return nil
		}
		createdAt, _ := time.Parse("2006-01-02 15:04:05", article.CreatedAt)
		updateAt, _ := time.Parse("2006-01-02 15:04:05", article.UpdateAt)
		(*reply).Articles = append((*reply).Articles, ArticleType{
			CreatedAt: createdAt.Unix(),
			UpdateAt: updateAt.Unix(),
			Title: article.Title,
			Source: article.Source,
			Content: article.Content,
		})
	}

	return nil
}

