package main

type Message struct {
	Magic   uint32
	Command uint32
	Value   uint32
}
