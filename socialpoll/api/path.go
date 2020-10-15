package main

import "strings"

//URLのパスの解析を自前で行う

const PathSeparator = "/"

type Path struct {
	Path string
	ID   string
}

func NewPath(p string) *Path {
	var id string
	//先頭末尾のスラッシュを削除
	p = strings.Trim(p, PathSeparator)
	//残りの文字列をスラッシュで区切って分割
	s := strings.Split(p, PathSeparator)

	if len(s) > 1 {
		//スライスが複数ある場合は最後の項目をIDとする
		id = s[len(s)-1]
		p = strings.Join(s[:len(s)-1], PathSeparator)
	}
	return &Path{Path: p, ID: id}
}

func (p *Path) HasID() bool {
	return len(p.ID) > 0
}
