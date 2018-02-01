package main

import (
	"log"
	"os"
	"strings"

	"github.com/mkideal/cli"
)

type argT struct {
	Inputs []string `cli:"I" usage:"Input files"`
	Output string   `cli:"O" usage:"Output file"`
}

func main() {
	//var inputs []string
	//var output string
	args := new(argT)
	cli.Run(args, func(ctx *cli.Context) error {
		//ctx.JSONln(args.Inputs)
		//ctx.String(ctx.Args()[0])
		inputs := args.Inputs
		if len(inputs) == 0 {
			ctx.Color().Red(ctx.Usage())
			return nil
		}
		output := args.Output
		if output == "" {
			if len(ctx.Args()) < 1 {
				ctx.Color().Red(ctx.Usage())
				return nil
			}
			output = strings.TrimSuffix(ctx.Args()[0], ".go") + "_sql2var.go"
		}
		sqlfns := make(map[string]string)

		return nil
	})
}

func extractall(inputfiles []string, m map[string]string) {
	invar := false
	var lvar string
	for _, fname := range inputfiles {
		fp, err := os.Open(fname)
		if err != nil {
			log.Println("could not open", fname, err.Error())
			continue
		}
		//bufio.
	}
}
