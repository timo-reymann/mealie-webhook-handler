package main

import (
	_ "embed"

	"go.deepl.dev/mealie-webhook-handler/cmd"
)

//go:embed NOTICE
var NoticeContent []byte

func main() {
	cmd.Execute(NoticeContent)
}
