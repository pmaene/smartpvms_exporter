package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pmaene/smartpvms_exporter/internal/collectors"
	"github.com/pmaene/smartpvms_exporter/internal/smartpvms"
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
	listenAddress        string
	metricsPath          string
	readHeaderTimeout    time.Duration
	spvmsBaseURL         string
	spvmsUsername        string
	spvmsPassword        string
	spvmsPasswordFile    string
	spvmsRefreshInterval time.Duration

	rootCmd = &cobra.Command{
		Use:          "smartpvms_exporter",
		Short:        "SmartPVMS Exporter",
		SilenceUsage: true,
		Run:          runRoot,
		PreRun:       runStartPre,
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

	rootCmd.Flags().DurationVar(
		&readHeaderTimeout,
		"web.read-header-timeout",
		5*time.Second,
		"timeout for reading request headers",
	)

	rootCmd.Flags().StringVar(
		&spvmsBaseURL,
		"smartpvms.base-url",
		"https://eu5.fusionsolar.huawei.com",
		"base url of the management system",
	)

	rootCmd.Flags().StringVar(
		&spvmsUsername,
		"smartpvms.username",
		"",
		"username to authenticate against the management system",
	)

	rootCmd.Flags().StringVar(
		&spvmsPassword,
		"smartpvms.password",
		"",
		"password to authenticate against the management system",
	)

	rootCmd.Flags().StringVar(
		&spvmsPasswordFile,
		"smartpvms.password-file",
		"",
		"path to the password to authenticate against the management system",
	)

	rootCmd.Flags().DurationVar(
		&spvmsRefreshInterval,
		"smartpvms.refresh-interval",
		5*time.Second,
		"interval at which to query the management system",
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

func runStartPre(cmd *cobra.Command, args []string) {
	if viper.GetString("smartpvms.password-file") != "" {
		buf, err := os.ReadFile(viper.GetString("smartpvms.password-file"))
		if err != nil {
			log.Fatal(err)
		}

		viper.Set("smartpvms.password", string(buf))
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	// arguments
	if viper.GetString("smartpvms.username") == "" {
		log.Fatal("management system username not set")
	}

	if viper.GetString("smartpvms.password") == "" {
		log.Fatal("management system password not set")
	}

	// main
	log.Infoln("starting", cmd.Name(), cmd.Version)

	cfg := &smartpvms.Config{
		BaseURL:  viper.GetString("smartpvms.base-url"),
		Username: viper.GetString("smartpvms.username"),
		Password: viper.GetString("smartpvms.password"),
	}

	{
		c := collectors.NewPlantsCollector(
			cfg.Client(),
			viper.GetDuration("smartpvms.refresh-interval"),
			log.Base(),
		)

		if err := prometheus.Register(c); err != nil {
			log.Fatal(err)
		}
	}

	{
		c := collectors.NewResidentialInvertersCollector(
			cfg.Client(),
			viper.GetDuration("smartpvms.refresh-interval"),
			log.Base(),
		)

		if err := prometheus.Register(c); err != nil {
			log.Fatal(err)
		}
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

	s := http.Server{
		ReadHeaderTimeout: viper.GetDuration("web.read-header-timeout"),
	}

	l, err := net.Listen(
		"tcp",
		viper.GetString("web.listen-address"),
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Infoln("listening on", viper.GetString("web.listen-address"))
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
