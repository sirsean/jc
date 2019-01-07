package commands

import (
    "log"
	"github.com/sirsean/jc/path"
	"github.com/codeskyblue/go-sh"
	"os"
)

func Body() {
	id, err := getArgId()
	if err != nil {
		log.Fatal(err)
	}
	path := path.RequestBodyPath(id)
	sh.Command("vim", path).SetStdin(os.Stdin).Run()
}
