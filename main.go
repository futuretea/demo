package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/mitchellh/colorstring"
	"github.com/urfave/cli"
)

var (
	colorNum = 5
	color    = 0
	colors   = []string{
		"light_red",
		"light_yellow",
		"light_green",
		"light_blue",
		"light_magenta",
	}
	configFile string
)

type Step struct {
	Description string            `csv:"description"`
	Command     string            `csv:"command"`
	OtherData   map[string]string `csv:"-"`
}

func ReadSteps(fileName string) []Step {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	dec, err := csvutil.NewDecoder(csv.NewReader(f))
	if err != nil {
		log.Fatal(err)
	}

	var steps []Step
	for {
		step := Step{}
		err = dec.Decode(&step)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		steps = append(steps, step)
	}
	return steps
}

func shell(command string) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(output))
	wait()
}

func wait() {
	if _, err := fmt.Scanln(); err != nil {
		log.Fatal(err)
	}
}

func showMessage(color, message string) {
	for _, i := range message {
		time.Sleep(100 * time.Millisecond)
		colorstring.Fprint(os.Stdout, "["+color+"]"+string(i))
	}
	wait()
}

func showDesc(color int, desc string) {
	showMessage(colors[color], desc)
}

func runCommand(color int, command string) {
	showMessage(colors[color], "$ "+command)
	shell(command)
}

func pickColor() int {
	if color == colorNum-1 {
		color = 0
	} else {
		color = color + 1
	}
	return color
}

func run(c *cli.Context) error {
	steps := ReadSteps(configFile)
	for _, step := range steps {
		color = pickColor()
		showDesc(color, step.Description)
		runCommand(color, step.Command)
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "demo"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "config file path (csv format)",
			Destination: &configFile,
			Required:    true,
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
