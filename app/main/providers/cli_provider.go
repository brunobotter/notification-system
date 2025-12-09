package providers

import (
	"fmt"

	"github.com/brunobotter/notification-system/app/main/app"
	"github.com/brunobotter/notification-system/app/main/container"
	"github.com/brunobotter/notification-system/app/main/server"
	"github.com/spf13/cobra"
)

type CliServiceProvider struct {
	command []any
}

func NewCliServiceProvider() *CliServiceProvider {
	return &CliServiceProvider{}
}

func (p *CliServiceProvider) Register(c container.Container) {
	c.Singleton(func(app *app.Application, container container.Container) *cobra.Command {
		return &cobra.Command{
			Use:   "int",
			Short: "command line",
			Run: func(cmd *cobra.Command, args []string) {
				srv, err := server.NewServer(container)
				if err != nil {
					panic(fmt.Errorf("nao pode inicializar o server: %v", err))
				}
				srv.Run(cmd.Context())
				app.WaitForShutdownSignal()
			},
		}
	})
}

func (p *CliServiceProvider) Boot(c container.Container, root *cobra.Command) {
	p.registerCommand(c, root)
}

func (p *CliServiceProvider) registerCommand(c container.Container, root *cobra.Command) {
	for _, construct := range p.command {
		cmd, ok := c.Call(construct).(*cobra.Command)
		if !ok {
			panic("command must be a instance")
		}
		root.AddCommand(cmd)
	}
}
