/*
 * Copyright (c) 2022. All rights reserved by XtraVisions.
 */

package xmgo

var (
	onConnected = make(map[string][]onConnectedCallback, 0)
	onOpened    = make(map[string]map[string][]onOpenedCallback, 0)
)

type onConnectedCallback struct {
	Name string
	Fn   func(*Client) error
}

type onOpenedCallback struct {
	Name string
	Fn   func(database *Database) error
}

func OnConnected(uri string, name string, fn func(*Client) error) {
	var hooks []onConnectedCallback
	var ok bool

	if hooks, ok = onConnected[uri]; !ok {
		hooks = make([]onConnectedCallback, 0)
	}

	hooks = append(hooks, onConnectedCallback{Name: name, Fn: fn})
	onConnected[uri] = hooks
}

func OnOpened(uri string, db string, name string, fn func(database *Database) error) {
	var cli map[string][]onOpenedCallback
	var hook []onOpenedCallback
	var ok bool

	if cli, ok = onOpened[uri]; !ok {
		cli = make(map[string][]onOpenedCallback, 0)
	}

	if hook, ok = cli[db]; !ok {
		hook = make([]onOpenedCallback, 0)
	}

	hook = append(hook, onOpenedCallback{Name: name, Fn: fn})
	cli[db] = hook
	onOpened[uri] = cli
}
