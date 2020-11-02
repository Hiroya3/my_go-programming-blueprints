package backup

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

type Archiver interface {
	//src：バックアップ対象のパス
	//dest：保存先のパス
	Archive(src, dest string) error
}

type zipper struct{}

//ZIPはファイルの圧縮とその解除にZIP形式を利用するArchiverです
var ZIP Archiver = (*zipper)(nil)

func (z *zipper) Archive(src, dest string) error {
	//保存先のディレクトリが存在するか確認。なければ0777の権限で作成。
	if err := os.MkdirAll(filepath.Dir(dest), 0777); err != nil {
		return err
	}
	//保存先のディレクトリにファイルを新規作成
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	w := zip.NewWriter(out)
	defer w.Close()
	//srcのファイルパスの全てのファイルについてfuncの関数が行われる
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil //フォルダなのでスキップします
		}
		if err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		f, err := w.Create(path)
		if err != nil {
			return err
		}
		io.Copy(f, in)
		return nil
	})
}
