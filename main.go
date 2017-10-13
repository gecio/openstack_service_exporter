package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/innovocloud/openstack_service_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// Build information. Populated ad build time.
var (
	Version   string
	Revision  string
	Branch    string
	BuildDate string
	GoVersion = runtime.Version()
)

// AuthenticatedClient Like openstack.AuthenticatedClient, but replaces client.HTTPClient.Transport by one with
// TLS verification disabled, if necessary
func AuthenticatedClient(options gophercloud.AuthOptions, insecure bool) (*gophercloud.ProviderClient, error) {
	client, err := openstack.NewClient(options.IdentityEndpoint)
	if err != nil {
		return nil, err
	}

	if insecure {
		config := &tls.Config{InsecureSkipVerify: true}
		transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}

		client.HTTPClient.Transport = transport
	}

	err = openstack.Authenticate(client, options)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func printVersion(w io.Writer) {
	fmt.Fprintf(w, "Version: %s\n", Version)
	fmt.Fprintf(w, "Revision: %s\n", Revision)
	fmt.Fprintf(w, "Branch: %s\n", Branch)
	fmt.Fprintf(w, "BuildDate: %s\n", BuildDate)
	fmt.Fprintf(w, "GoVersion: %s\n", GoVersion)
}

func main() {
	listenAddress := flag.String("web.listen-address", ":9177", "Address to listen on for web interface and telemetry.")
	version := flag.Bool("version", false, "show version.")
	insecure := flag.Bool("insecure", false, "Skip TLS verify.")
	endpointType := flag.String("openstack.endpoint_type", "public", "OpenStack endpoint to use (admin, public, internal).")
	region := flag.String("openstack.region", "", "OpenStack region of the services.")
	metricsPath := flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	enableBlockstorage := flag.Bool("collector.blockstorage", true, "Enable collector for blockstorage services.")
	enableCompute := flag.Bool("collector.compute", true, "Enable collector for compute services.")
	enableNetwork := flag.Bool("collector.network", true, "Enable collector for network agents.")
	enableOrchestration := flag.Bool("collector.orchestration", false, "Enable collector for orchestration services.")

	flag.Parse()

	if *version {
		printVersion(os.Stdout)
		os.Exit(0)
	}
	endpointOpts := gophercloud.EndpointOpts{Region: *region}

	endpointTypeEnv := os.Getenv("OS_INTERFACE")
	if endpointTypeEnv == "" {
		endpointTypeEnv = "public"
	}
	if *endpointType != "public" {
		endpointTypeEnv = *endpointType
	}
	switch endpointTypeEnv {
	case "admin":
		endpointOpts.Availability = gophercloud.AvailabilityAdmin
	case "public":
		endpointOpts.Availability = gophercloud.AvailabilityPublic
	case "internal":
		endpointOpts.Availability = gophercloud.AvailabilityInternal
	default:
		log.Fatalf("No such endpoint_type %s. Use one of admin, public, internal.", *endpointType)
	}

	collectors := []string{}
	if *enableBlockstorage {
		collectors = append(collectors, "blockstorage")
	}
	if *enableCompute {
		collectors = append(collectors, "compute")
	}
	if *enableNetwork {
		collectors = append(collectors, "network")
	}
	if *enableOrchestration {
		collectors = append(collectors, "orchestration")
	}
	if len(collectors) == 0 {
		log.Fatalf("No collector enabled. Bailing out.")
	}

	authOpts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		log.Fatalf("Unable to get openstack auth options. Set openstack environment: %v", err)
	}
	authOpts.AllowReauth = true

	provider, err := AuthenticatedClient(authOpts, *insecure)
	if err != nil {
		log.Fatalf("Unable to get openstack provider: %v", err)
	}

	provider.HTTPClient.Timeout = 2 * time.Second

	exporter, err := collector.NewCollector(provider, endpointOpts, collectors...)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Using collectors: %s", strings.Join(collectors, ","))
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<html>
<head><title>OpenStack Service Exporter</title></head>
<body>
<h1>OpenStack Service Exporter</h1>
<p><a href='%s'>Metrics</a></p>
<h2>Config</h2>
<pre>`, *metricsPath)

		fmt.Fprintf(w, "UserID: %s\n", authOpts.UserID)
		fmt.Fprintf(w, "Username: %s\n", authOpts.Username)
		fmt.Fprintf(w, "DomainID: %s\n", authOpts.DomainID)
		fmt.Fprintf(w, "DomainName: %s\n", authOpts.DomainName)
		fmt.Fprintf(w, "TenantID: %s\n", authOpts.TenantID)
		fmt.Fprintf(w, "TenantName: %s\n", authOpts.TenantName)
		fmt.Fprintf(w, "AllowReauth: %v\n", authOpts.AllowReauth)
		fmt.Fprintf(w, "EndpointType: %v\n", endpointTypeEnv)

		fmt.Fprintf(w, `</pre>
<h2>Build</h2>
<pre>`)

		printVersion(w)

		fmt.Fprintln(w, `</pre>
</body>
</html>`)
	})

	log.Infof("Listen: %s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))

	os.Exit(0)
}
