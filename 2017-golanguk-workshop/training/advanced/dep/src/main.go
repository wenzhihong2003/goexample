package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Error(errors.New("boom!"))
}
