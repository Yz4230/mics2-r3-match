package main

import "fmt"

type PieceType int8

// 55将棋
const (
	NONE    PieceType = iota // 空
	KING                     // 玉
	HISHA                    // 飛
	KAKU                     // 角
	KIN                      // 金
	GIN                      // 銀
	HU                       // 歩
	RYU                      // 龍(成飛)
	UMA                      // 馬(成角)
	NARIGIN                  // 成銀
	TOKIN                    // と金
)

func (p PieceType) ToPromoted() PieceType {
	switch p {
	case HISHA:
		return RYU
	case KAKU:
		return UMA
	case GIN:
		return NARIGIN
	case HU:
		return TOKIN
	}
	panic("invalid piece type")
}

func (p PieceType) String() string {
	switch p {
	case NONE:
		return "・"
	case KING:
		return "玉"
	case HISHA:
		return "飛"
	case KAKU:
		return "角"
	case KIN:
		return "金"
	case GIN:
		return "銀"
	case HU:
		return "歩"
	case RYU:
		return "龍"
	case UMA:
		return "馬"
	case NARIGIN:
		return "全"
	case TOKIN:
		return "と"
	}
	panic("invalid piece type")
}

type Position struct {
	// 先手の盤面
	FirstBoard [5][5]PieceType
	// 後手の盤面
	SecondBoard [5][5]PieceType
	// 先手の持ち駒
	FirstHand []PieceType
	// 後手の持ち駒
	SecondHand []PieceType
	// 手番
	Side Side
	// 手数
	Count int
}

func (p Position) String() string {
	s := ""
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			fp := p.FirstBoard[y][x]
			sp := p.SecondBoard[y][x]
			if fp != NONE {
				s += firstStyle.Render(fmt.Sprintf("▲%s", fp))
			} else if sp != NONE {
				s += secondStyle.Render(fmt.Sprintf("▼%s", sp))
			} else {
				s += " ・"
			}
		}
		s += "\n"
	}
	s += "先手の持ち駒: "
	for _, piece := range p.FirstHand {
		s += firstStyle.Render(piece.String())
	}
	s += "\n"
	s += "後手の持ち駒: "
	for _, piece := range p.SecondHand {
		s += secondStyle.Render(piece.String())
	}
	s += "\n"
	s += fmt.Sprintf("手番: %s\n", p.Side)
	s += fmt.Sprintf("手数: %d", p.Count)
	return s
}

type Side int8

const (
	FIRST  Side = 0
	SECOND Side = 1
)

func (s Side) String() string {
	switch s {
	case FIRST:
		return firstStyle.Render("先手")
	case SECOND:
		return secondStyle.Render("後手")
	}
	panic("invalid side")
}

func (s Side) Log(str string) {
	if args.Verbose {
		fmt.Printf("%s: %s\n", s, str)
	}
}
