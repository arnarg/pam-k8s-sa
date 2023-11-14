package main

type logger interface {
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errf(string, ...interface{})
}
