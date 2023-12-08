package main

import (
	"reflect"
	"testing"
)

func Test_parseSfen(t *testing.T) {
	type args struct {
		sfen string
	}
	tests := []struct {
		name string
		args args
		want Position
	}{
		{
			name: "rbsgk/4p/5/P4/KGSBR b - 1",
			args: args{
				sfen: "rbsgk/4p/5/P4/KGSBR b - 1",
			},
			want: Position{
				FirstBoard: [5][5]PieceType{
					{NONE, NONE, NONE, NONE, NONE},
					{NONE, NONE, NONE, NONE, NONE},
					{NONE, NONE, NONE, NONE, NONE},
					{HU, NONE, NONE, NONE, NONE},
					{KING, KIN, GIN, KAKU, HISHA},
				},
				SecondBoard: [5][5]PieceType{
					{HISHA, KAKU, GIN, KIN, KING},
					{NONE, NONE, NONE, NONE, HU},
					{NONE, NONE, NONE, NONE, NONE},
					{NONE, NONE, NONE, NONE, NONE},
					{NONE, NONE, NONE, NONE, NONE},
				},
				FirstHand:  []PieceType{},
				SecondHand: []PieceType{},
				Side:       FIRST,
				Count:      1,
			},
		},
		{
			name: "r3b/1k1p1/s1g1b/P4/KGSBR w - 32",
			args: args{
				sfen: "r3b/1k1p1/s3g/P4/KGSBR w - 32",
			},
			want: Position{
				FirstBoard: [5][5]PieceType{
					{NONE, NONE, NONE, NONE, NONE},
					{NONE, NONE, NONE, NONE, NONE},
					{NONE, NONE, NONE, NONE, NONE},
					{HU, NONE, NONE, NONE, NONE},
					{KING, KIN, GIN, KAKU, HISHA},
				},
				SecondBoard: [5][5]PieceType{
					{HISHA, NONE, NONE, NONE, KAKU},
					{NONE, KING, NONE, HU, NONE},
					{GIN, NONE, NONE, NONE, KIN},
					{NONE, NONE, NONE, NONE, NONE},
					{NONE, NONE, NONE, NONE, NONE},
				},
				FirstHand:  []PieceType{},
				SecondHand: []PieceType{},
				Side:       SECOND,
				Count:      32,
			},
		},
		{
			// 先手側が銀１枚歩２枚、後手側が角１枚歩３枚
			name: "rbsgk/4p/5/P4/KGSBR w S2Pb3p 32",
			args: args{
				sfen: "rbsgk/4p/5/P4/KGSBR w S2Pb3p 32",
			},
			want: Position{
				FirstBoard: [5][5]PieceType{
					{NONE, NONE, NONE, NONE, NONE},
					{NONE, NONE, NONE, NONE, NONE},
					{NONE, NONE, NONE, NONE, NONE},
					{HU, NONE, NONE, NONE, NONE},
					{KING, KIN, GIN, KAKU, HISHA},
				},
				SecondBoard: [5][5]PieceType{
					{HISHA, KAKU, GIN, KIN, KING},
					{NONE, NONE, NONE, NONE, HU},
					{NONE, NONE, NONE, NONE, NONE},
					{NONE, NONE, NONE, NONE, NONE},
					{NONE, NONE, NONE, NONE, NONE},
				},
				FirstHand:  []PieceType{GIN, HU, HU},
				SecondHand: []PieceType{KAKU, HU, HU, HU},
				Side:       SECOND,
				Count:      32,
			},
		},
		{
			// 先手側が銀１枚歩２枚、後手側が角１枚歩３枚
			name: "+r+b+sgk/4+p/5/P4/KGSBR w S2Pb3p 32",
			args: args{
				sfen: "+r+b+sgk/4+p/5/P4/KGSBR w S2Pb3p 32",
			},
			want: Position{
				FirstBoard: [5][5]PieceType{
					{NONE, NONE, NONE, NONE, NONE},
					{NONE, NONE, NONE, NONE, NONE},
					{NONE, NONE, NONE, NONE, NONE},
					{HU, NONE, NONE, NONE, NONE},
					{KING, KIN, GIN, KAKU, HISHA},
				},
				SecondBoard: [5][5]PieceType{
					{RYU, UMA, NARIGIN, KIN, KING},
					{NONE, NONE, NONE, NONE, TOKIN},
					{NONE, NONE, NONE, NONE, NONE},
					{NONE, NONE, NONE, NONE, NONE},
					{NONE, NONE, NONE, NONE, NONE},
				},
				FirstHand:  []PieceType{GIN, HU, HU},
				SecondHand: []PieceType{KAKU, HU, HU, HU},
				Side:       SECOND,
				Count:      32,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseSfen(tt.args.sfen); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSfen() = \n%v\nwant = \n%v", got, tt.want)
			}
		})
	}
}
