package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
)

var cards = [55]string{
	"Ac", "As", "Ah", "Ad",
	"2c", "2s", "2h", "2d",
	"3c", "3s", "3h", "3d",
	"4c", "4s", "4h", "4d",
	"5c", "5s", "5h", "5d",
	"6c", "6s", "6h", "6d",
	"7c", "7s", "7h", "7d",
	"8c", "8s", "8h", "8d",
	"9c", "9s", "9h", "9d",
	"Tc", "Ts", "Th", "Td",
	"Jc", "Js", "Jh", "Jd",
	"Qc", "Qs", "Qh", "Qd",
	"Kc", "Ks", "Kh", "Kd",
}

type player struct {
	in_play  bool
	fold     bool
	name     string
	chips    int
	curr_bet int
	card     [2]string
}

type table struct {
	card       [5]string
	pot_number int
	pot        [8]int
}

type coord struct {
	x int
	y int
}

var coords = [8]coord{
	{480, 350},
	{480, 630},
	{385, 80},
	{270, 50},
	{385, 915},
	{270, 960},
}

var board = table{}
var use = [55]bool{}
var deck = [55]string{}
var players = [8]player{}
var curr, act = 0, 0

var active = [8]bool{}
var tr_cnt, pl_cnt, wait = 0, 0, false

var t = template.Must(template.ParseGlob("template/*.html"))

func deal() {
	for i := 0; i < 8; i++ {
		if players[i].in_play == false {
			continue
		}
		players[i].card[0] = deck[curr]
		curr++
	}
	for i := 0; i < 8; i++ {
		if players[i].in_play == false {
			continue
		}
		players[i].card[1] = deck[curr]
		curr++
	}
}

func flop() {
	curr++
	board.card[0] = deck[curr]
	curr++
	board.card[1] = deck[curr]
	curr++
	board.card[2] = deck[curr]
}

func river() {
	curr += 2
	board.card[3] = deck[curr]
}

func turn() {
	curr += 2
	board.card[4] = deck[curr]
}
func shuffle() {
	for i := 0; i <= 52; i++ {
		use[52] = true
		var n = rand.Int63n(52)
		var j = n
		for ; j <= 52 && use[j] == true; j++ {
			if j == 52 {
				use[52] = false
				j = 0
			}
		}
		deck[i] = cards[j]
		use[j] = true
	}
}

func check() {
	act++
}

func fold() {
	players[act].fold = true
	act++
}

func call() {
	players[act].chips -= players[act-1].curr_bet - players[act].curr_bet
	players[act].curr_bet = players[act-1].curr_bet
	act++
}

func raise(bet int) {
	players[act].chips -= bet
	players[act].curr_bet += bet
}

func loadHTML(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "index.html", nil)
}

func add_player(w http.ResponseWriter, r *http.Request) {
	pl_cnt++
	wait = true
	for i := 0; i < 8; i++ {
		if players[i].in_play == true {
			continue
		}
		players[i].in_play = true
		players[i].name = r.PostFormValue("name")
		var Namestr, Butstr = "<div>", ""
		for j := 0; j < 8; j++ {
			if players[j].in_play {
				Namestr += "<div><h3 class='player' style='position:absolute; top:" + strconv.Itoa(coords[j].x) + "px; left:" + strconv.Itoa(coords[j].y) + "px;'>%s</h3>\n"
				Namestr += "<input class='invisible' hx-post='/check' hx-target='.back-end'  hx-swap='outterHTML' hx-trigger='every 200ms' type='text' name='id' class='input-text' value=" + strconv.Itoa(i) + "></div>\n"

				Butstr += "<button class='crf' class='check'>Check</button>"
				Butstr += "<button class='crf' class='call'>Call</button>"
				Butstr += "<button class='crf' class='raise'>Raise</button>"
			} else {
				Namestr += "%s\n"
			}
		}
		Namestr += "</div>"
		htmlStr := fmt.Sprintf(Namestr+Butstr, players[0].name, players[1].name, players[2].name, players[3].name, players[4].name, players[5].name, players[6].name, players[7].name)
		tmpl, _ := template.New("t").Parse(htmlStr)
		tmpl.Execute(w, nil)
		//fmt.Println(str)
		break
	}

}

func update(w http.ResponseWriter, r *http.Request) {
	tr_cnt++
	fmt.Println("hey, can i play", r.PostFormValue("id"))
	var i = r.PostFormValue("id")
	var id, _ = strconv.Atoi(i)
	active[id] = true
	fmt.Println("a ", wait, tr_cnt, pl_cnt)
	if tr_cnt >= pl_cnt+1 {
		tr_cnt = 0
		wait = false
		for i := 0; i < 8; i++ {
			if active[i] {
				active[i] = false
				continue
			}
			if players[i].in_play {
				players[i].in_play = false
				players[i].name = ""
				pl_cnt--
			}
		}
	}
	//fmt.Println(tr_cnt, pl_cnt)

	str, Butstr := "<div>", ""
	for j := 0; j < 8; j++ {
		if players[j].in_play {
			str += "<div><h3 class='player' style='position:absolute; top:" + strconv.Itoa(coords[j].x) + "px; left:" + strconv.Itoa(coords[j].y) + "px;'>%s</h3>\n"
			if j == id {
				str += "<input class='invisible' hx-post='/check' hx-target='.back-end'  hx-swap='outterHTML' hx-trigger='every 200ms' type='text' name='id' class='input-text' value='" + i + "'></div>\n"
				Butstr += "<button class='crf'>Check</button>"
				Butstr += "<button class='crf' class='call'>Call</button>"
				Butstr += "<button class='crf' class='raise'>Raise</button>"
			}
		} else {
			str += "%s\n"
		}
	}
	str += "</div>"
	str += Butstr
	fmt.Println(str)
	htmlStr := fmt.Sprintf(str, players[0].name, players[1].name, players[2].name, players[3].name, players[4].name, players[5].name, players[6].name, players[7].name)
	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(w, nil)
}

func main() {
	/*
		players[0].in_play = true
		players[0].name = "Vatman"

		players[1].in_play = true
		players[1].name = "Albert"

		players[2].in_play = true
		players[2].name = "Darius"

		players[3].in_play = true
		players[3].name = "Batei"

		players[4].in_play = true
		players[4].name = "Faso"
	*/
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))

	http.HandleFunc("/", loadHTML)
	http.HandleFunc("/add-player/", add_player)
	http.HandleFunc("/check", update)
	shuffle()

	deal()
	flop()
	river()
	turn()
	http.ListenAndServe(":8080", nil)

}
