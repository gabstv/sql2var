package main

import (
	"bufio"
	"bytes"
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
	var lvar bytes.Buffer
	var lcontent bytes.Buffer
	var rprev, rnow rune
	incomment := false
	var lcomment bytes.Buffer
	for _, fname := range inputfiles {
		fp, err := os.Open(fname)
		if err != nil {
			log.Println("could not open", fname, err.Error())
			continue
		}
		r := bufio.NewReader(fp)
		for {
			rprev = rnow
			rnow, _, err = r.ReadRune()
			if err != nil {
				break
			}
			if rnow == '\n' {
				if incomment {
					// end comment
					// parse if begin or end var
				} else if invar {
					// write '\\n'
				}
				continue
			}
			if incomment {
				// write to last comment
				lcomment.WriteRune(rnow)
				continue
			}
			if rnow == '-' {
				if rprev == '-' {
					//FIXME: handle inside quotes
					incomment = true
					continue
				}
			}
			if invar {
				// write to last var
				// TODO
				lvar.WriteRune(rnow)
				continue
			}
		}
	}
}
