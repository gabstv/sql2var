package main

/*
$GOARCH
	The execution architecture (arm, amd64, etc.)
$GOOS
	The execution operating system (linux, windows, etc.)
$GOFILE
	The base name of the file.
$GOLINE
	The line number of the directive in the source file.
$GOPACKAGE
	The name of the package of the file containing the directive.
$DOLLAR
	A dollar sign.

https://golang.org/cmd/go/#hdr-Generate_Go_files_by_processing_source
*/

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mkideal/cli"
)

type argT struct {
	Inputs  []string `cli:"I" usage:"Input files"`
	Output  string   `cli:"O" usage:"Output file"`
	Package string   `cli:"P" usage:"Go Package"`
}

func main() {
	args := new(argT)
	cli.Run(args, func(ctx *cli.Context) error {
		inputs := args.Inputs
		if len(inputs) == 0 {
			ctx.Color().Red(ctx.Usage())
			return nil
		}
		source := ""
		output := args.Output
		if output == "" {
			if len(ctx.Args()) < 1 && os.Getenv("GOFILE") == "" {
				ctx.Color().Red(ctx.Usage())
				return nil
			} else if len(ctx.Args()) > 0 {
				output = strings.TrimSuffix(ctx.Args()[0], ".go") + "_sql2var.go"
				source = ctx.Args()[0]
			} else if os.Getenv("GOFILE") != "" {
				output = strings.TrimSuffix(os.Getenv("GOFILE"), ".go") + "_sql2var.go"
				source = os.Getenv("GOFILE")
			}
		} else {
			source = os.Getenv("GOFILE")
		}
		sqlk := make([]string, 0)
		sqlv := make([]string, 0)
		sqltags := make([][]string, 0)
		extractall(inputs, &sqlk, &sqlv, &sqltags)
		//
		if args.Package == "" {
			args.Package = os.Getenv("GOPACKAGE")
		}
		//fmt.Println("PACKAGE", args.Package)
		fmt.Println(ctx.Color().Green("sql2var"), output)

		ef, err := os.OpenFile(output, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer ef.Close()
		eb := bufio.NewWriter(ef)
		eb.WriteString("// Code generated by sql2var <https://github.com/gabstv/sql2var>. DO NOT EDIT.\n")
		eb.WriteString("// source: ")
		eb.WriteString(source)
		eb.WriteString("\n\npackage ")
		eb.WriteString(args.Package)
		eb.WriteString("\n\n")
		for k := range sqlk {
			skipit := false
			if len(sqltags[k]) > 0 {
				for _, vv := range sqltags[k] {
					if vv == "skip" {
						skipit = true
					}
				}
			}
			if !skipit {
				eb.WriteString(fmt.Sprintf("const %s = %q\n", sqlk[k], sqlv[k]))
			}
		}
		eb.Flush()
		return nil
	})
}

func extractall(inputfiles []string, mk, mv *[]string, mtags *[][]string) {
	for _, fname := range inputfiles {
		fp, err := os.Open(fname)
		if err != nil {
			log.Println("could not open", fname, err.Error())
			continue
		}
		invar := false
		incomment := false
		curline := 1 // will be used to throw errors
		curpos := -1 // will be used to throw errors
		var lvar bytes.Buffer
		var lcontent bytes.Buffer
		var rprev, rnow rune
		var lcomment bytes.Buffer
		lasttags := make([]string, 0)
		r := bufio.NewReader(fp)
		hasEOF := false
		for !hasEOF {
			curpos++
			rprev = rnow
			rnow, _, err = r.ReadRune()
			if err != nil {
				if err.Error() == "EOF" {
					hasEOF = true
					rnow = '\n'
					fp.Close()
					//fmt.Print("(EOF)")
					//spew.Dump(curline)
				} else {
					break
				}
			}
			//fmt.Print(string(rnow))
			if rnow == '\n' {
				if incomment {
					// end comment
					// parse if begin or end var
					incomment = false
					lc := strings.TrimSpace(lcomment.String())
					if strings.HasPrefix(lc, "sql:") && invar {
						skey := strings.TrimSpace(lc[4:])
						cpk0 := *mk
						cpv0 := *mv
						found := false
						for k := range cpk0 {
							if cpk0[k] == skey {
								lcontent.WriteString(cpv0[k])
								found = true
								break
							}
						}
						if !found {
							fmt.Println("key not found (sql:):", skey)
						}
						lcomment.Reset()
					} else if strings.HasPrefix(lc, "define:") && !invar {
						rawt := strings.Split(strings.TrimSpace(lc[7:]), ";")
						lvar.WriteString(rawt[0])
						if len(rawt) > 1 {
							for k, v := range rawt {
								if k == 0 {
									continue
								}
								lasttags = append(lasttags, strings.TrimSpace(v))
							}
						}
						//fmt.Println("[VAR:", lc[7:], "]")
						invar = true
						lcomment.Reset()
					} else if lc == "end" {
						//fmt.Println("[END]")
						if invar {
							//fmt.Println("[PUSH", lvar.String(), "]")
							// push everything
							*mk = append(*mk, lvar.String())
							*mv = append(*mv, lcontent.String())
							*mtags = append(*mtags, lasttags)
							invar = false
							incomment = false
							lvar.Reset()
							lcontent.Reset()
							lcomment.Reset()
							lasttags = make([]string, 0)
						} else {
							incomment = false
							lcomment.Reset()
						}
					} else {
						//fmt.Println("[?COMMENT:", lcomment.String(), "]")
						incomment = false
						lcomment.Reset()
					}
				} else if invar {
					// write '\n'
					lcontent.WriteString("\n")
				}
				curline++
				curpos = -1
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
					//fmt.Print("(C)")
					incomment = true
					continue
				}
			}
			if invar {
				if rprev == '-' {
					lcontent.WriteRune(rprev)
					//fmt.Println("IN VAR, WROTE -", string(rnow), "???")
				}
				// write to last var
				if rnow != '-' {
					lcontent.WriteRune(rnow)
				}
				continue
			}
		}
	}
}
