package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/buengese/qrterminal"
	"rsc.io/qr"
)

var verboseFlag bool
var halfsizeFlag bool
var levelFlag string
var quietZoneFlag int

func getLevel(s string) qr.Level {
	switch l := strings.ToLower(s); l {
	case "l":
		return qr.L
	case "m":
		return qr.M
	case "h":
		return qr.H
	default:
		return -1
	}
}

func main() {
	flag.BoolVar(&verboseFlag, "v", false, "Output debugging information")
	flag.BoolVar(&halfsizeFlag, "s", false, "Use smaller characters")
	flag.StringVar(&levelFlag, "l", "L", "Error correction level")
	flag.IntVar(&quietZoneFlag, "q", 2, "Size of quietzone border")

	flag.Parse()
	level := getLevel(levelFlag)

	if len(flag.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "usage of %s: \"[arguments]\"\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	} else if level < 0 {
		fmt.Fprintf(os.Stderr, "Invalid error correction level: %s\n", levelFlag)
		fmt.Fprintf(os.Stderr, "Valid options are [L, M, H]\n")
		os.Exit(1)
	}

	var cfg qrterminal.Config
	if !halfsizeFlag {
		cfg = qrterminal.Config{
			Level:     level,
			Writer:    os.Stdout,
			QuietZone: quietZoneFlag,
			BlackChar: qrterminal.BLACK,
			WhiteChar: qrterminal.WHITE,
		}
	} else {
		cfg = qrterminal.Config{
			Level:     level,
			Writer:    os.Stdout,
			QuietZone: quietZoneFlag,
			HalfBlocks:     true,
			BlackChar:      qrterminal.BLACK_BLACK,
			WhiteBlackChar: qrterminal.WHITE_BLACK,
			WhiteChar:      qrterminal.WHITE_WHITE,
			BlackWhiteChar: qrterminal.BLACK_WHITE,
		}
	}

	if verboseFlag {
		fmt.Fprintf(os.Stdout, "Level: %s \n", levelFlag)
		fmt.Fprintf(os.Stdout, "Quietzone Border Size: %d \n", quietZoneFlag)
		fmt.Fprintf(os.Stdout, "Encoded data: %s \n", strings.Join(flag.Args(), "\n"))
		fmt.Println("")
	}

	if runtime.GOOS == "windows" {
		cfg.Writer = colorable.NewColorableStdout()
		cfg.BlackChar = qrterminal.BLACK
		cfg.WhiteChar = qrterminal.WHITE
	}

	qrterminal.GenerateWithConfig(strings.Join(flag.Args(), "\n"), cfg)
}
