package main

import "strings"

//URLのパスの解析を自前で行う

const PathSeparator = "/"

type Path struct {
	Path string
	ID   string
}

// /people/1/をイメージ
func NewPath(p string) *Path {
	var id string
	//先頭末尾のスラッシュを削除
	p = strings.Trim(p, PathSeparator) // people/1
	//残りの文字列をスラッシュで区切って分割
	s := strings.Split(p, PathSeparator) // ["people",1]

	if len(s) > 1 {
		//スライスが複数ある場合は最後の項目をIDとする
		id = s[len(s)-1]                              // 1
		p = strings.Join(s[:len(s)-1], PathSeparator) // people/のみ。スライスは[{0}"people"{1}1{2}]なので全部含むには2が必要！
	}
	return &Path{Path: p, ID: id}
}

func (p *Path) HasID() bool {
	return len(p.ID) > 0
}
