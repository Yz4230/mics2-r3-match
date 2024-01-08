package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"os"
	"path"
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
	// first player program
	Fp string
	// first player random
	Fr float64
	// second player program
	Sp string
	// second player random
	Sr float64
	// number of games
	Ngames int
	// verbose
	Verbose bool
	// byoyomi
	Byoyomi int
	// no output
	NoOutput bool
	// output directory
	Outdir string
}

var args Args

func init() {
	flag.StringVar(&args.Fp, "fp", "./minishogi", "first player program")
	flag.Float64Var(&args.Fr, "fr", 0, "first player random rate")
	flag.StringVar(&args.Sp, "sp", "./minishogi", "second player program")
	flag.Float64Var(&args.Sr, "sr", 0, "second player random rate")
	flag.IntVar(&args.Ngames, "n", 1, "number of games")
	flag.BoolVar(&args.Verbose, "v", false, "verbose")
	flag.IntVar(&args.Byoyomi, "b", 10000, "byoyomi")
	flag.BoolVar(&args.NoOutput, "no-output", false, "no output")
	flag.StringVar(&args.Outdir, "outdir", "", "output directory")

	flag.Parse()
}

func checkRandomPlayer() {
	// check './minishogi-ramdom' exists
	if _, err := os.Stat("./minishogi-random"); os.IsNotExist(err) {
		fmt.Println("error: './minishogi-random' not found")
		os.Exit(1)
	}
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
	filename := getFilename()
	f, err := os.Create(path.Join(args.Outdir, filename))
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
	checkRandomPlayer()

	if !args.NoOutput {
		createDirIfNotExists(args.Outdir)
		fmt.Printf("outdir: %s\n", args.Outdir)
	}

	for i := 0; i < args.Ngames || args.Ngames == 0; i++ {
		startTime := time.Now()
		matcher := &Matcher{
			Fp:      args.Fp,
			Sp:      args.Sp,
			Fr:      args.Fr,
			Sr:      args.Sr,
			Byoyomi: args.Byoyomi,
		}
		result, err := matcher.Match()
		if err != nil {
			panic(err)
		}
		elapsed := time.Since(startTime)
		elapsed = elapsed.Round(time.Second)
		fmt.Printf("game %d/%d: %s (%s)\n", i+1, args.Ngames, result.winner, elapsed)
		if !args.NoOutput {
			exportToFile(result)
		}
	}
}
