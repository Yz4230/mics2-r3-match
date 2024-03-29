package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

type Result struct {
	moves  []string
	sfens  []string
	winner Side
}

type Matcher struct {
	// first player program
	Fp string
	// first player random
	Fr int
	// first player depth
	Fd int
	// second player program
	Sp string
	// second player random
	Sr int
	// second player depth
	Sd int
	// byoyomi
	Byoyomi int
}

func (m *Matcher) Match() (*Result, error) {
	firstCmd := exec.Command(m.Fp, "-r", fmt.Sprint(m.Fr), "-d", fmt.Sprint(m.Fd))
	secondCmd := exec.Command(m.Sp, "-r", fmt.Sprint(m.Sr), "-d", fmt.Sprint(m.Sd))

	FIRST.Log(firstCmd.String())
	SECOND.Log(secondCmd.String())

	firstCmdStdin, _ := firstCmd.StdinPipe()
	firstCmdStdout, _ := firstCmd.StdoutPipe()
	firstStdinReader := bufio.NewReader(firstCmdStdout)

	secondCmdStdin, _ := secondCmd.StdinPipe()
	secondCmdStdout, _ := secondCmd.StdoutPipe()
	secondStdinReader := bufio.NewReader(secondCmdStdout)

	firstCmd.Start()
	secondCmd.Start()

	firstCmdStdin.Write([]byte("isready\n"))
	secondCmdStdin.Write([]byte("isready\n"))

	// wait for readyok
	for {
		line, _ := firstStdinReader.ReadString('\n')
		line = strings.TrimSpace(line)
		FIRST.Log(line)
		if line == "readyok" {
			break
		}
	}
	for {
		line, _ := secondStdinReader.ReadString('\n')
		line = strings.TrimSpace(line)
		SECOND.Log(line)
		if line == "readyok" {
			break
		}
	}

	moves := []string{}
	sfens := []string{}

	side := FIRST
	currentStdinReader := firstStdinReader
	currentStdinWriter := firstCmdStdin
	stop := false
	for maxCount := 500; maxCount > 0 && !stop; maxCount-- {
		verbose(strings.Repeat("-", 20))

		usi := "position startpos"
		if len(moves) > 0 {
			usi += " moves " + strings.Join(moves, " ")
		}
		currentStdinWriter.Write([]byte(usi + "\n"))
		currentStdinWriter.Write([]byte(fmt.Sprintf("go byoyomi %d\n", m.Byoyomi)))
		for {
			// bestmoveを受け取るまで待つ
			line, _ := currentStdinReader.ReadString('\n')
			line = strings.TrimSpace(line)
			side.Log(line)
			if strings.HasPrefix(line, "bestmove") {
				// format: bestmove [move]
				//   move: 'resign', sfen
				move := strings.TrimPrefix(line, "bestmove ")
				if move == "resign" {
					stop = true
					break
				}
				moves = append(moves, move)
				break
			}
		}

		usi = fmt.Sprintf("position startpos moves %s\n", strings.Join(moves, " "))
		currentStdinWriter.Write([]byte(usi))

		currentStdinWriter.Write([]byte("sfen\n"))

		// sfenを受け取るまで待つ
		line, _ := currentStdinReader.ReadString('\n')
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "sfen") {
			// format: sfen [sfen]
			sfen := strings.TrimPrefix(line, "sfen ")
			sfens = append(sfens, sfen)
			verbose(fmt.Sprintf("%s", parseSfen(sfen)))
		} else {
			verbose(fmt.Sprintf("invalid sfen format: %s", line))
			break
		}

		if side == FIRST {
			side = SECOND
			currentStdinReader = secondStdinReader
			currentStdinWriter = secondCmdStdin
		} else {
			side = FIRST
			currentStdinReader = firstStdinReader
			currentStdinWriter = firstCmdStdin
		}
	}

	verbose(strings.Repeat("-", 20))
	verbose(fmt.Sprintf("winner: %s, moves: %s", side, strings.Join(moves, " ")))

	firstCmdStdin.Write([]byte("quit\n"))
	secondCmdStdin.Write([]byte("quit\n"))
	firstCmd.Wait()
	secondCmd.Wait()

	return &Result{
		moves:  moves,
		sfens:  sfens,
		winner: side,
	}, nil
}
