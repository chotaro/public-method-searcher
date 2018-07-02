package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	flag.Parse()
	targetDirectoryName := flag.Arg(0)
	println(targetDirectoryName)
	//対象ディレクトリ配下の一覧を作成し、javaファイルのみ残す
	targets := dirwalk(targetDirectoryName)
	matchedTagets, error := matchedStrings(targets, "java")
	if error != nil {
		// Openエラー処理
		fmt.Println(error)
	}

	//各ファイルについて、抽出とresultへの書き出しを行う。
	for _, filename := range matchedTagets {
		// println(strconv.Itoa(counter) + ":")
		analize(filename)
	}
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}

//matchedStrings string配列から、regexにマッチしたもののみを取り出す。
func matchedStrings(slice []string, regex string) ([]string, error) {
	ret := make([]string, len(slice))
	i := 0
	for _, targetString := range slice {
		if isMatchRegex(targetString, regex) {
			ret[i] = targetString
			i++
		}
	}
	// if len(ret[:i]) == len(slice) {
	// return slice, fmt.Errorf("Couldn't find" + ret[i])
	// }
	return ret[:i], nil
}

// 正規表現でマッチするかどうか判断を返す
func isMatchRegex(target string, regex string) bool {
	condition := regexp.MustCompile(regex)
	return condition.MatchString(target)
}

//analize 引数で渡されたファイルを解析する。publicに当たればそれを書き込んで追記する。
func analize(filepath string) {
	println("[ filepath: " + filepath + " ]")
	file, error := os.Open(filepath)
	if error != nil {
		// Openエラー処理
		fmt.Print("error occured\n")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// 一行ずつ読み込んで文字列の中身を判定し、ヒットしたものについてresultへ書き出しを行う。
	for scanner.Scan() {
		readLine := scanner.Text()
		//正規表現でpublicメソッドを抜き出す
		if isMatchRegex(readLine, "public.*\\(.*\\{") {
			fmt.Println("    :" + readLine)
		}
	}
	println("[ file end ]")
	println("")
}
