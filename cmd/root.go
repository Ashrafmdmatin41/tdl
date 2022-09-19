package cmd

import (
	"context"
	"fmt"
	"github.com/iyear/tdl/cmd/chat"
	"github.com/iyear/tdl/cmd/dl"
	"github.com/iyear/tdl/cmd/login"
	"github.com/iyear/tdl/cmd/up"
	"github.com/iyear/tdl/cmd/version"
	"github.com/iyear/tdl/pkg/consts"
	"github.com/iyear/tdl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"path/filepath"
)

var cmd = &cobra.Command{
	Use:               "tdl",
	Short:             "Telegram Downloader, but more than a downloader",
	Example:           "tdl -h",
	DisableAutoGenTag: true,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

func init() {
	cmd.AddCommand(version.Cmd, login.Cmd, dl.CmdDL, chat.Cmd, up.Cmd)
	cmd.PersistentFlags().String(consts.FlagProxy, "", "proxy address, only socks5 is supported, format: protocol://username:password@host:port")
	cmd.PersistentFlags().StringP(consts.FlagNamespace, "n", "", "namespace for Telegram session")

	cmd.PersistentFlags().IntP(consts.FlagPartSize, "s", 512*1024, "part size for transfer, max is 512*1024")
	cmd.PersistentFlags().IntP(consts.FlagThreads, "t", 8, "threads for transfer one item")
	cmd.PersistentFlags().IntP(consts.FlagLimit, "l", 2, "max number of concurrent tasks")

	_ = viper.BindPFlags(cmd.PersistentFlags())
	viper.AutomaticEnv()
	viper.SetEnvPrefix("tdl")

	docs := filepath.Join(consts.DocsPath, "command")
	if utils.FS.PathExists(docs) {
		if err := doc.GenMarkdownTree(cmd, docs); err != nil {
			panic(fmt.Errorf("generate cmd docs failed: %v", err))
		}
	}
}

func Execute() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	if err := cmd.ExecuteContext(ctx); err != nil {
		return err
	}
	return nil
}
