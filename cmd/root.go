package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/pmaene/smartpvms_exporter/internal"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version string
	commit  string
)

var (
	listenAddress string
	metricsPath   string

	rootCmd = &cobra.Command{
		Use:          "smartpvms_exporter",
		Short:        "SmartPVMS Exporter",
		SilenceUsage: true,
		Run:          runRoot,
	}
)

func GetVersion() string {
	if version == "" {
		return "(devel)"
	}

	return version
}

func SetVersion(v string) {
	version = v
}

func GetCommit() string {
	return commit
}

func SetCommit(c string) {
	commit = c
}

func Execute() {
	rootCmd.Version = fmt.Sprintf(
		"%s (%s)",
		GetVersion(),
		GetCommit(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(
		&listenAddress,
		"web.listen-address",
		":9867",
		"address on which to expose metrics and web interface",
	)

	rootCmd.Flags().StringVar(
		&metricsPath,
		"web.telemetry-path",
		"/metrics",
		"path under which to expose metrics",
	)

	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("SMARTPVMS_EXPORTER")
	viper.SetEnvKeyReplacer(
		strings.NewReplacer(".", "_", "-", "_"),
	)
}

func runRoot(cmd *cobra.Command, args []string) {
	log.Infoln("starting", cmd.Name(), cmd.Version)

	c := internal.NewCollector()
	if err := prometheus.Register(c); err != nil {
		log.Fatal(err)
	}

	http.Handle(
		viper.GetString("web.telemetry-path"),
		promhttp.Handler(),
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(
			[]byte(
				`<html>
				<head><title>SmartPVMS Exporter</title></head>
				<body>
				<h1>SmartPVMS Exporter</h1>
				<p><a href='/metrics'>Metrics</a></p>
				</body>
				</html>`,
			),
		)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	log.Infoln("listening on", viper.GetString("web.listen-address"))
	if err := http.ListenAndServe(viper.GetString("web.listen-address"), nil); err != nil {
		log.Fatal(err)
	}
}
