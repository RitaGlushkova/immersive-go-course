package main

import (
	cmd "github.com/RitaGlushkova/raft-otel/command"
	rt "github.com/RitaGlushkova/raft-otel/raft"
)

type ServerClient struct {
	cmd.UnimplementedCommandServer
}

type ServerRaft struct {
	rt.UnimplementedRaftServer
}

type Entry struct {
	Key   string
	Value int64
}

type PersistentState struct {
	currentTerm int64
	votedFor    int64
	log         []Entry
}
type VolatileState struct {
	commitIndex int64
	lastApplied int64
}

type VolatileStateLeader struct {
	nextIndex  []int64
	matchIndex []int64
}
