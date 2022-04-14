package main

import (
	"bytes"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
	"github.com/valyala/fasttemplate"
)

func readFile(filename string) (io.Reader, bool) {
	v, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	if len(v) > 0 {
		return bytes.NewBuffer(v), true
	}
	return nil, false
}

func command() {
	command := cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, args)
		},
	}
	flag := command.Flags()
	flag.StringArrayP("inputs", "i", nil, "input files")
	flag.StringP("outout", "o", "outout.yaml", "output file name")
	flag.StringArrayP("references", "d", nil, "anchor file directory")
	flag.IntP("mode", "m", 0644, "write file mode")
	exitOrEmpty(command.Execute())
}

func exitOrEmpty(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func isDirectory(name string) bool {
	stat, err := os.Stat(name)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func run(cmd *cobra.Command, args []string) {
	flag := cmd.Flags()
	mode, _ := flag.GetInt("mode")
	inputs, err := flag.GetStringArray("inputs")
	exitOrEmpty(err)
	output, err := flag.GetString("outout")
	exitOrEmpty(err)
	references, err := flag.GetStringArray("references")
	exitOrEmpty(err)

	files := make([]string, 0, len(inputs))
	for _, file := range inputs {
		if isDirectory(file) {
			for _, ext := range []string{"*.yaml", "*.yml"} {
				fs, _ := filepath.Glob(path.Join(file, ext))
				if len(fs) > 0 {
					files = append(files, fs...)
				}
			}
		} else {
			files = append(files, file)
		}
	}

	dest := map[string]interface{}{}
	for i, input := range files {
		src := map[string]interface{}{}
		reader, ok := readFile(input)
		if !ok {
			continue
		}
		dec := yaml.NewDecoder(reader, yaml.RecursiveDir(true), yaml.ReferenceDirs(references...))
		if i == 0 {
			if err := dec.Decode(&dest); err != nil {
				log.Fatalf("file name: %s, error: %s", input, err)
			}
		} else {
			if err := dec.Decode(&src); err != nil {
				log.Fatalf("file name: %s, error: %s", input, err)
			}
		}
		exitOrEmpty(mergo.Map(&dest, src, mergo.WithAppendSlice))
	}
	buf, err := yaml.Marshal(&dest)
	exitOrEmpty(err)
	inputBuffer := bytes.NewBuffer(buf)
	writeBuffer := &bytes.Buffer{}
	tpl := fasttemplate.New(inputBuffer.String(), "${{", "}}")
	_, err = tpl.ExecuteFunc(writeBuffer, func(w io.Writer, tag string) (int, error) {
		v := os.Getenv(tag)
		return w.Write([]byte(v))
	})
	exitOrEmpty(err)

	err = os.WriteFile(output, writeBuffer.Bytes(), fs.FileMode(mode))
	exitOrEmpty(err)
}

func main() {
	command()
}
