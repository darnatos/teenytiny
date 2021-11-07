package emitter

import (
	"fmt"
	"os"
)

type Emitter interface {
	WriteFile()
	HeaderLine(string)
	Emit(string)
	EmitLine(string)
}

type emitterImpl struct {
	fullPath string
	header   []rune
	code     []rune
}

func Constructor(fullPath string) Emitter {
	return Emitter(&emitterImpl{fullPath, []rune{}, []rune{}})
}

func (emt *emitterImpl) WriteFile() {
	f, err := os.Create(emt.fullPath)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(f)
	for i := 0; i < len(emt.header); i++ {
		j := i
		for i < len(emt.header) && emt.header[i] != '\n' {
			i++
		}
		_, err := f.WriteString(string(emt.header[j : i+1]))
		if err != nil {
			panic(err)
		}
	}
	for i := 0; i < len(emt.code); i++ {
		j := i
		for i < len(emt.code) && emt.code[i] != '\n' {
			i++
		}
		_, err := f.WriteString(string(emt.code[j : i+1]))
		if err != nil {
			panic(err)
		}
	}
}

func (emt *emitterImpl) Emit(code string) {
	emt.code = append(emt.code, []rune(code)...)
}

func (emt *emitterImpl) EmitLine(code string) {
	emt.code = append(emt.code, []rune(code)...)
	emt.code = append(emt.code, '\n')
}

func (emt *emitterImpl) HeaderLine(code string) {
	emt.header = append(emt.header, []rune(code)...)
	emt.header = append(emt.header, '\n')
}
