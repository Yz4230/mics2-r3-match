package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var firstStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("2")). // green
	Background(lipgloss.Color("0")). // black
	Bold(true)

var secondStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("3")). // yellow
	Background(lipgloss.Color("0")). // black
	Bold(true)

type Args struct {
	// first player random
	Fr int
	// first player depth
	Fd int
	// second player random
	Sr int
	// second player depth
	Sd int
	// number of games
	Ngames int
	// verbose
	Verbose bool
	// byoyomi
	Byoyomi int
}

var args Args

func init() {
	flag.IntVar(&args.Fr, "fr", 0, "first player random")
	flag.IntVar(&args.Fd, "fd", 7, "first player depth")
	flag.IntVar(&args.Sr, "sr", 0, "second player random")
	flag.IntVar(&args.Sd, "sd", 7, "second player depth")
	flag.IntVar(&args.Ngames, "n", 1, "number of games")
	flag.BoolVar(&args.Verbose, "v", false, "verbose")
	flag.IntVar(&args.Byoyomi, "b", 10000, "byoyomi")

	flag.Parse()
}

func getOutdirName() string {
	// example: fr5-fd6-sr7-sd8
	return fmt.Sprintf("fr%d-fd%d-sr%d-sd%d", args.Fr, args.Fd, args.Sr, args.Sd)
}

func getFilename() string {
	// unix timestamp + random string
	buf := make([]byte, 4)
	rand.Read(buf)
	return fmt.Sprintf("%d-%x.txt", time.Now().Unix(), buf)
}

func createDirIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.Mkdir(dir, 0755)
	}
	return nil
}

func exportToFile(result *Result) error {
	outdir := getOutdirName()
	filename := getFilename()
	f, err := os.Create(fmt.Sprintf("%s/%s", outdir, filename))
	if err != nil {
		return err
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	winner := "first"
	if result.winner == SECOND {
		winner = "second"
	}
	bw.WriteString(fmt.Sprintf("winner:%s\n", winner))
	bw.WriteString(fmt.Sprintf("moves:%d\n", len(result.moves)))
	for i := range result.moves {
		bw.WriteString(fmt.Sprintf("%s,%s\n", result.moves[i], result.sfens[i]))
	}
	bw.Flush()

	return nil
}

func main() {
	outdir := getOutdirName()
	createDirIfNotExists(outdir)
	fmt.Printf("outdir: %s\n", outdir)

	for i := 0; i < args.Ngames || args.Ngames == 0; i++ {
		startTime := time.Now()
		matcher := &Matcher{
			Fr:      args.Fr,
			Fd:      args.Fd,
			Sr:      args.Sr,
			Sd:      args.Sd,
			Byoyomi: args.Byoyomi,
		}
		result, err := matcher.Match()
		if err != nil {
			panic(err)
		}
		elapsed := time.Since(startTime)
		elapsed = elapsed.Round(time.Second)
		fmt.Printf("game %d/%d: %s (%s)\n", i+1, args.Ngames, result.winner, elapsed)
		exportToFile(result)
	}
}
