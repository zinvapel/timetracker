package main

import (
	"context"
	"github.com/zinvapel/timetracker/contract"
	"github.com/zinvapel/timetracker/cron"
	sreact "github.com/zinvapel/timetracker/sheets/reaction"
	"github.com/zinvapel/timetracker/telegram"
	tgreact "github.com/zinvapel/timetracker/telegram/reaction"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	pc := contract.ProviderChain{
		Chain: []contract.Provider{
			&telegram.Provider{},
			&cron.Provider{},
			contract.GetEventBus(),
		},
	}
	rc := contract.ReactionChain{
		Reactions: []contract.Reaction{
			&contract.ReactionHealth{},
			&contract.ReactionLog{},
			&tgreact.CheckUpdate{},
			&tgreact.Log{},
			&tgreact.SendMessage{},
			&tgreact.ResolveTag{},
			tgreact.NewCommand(),
			&tgreact.AskForAuthCode{},
			&tgreact.ResponseForAuthCode{},
			&sreact.CurrentTask{},
			&sreact.UpdateMainPage{},
			&sreact.SendCurrentImmediately{},
		},
	}

	go health()

	ctx, cancelFunc := context.WithCancel(context.Background())

	sChan := make(chan os.Signal, 1)
	signal.Notify(sChan, syscall.SIGTERM, syscall.SIGKILL)

	for {
		select {
		case event := <-pc.GetEventChannel(ctx):
			rc.React(event)
		case <-sChan:
			cancelFunc()
			contract.Finish(contract.Success)
		}
	}
}