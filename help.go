package plugins

import (
    "fmt"
    "strings"
)

func init() {
    Register("help", &Help{})
}

type Help struct{}

func (h *Help) Execute(comm Communicator, pluginDataChan chan<- string) {
    var helpMessage strings.Builder
    helpMessage.WriteString("[Help] Available Commands:\n")
    for name, plugin := range Commands {
        helpMessage.WriteString(fmt.Sprintf("- %s: %s\n", name, plugin.Description()))
    }
    comm.Send(helpMessage.String())
}

func (h *Help) Description() string {
    return "Lists all available commands and their descriptions."
}
