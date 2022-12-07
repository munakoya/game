package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// 3*3の盤面
type Board [3][3]string

// 盤面などで使用する文字列
const maru, batsu = "〇", "✕"

// htmlに送信するデータ
type ViewData struct {
	Turn   string
	Board  *Board
	Win    bool   // 勝敗が付いた場合にtrueを設定
	Draw   bool   // 引き分けの場合にtrueを設定
	Winner string // 勝者を設定　"〇"か"✕"
}

// 変数nextTurnMapに、次の手番を取得するマップを設定
var nextTurnMap = map[string]string{
	maru:  batsu, // ○ の次は ×
	batsu: maru,  // × の次は ○
	"":    maru,  // 「""」の場合、ゲーム開始時として「"〇"」を取得
}

// テンプレートの設定
// template構造体のポインタを返す
var tpl *template.Template = template.Must(template.ParseFiles("index.html"))

// Executeメソッドの宣言
func (v *ViewData) Execute(w http.ResponseWriter) {
	// HTMLをクライアント（ブラウザ）に送信する
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.Execute(w, v)
	if err != nil {
		panic(err)
	}
}

// gameHandle関数の宣言 … ②
func gameHandle(w http.ResponseWriter, r *http.Request) {
	// 手番の入力値を取得する
	turn, nextTurn := turnFormValue(r)
	// 盤面の入力値を取得する
	board := boardFormValue(r)

	// 勝敗、引き分け、勝者の変数宣言と初期化
	win, draw, winner := false, false, ""

	// turnが「""」の場合、ゲーム開始時とする
	if turn != "" { // ゲーム開始時以外

		// winメソッド → Boardの入力値から勝利判定
		win = board.win(turn) // 勝敗の判定

		if win { // 勝敗が付いた場合（win の値が true のとき）
			winner = turn // 勝者を設定

			// 勝利が確定 → 空いたマスに-を設定
			// setBarメソッド
			board.setBar() // 未入力の項目に「"-"」を設定
		} else { // 勝敗が付かない場合（win の値が false のとき）
			draw = board.draw() // 引き分けの判定
		}
	}

	// 値を設定してHTMLを送信する
	v := ViewData{nextTurn, board, win, draw, winner}
	v.Execute(w)
}

// turnFormValue関数の宣言（手番の値を取得）
func turnFormValue(r *http.Request) (string, string) {
	// 現在の手番を取得
	turn := r.FormValue("turn")
	// マップを使用して次の手番を取得
	nextTurn := nextTurnMap[turn]
	return turn, nextTurn
}

// boardFormValue関数の宣言（盤面の値を取得）
func boardFormValue(r *http.Request) *Board {
	var board Board
	// Board型配列 board に保存されている要素とインデックス
	// を取り出す
	for i, array := range board {
		for j, _ := range array {
			// テンプレート上のボードの name 属性値を作成
			name := "c" + strconv.Itoa(i) + "," + strconv.Itoa(j)
			// 盤面の各項目を取得
			board[i][j] = r.FormValue(name)
			// デバッグ用
			fmt.Printf("%v = %v\n", name, board[i][j])

		}
	}

	// 構造体 board の参照を返す
	return &board
}

func main() {
	//ローカルサーバ立てる

	// ハンドラーcalcの設定
	http.HandleFunc("/game", gameHandle)

	// http.ListenAndServe関数 → httpサーバー起動
	// http.ListenAndServe(サーバアドレス, ルーティングハンドラ)

	// ルーティングハンドラが設定されているのならサーバ起動？
	result := http.ListenAndServe(":8888", nil)
	if result != nil {
		fmt.Println(result)
	}
}

// winメソッドの宣言（勝敗の判定）
func (b *Board) win(turn string) bool {
	// 横 上
	if b[0][0] == turn && b[0][1] == turn && b[0][2] == turn {
		return true
	}
	// 横 真ん中
	if b[1][0] == turn && b[1][1] == turn && b[1][2] == turn {
		return true
	}
	// 横 下
	if b[2][0] == turn && b[2][1] == turn && b[2][2] == turn {
		return true
	}
	// 縦 左
	if b[0][0] == turn && b[1][0] == turn && b[2][0] == turn {
		return true
	}
	// 縦 真ん中
	if b[0][1] == turn && b[1][1] == turn && b[2][1] == turn {
		return true
	}
	// 縦 右
	if b[0][2] == turn && b[1][2] == turn && b[2][2] == turn {
		return true
	}
	// 斜め1
	if b[0][0] == turn && b[1][1] == turn && b[2][2] == turn {
		return true
	}
	// 斜め2
	if b[2][0] == turn && b[1][1] == turn && b[0][2] == turn {
		return true
	}
	return false
}

// drawメソッドの宣言（引き分けの判定）
func (b *Board) draw() bool {
	for i, array := range b {
		for j, _ := range array {
			// 空白がある → どっちか勝ちorゲーム続行
			// winメソッドがfalseの際に実行される
			if b[i][j] == "" {
				return false // 未入力の項目がある場合、ゲームを続行}
			}
		}
	}
	return true // 未入力の項目がない場合、引き分け
}

// setBarメソッドの宣言（「"-"」の設定）
func (b *Board) setBar() {
	for i, array := range b {
		for j, _ := range array {
			// 二次元配列の空白をループで調べる
			if b[i][j] == "" {
				b[i][j] = "-" // 未入力の項目は「"-"」を設定
			}
		}

	}
}
