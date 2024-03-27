package plugins

import (
    "strings"
)

func init() {
    Register("template", &Template{})
}

type Template struct{}

func (t *Template) Description() string {
    return "A basic boilerplate template for a netsquirrel plugin."
}

func (t *Template) Execute(comm Communicator, pluginDataChan chan<- string) {
    comm.Send("[Plugin] This is a basic boilerplate template for a netsquirrel plugin. Type 'exit' to quit.")
    
    for {
        comm.Send("> ")
        input, err := comm.Receive()
        if err != nil {
            break
        }

        if strings.TrimSpace(input) == "exit" {
            comm.Send("[Plugin] Goodbye!")
            break
        }

        pluginInput := "[Plugin] " + input
        pluginDataChan <- pluginInput
        comm.Send(pluginInput)
    }
}
