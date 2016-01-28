package main

import (
	"errors"
)

const (
	errSuccess = 1
	errSrcIP   = 2
	errJson    = 3
	errUnkown  = 99
)

func checkSrcIP(ip string) error {
	for _, v := range srcIP {
		if v == ip {
			return nil
		}
	}
	return errors.New("invalid src ip")
}
