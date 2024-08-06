package main

import ws "github.com/gorilla/websocket"

// type s = string
type ss = map[string]string
type sss = map[string]map[string]string
type a = any
type sa = map[string]any
type ssa = map[string]map[string]any
type sssa = map[string]map[string]map[string]any
type slss = map[string][]ss             // private chat repertories
type sslss = map[string]map[string][]ss // private chat repertories
type b = bool
type sb = map[string]bool
type ssb = map[string]map[string]bool
type cb map[*ws.Conn]bool
type cs map[*ws.Conn]string
type csa map[*ws.Conn]map[string]any

type nvm = struct{}
