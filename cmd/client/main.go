package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go-keeper/config/client"
	"go-keeper/internal/api"
	"go-keeper/internal/command"
	"go-keeper/internal/session"
)

const (
	RegisterCMD = "reg"
	AuthCMD     = "auth"
	LogPassCMD  = "logpass"
	CardCMD     = "card"
	TextCMD     = "text"
	FileCMD     = "file"
)

func main() {

	cfg, err := client.LoadConfig()
	if err != nil {
		return
	}

	httpClient := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   time.Second * time.Duration(cfg.Timeout),
	}

	clientSession := session.NewClientSession()
	dataAPI := api.NewDataAPI(cfg.ServerAddr, httpClient, clientSession)

	cmdManager := command.NewManager()
	cmdManager.AddCommand(
		RegisterCMD, command.NewRegisterCMD(httpClient, clientSession, cfg.ServerAddr),
	)
	cmdManager.AddCommand(
		AuthCMD, command.NewAuthCMD(httpClient, clientSession, cfg.ServerAddr),
	)
	cmdManager.AddCommand(
		LogPassCMD, command.NewLogPassCMD(dataAPI),
	)
	cmdManager.AddCommand(
		CardCMD, command.NewCardCMD(dataAPI),
	)
	cmdManager.AddCommand(
		TextCMD, command.NewTextCMD(dataAPI),
	)
	cmdManager.AddCommand(
		FileCMD, command.NewFileCMD(dataAPI),
	)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Hello, Keeper...")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		for {
			if !clientSession.IsAuth() {
				fmt.Print("no auth")
			}

			fmt.Print("> ")
			text, _ := reader.ReadString('\n')
			text = strings.Trim(text, "\n\r")
			parts := strings.Split(text, " ")
			if parts[0] == "q" {
				fmt.Println("exit...")
				sig <- syscall.SIGQUIT
				break
			}

			err = cmdManager.RunCommand(parts[0], parts[1:])

			if errors.Is(err, command.ErrCommandNotFound) {
				fmt.Println("command not found")
				continue
			}

			if err != nil {
				fmt.Println("executed with error -", err.Error())
			}
		}
	}()

	<-sig

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer shutdownCtxCancel()

	go func() {
		<-shutdownCtx.Done()
		if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
			fmt.Println("Something goes wrong in exiting from app...forcing exit")
			os.Exit(1)
		}
	}()
}
