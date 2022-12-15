package main

import ("fmt")

type Buffer struct {
	buf []byte
}

func (b *Buffer) Read(p []byte) (int, error) {
	n := copy(p, b.buf)
	return n, nil
}

func (b *Buffer) Write(p []byte) (int, error) {
	b.buf = append(b.buf, p...)
	return len(p), nil
}

func (b *Buffer) AsString () string {
	return string(b.buf)
}

func main() {
	var buf Buffer
	slice := make([]byte, 100)
	n, err := buf.Read(slice)
	fmt.Println(n, err, slice)

	writeN, writeErr := buf.Write([]byte("Rita"))
	fmt.Println(writeN, writeErr)
	n, err = buf.Read(slice)
	fmt.Println(n, err, slice)

	writeN, writeErr = buf.Write([]byte("Glushkova"))
	fmt.Println(writeN, writeErr)
	n, err = buf.Read(slice)
	fmt.Println(n, err, slice)
}