package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/martinnirtl/dynatrace-kube-hunter-ingester/internal/dynatrace"
	"github.com/martinnirtl/dynatrace-kube-hunter-ingester/pkg/kubehunter"
	"github.com/spf13/viper"
)

func report(w http.ResponseWriter, r *http.Request) {
	defer func() {
		StopServer()
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	// unquote is necessary as JSON received from kube-hunter is not clean
	unqotedBodyString, _ := strconv.Unquote(string(body))

	report, err := kubehunter.ParseReport([]byte(unqotedBodyString))
	if err != nil {
		log.Printf("BODY:\n%s", string(body))

		log.Fatal(err.Error())
	}

	apiBaseUrl := viper.GetString("api-url")
	token := viper.GetString("token")

	// TODO run in goroutines
	if viper.GetString("ingest-as") == "events" || viper.GetString("ingest-as") == "both" {
		dynatrace.IngestReportAsEventsV2(apiBaseUrl, token, report)
	}
	if viper.GetString("ingest-as") == "logs" || viper.GetString("ingest-as") == "both" {
		dynatrace.IngestReportAsLogs(apiBaseUrl, token, report)
	}
}