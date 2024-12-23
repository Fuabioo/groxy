package cmd

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Fuabioo/groxy/internal/service"
	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed signature.txt
var signature string

var rootCmd = &cobra.Command{
	Use:   "groxy",
	Short: "A simple tester proxy",
	Long:  `A proxy whose purpose is to test the behavior of a web application`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if !viper.GetBool("colors") {
			log.SetColorProfile(termenv.Ascii)
		}

		fmt.Println(colorizeSignature(signature))

		// check if domain is provided
		domain := viper.GetString("domain")
		if domain == "" {
			if len(args) == 0 {
				log.Fatal("a domain to proxy to is required")
			}
			viper.Set("domain", args[0])
		}

		// set up logging
		if viper.GetBool("debug") {
			log.SetLevel(log.DebugLevel)
		} else if viper.GetBool("verbose") {
			log.SetLevel(log.InfoLevel)
		} else {
			log.SetLevel(log.WarnLevel)
		}
	},
	Run: func(_ *cobra.Command, _ []string) {

		srv, err := service.New()
		if err != nil {
			log.Fatal("could not create service", "err", err)
		}

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", viper.GetInt("port")),
			Handler: srv,
		}

		ctx, cancel := signal.NotifyContext(context.Background(),
			os.Interrupt,
			os.Kill,
		)
		defer cancel()

		log.Infof("🖴 serving at http://localhost%s", server.Addr)
		log.Infof("\t press %s to stop server", highlight("CTRL+C"))
		go func() {
			err = server.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}()

		<-ctx.Done()
		fmt.Println()

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		log.Info("🔻 shutting down")

		err = server.Shutdown(ctx)
		if err != nil {
			log.Error("could not shutdown server", "err", err)
		}
		cancel()

		<-ctx.Done()
		log.Info("👋 bye!")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	var cfgFile string

	flagSet := rootCmd.PersistentFlags()

	flagSet.StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.groxy.yaml)")
	flagSet.Bool("debug", false, "Enable debug mode")
	flagSet.Bool("verbose", false, "Enable verbose mode")
	flagSet.Bool("colors", true, "Use terminal colors on or off")
	flagSet.Bool("insecure", false, "Skip TLS verification")
	flagSet.Int("port", 8080, "Port to listen on")
	flagSet.String("endpoints", "", "yml raw endpoints configuration")

	// bind the falgs to viper
	if err := viper.BindPFlags(flagSet); err != nil {
		log.Fatal("could not bind flags",
			"err", err,
		)
	}

	cobra.OnInitialize(
		viper.AutomaticEnv,
		func() {
			if cfgFile == "" {
				return
			}

			log.Debug("using config file", "file", cfgFile)

			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)

			// If a config file is found, read it in.
			if err := viper.ReadInConfig(); err != nil {
				log.Fatal("could not read config file", "err", err)
			}

			log.Debug("read config file", "file", viper.ConfigFileUsed())
		},
	)
}
