package main

import (
	"strconv"
	"strings"
)

func parsePiece(piece string) (PieceType, Side) {
	switch piece {
	case "R":
		return HISHA, FIRST
	case "B":
		return KAKU, FIRST
	case "G":
		return KIN, FIRST
	case "S":
		return GIN, FIRST
	case "K":
		return KING, FIRST
	case "P":
		return HU, FIRST
	case "r":
		return HISHA, SECOND
	case "b":
		return KAKU, SECOND
	case "g":
		return KIN, SECOND
	case "s":
		return GIN, SECOND
	case "k":
		return KING, SECOND
	case "p":
		return HU, SECOND
	}
	panic("invalid piece type")
}

func parseSfen(sfen string) Position {
	// sfen format: rbsgk/4p/5/P4/KGSBR b - 1

	sfen = strings.TrimSpace(sfen)

	firstBoard := [5][5]PieceType{}
	secondBoard := [5][5]PieceType{}
	firstRow := 0
	firstCol := 0
	secondRow := 0
	secondCol := 0
	cur := 0

	// 1. 盤上の駒を読み取る
	for cur < len(sfen) {
		c := string(sfen[cur])
		if i, err := strconv.Atoi(c); err == nil {
			// is number
			firstCol += i
			secondCol += i
			cur++
			continue
		} else if sfen[cur] == '/' {
			firstRow++
			firstCol = 0
			secondRow++
			secondCol = 0
			cur++
			continue
		} else if sfen[cur] == ' ' {
			cur++
			break
		} else {
			isPromoted := c == "+"
			if isPromoted {
				cur++
				c = string(sfen[cur])
			}
			piece, side := parsePiece(c)
			if isPromoted {
				piece = piece.ToPromoted()
			}
			if side == FIRST {
				firstBoard[firstRow][firstCol] = piece
			} else {
				secondBoard[secondRow][secondCol] = piece
			}
			firstCol++
			secondCol++
			cur++
		}
	}

	// 2. 手番を読み取る
	side := FIRST
	if sfen[cur] == 'w' {
		side = SECOND
	}
	cur += 2

	// 3. 持ち駒を読み取る
	firstHand := []PieceType{}
	secondHand := []PieceType{}
	if sfen[cur] == '-' {
		cur += 2
	} else {
		for cur < len(sfen) {
			c := string(sfen[cur])
			if c == " " {
				cur++
				break
			}
			n := 1
			if num, err := strconv.Atoi(c); err == nil {
				n = num
				cur++
				c = string(sfen[cur])
			}

			piece, side := parsePiece(c)
			for j := 0; j < n; j++ {
				if side == FIRST {
					firstHand = append(firstHand, piece)
				} else {
					secondHand = append(secondHand, piece)
				}
			}
			cur++
		}
	}

	// 4. 手数を読み取る
	count, err := strconv.Atoi(sfen[cur:])
	if err != nil {
		panic(err)
	}

	return Position{
		FirstBoard:  firstBoard,
		SecondBoard: secondBoard,
		FirstHand:   firstHand,
		SecondHand:  secondHand,
		Side:        side,
		Count:       count,
	}
}
