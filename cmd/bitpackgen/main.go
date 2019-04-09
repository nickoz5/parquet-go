package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"text/template"
)

var (
	pkg    = flag.String("package", "main", "package of the generated code")
	max    = flag.Int("maxwidth", 3, "the bit width at which to stop")
	outPth = flag.String("output", "bitpack.go", "name of the file that is produced, defaults to parquet.go")
)

func main() {
	flag.Parse()
	pb := bitback{Package: *pkg, Max: *max}
	tmpl := template.New("output").Funcs(funcs)
	var err error
	tmpl, err = tmpl.Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range []string{
		bytesTpl,
		intsTpl,
	} {
		var err error
		tmpl, err = tmpl.Parse(t)
		if err != nil {
			log.Fatal(err)
		}
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, pb)
	if err != nil {
		log.Fatal(err)
	}

	gocode, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(*outPth)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(gocode)
	if err != nil {
		log.Fatal(err)
	}

	f.Close()
}

type bitback struct {
	Package string
	Max     int
}

/*
end := 8 / width
			if width > 2 && width%2 > 0 {
				end++
			}
			var out string
			for i := 0; i < end; i++ {
				index := (width * byt) + i
				if index > 7 {
					break
				}

				shift := (index * width) % 8
				and := (1<<uint(width) - 1)
				or := " |\n"
				if index > 0 && width%2 != 0 && i == end-1 {
					and = 7 >> uint(width-(8-shift))
					or = ""
				} else if index > 0 && width%2 != 0 && i == 0 {
					s := 8 - (((index - 1) * width) % 8)
					a := 7 - (7 >> uint(width-s))
					out += fmt.Sprintf("byte((vals[%d]&%d)%s%d) |\n", index-1, a, ">>", s)
				} else if index == 7 || i == end-1 {
					or = ""
				}
				out += fmt.Sprintf("byte((vals[%d]&%d)%s%d)%s", index, and, "<<", shift, or)
			}
			return out
*/

type byt struct {
	I     int
	Or    string
	And   int
	Shift int
	Dir   string
}

var (
	funcs = template.FuncMap{
		"pack": func(width int) [][]byt {
			bs := [][]byt{[]byt{}}
			var x int
			and := 1<<uint(width) - 1
			for i := 0; i < 8; i++ {
				shift := (i * width) % 8
				if shift+width > 8 {
					a1 := 7 >> uint(width-(8-shift))
					a2 := 7 - a1
					s2 := 8 - shift
					bs[x] = append(bs[x],
						byt{
							I:     i,
							And:   a1,
							Shift: shift,
							Dir:   "<<",
						})
					x++
					bs = append(bs, []byt{})
					bs[x] = append(bs[x],
						byt{
							I:     i,
							And:   a2,
							Or:    " |\n",
							Shift: s2,
							Dir:   ">>",
						},
					)
				} else {
					o := " |\n"
					if shift+width == 8 {
						o = ""
					}
					bs[x] = append(bs[x], byt{
						I:     i,
						And:   and,
						Or:    o,
						Shift: shift,
						Dir:   "<<",
					})
					if shift+width == 8 && i < 7 {
						bs = append(bs, []byt{})
						x++
					}
				}
			}
			return bs
		},
		"int64": func(width, i int) string {
			shift := (i * width) % 8
			index := (i * width) / 8
			mask := ((1 << uint(width)) - 1) << uint(shift)
			if mask < (1 << 8) {
				return fmt.Sprintf("(int64(vals[%d] & %d) >> %d),", index, mask, shift)
			}

			return fmt.Sprintf(
				"%s | %s,",
				fmt.Sprintf("(int64(vals[%d] & %d) >> %d)", index, mask&((1<<8)-1), shift),
				fmt.Sprintf("(int64(vals[%d] & %d) << %d)", index+1, mask>>8, 8-shift),
			)
		},
		"N": func(start, end int) (stream chan int) {
			stream = make(chan int)
			go func() {
				for i := start; i <= end; i++ {
					stream <- i
				}
				close(stream)
			}()
			return
		},
	}

	/*



	 */

	tpl = `package {{.Package}}

// This code is generated by github.com/parsyl/parquet.

func Pack(width int, vals []int64) []byte {
	switch width {
		{{range $i := N 1 .Max }}case {{$i}}:
			return pack{{$i}}(vals)
		{{end}}default:
			return []byte{}
	}
}

{{range $i := N 1 .Max}}
func pack{{$i}}(vals []int64) []byte {
return []byte{ {{template "bytes" $i}} }
}
{{end}}

func Unpack(width int, vals []byte) []int64 {
	switch width {
		{{range $i := N 1 .Max }}case {{$i}}:
			return unpack{{$i}}(vals)
		{{end}}default:
			return []int64{}
	}
}

{{range $i := N 1 .Max }}
	   func unpack{{$i}}(vals []byte) []int64 { {{template "ints" .}}
	   }
{{end}}
`

	bytesTpl = `{{define "bytes"}}
{{ $bytes := pack .}} {{range $byte := $bytes}} ( {{ range $b := $byte}} byte((vals[{{$b.I}}]&{{$b.And}}){{$b.Dir}}{{$b.Shift}}){{$b.Or}}{{end}} ),
{{end}}
{{end}}`
	intsTpl = `{{define "ints"}}{{$width := .}}
return []int64{
{{range $i := N 0 7}} {{int64 $width $i}}
{{end}} }{{end}}`
)
