package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/msoerjanto/thepikaso/data"
	"github.com/msoerjanto/thepikaso/exhibition"
	"github.com/msoerjanto/thepikaso/piece"
	"github.com/msoerjanto/thepikaso/space"
)

const (
	defaultPort              = "8080"
	defaultRoutingServiceURL = "http://localhost:7878"
)

func main() {
	var (
		addr = envString("PORT", defaultPort)

		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")
	)

	flag.Parse()

	data.CreateTables()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	var (
		pieces  = piece.NewPieceRepository()
		artists = piece.NewArtistRepository()
		spaces  = space.NewSpaceRepository()
	)

	var es exhibition.ExhibitionService
	es = exhibition.NewService(pieces, artists, spaces)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/exhibition/v1/", exhibition.MakeHandlers(es, httpLogger))

	http.Handle("/", accessControl(mux))

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

// import (
// 	"net/http"
// 	"os"

// 	"github.com/go-kit/kit/log"

// 	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
// 	stdprometheus "github.com/prometheus/client_golang/prometheus"
// 	"github.com/prometheus/client_golang/prometheus/promhttp"

// 	httptransport "github.com/go-kit/kit/transport/http"
// )

// func main() {
// 	logger := log.NewLogfmtLogger(os.Stderr)

// 	fieldKeys := []string{"method", "error"}
// 	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
// 		Namespace: "my_group",
// 		Subsystem: "piece_service",
// 		Name:      "request_count",
// 		Help:      "Number of requests received.",
// 	}, fieldKeys)
// 	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
// 		Namespace: "my_group",
// 		Subsystem: "piece_service",
// 		Name:      "request_latency_microseconds",
// 		Help:      "Total duration of requests in microseconds.",
// 	}, fieldKeys)
// 	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
// 		Namespace: "my_group",
// 		Subsystem: "piece_service",
// 		Name:      "count_result",
// 		Help:      "The result of each count method.",
// 	}, []string{}) // no fields here

// 	var svc StringService
// 	svc = stringService{}
// 	svc = loggingMiddleware{logger, svc}
// 	svc = instrumentingMiddleware{requestCount, requestLatency, countResult, svc}

// 	createPieceHandler := httptransport.NewServer(
// 		makeCreatePieceEndpoint(svc),
// 		decodeCreatePieceRequest,
// 		encodeResponse,
// 	)

// 	getPiecesHandler := httptransport.NewServer(
// 		makeGetPiecesEndpoint(svc),
// 		decodeGetPiecesRequest,
// 		encodeResponse,
// 	)

// 	http.Handle("/uppercase", uppercaseHandler)
// 	http.Handle("/count", countHandler)
// 	http.Handle("/metrics", promhttp.Handler())
// 	logger.Log("msg", "HTTP", "addr", ":8080")
// 	logger.Log("err", http.ListenAndServe(":8080", nil))
// }

// package main

// import (
// 	"fmt"

// 	"github.com/msoerjanto/thepikaso/data"
// 	"github.com/msoerjanto/thepikaso/data/artist"
// 	"github.com/msoerjanto/thepikaso/data/piece"
// 	"github.com/msoerjanto/thepikaso/data/space"
// )

// func main() {
// 	fmt.Println("Starting Pikaso...")
// 	data.CreateTables()

// 	samplePiece := piece.Piece{
// 		PieceId: "1+1",
// 		Year:    1999,
// 		Length:  123,
// 		Height:  123,
// 	}

// 	piece.CreateNewPiece(samplePiece)

// 	sampleSpace := space.Space{
// 		Location:    "Main room",
// 		SpaceNumber: 1,
// 	}

// 	space.CreateNewSpace(sampleSpace)

// 	sampleArtist := artist.Artist{
// 		ArtistId:  1,
// 		FirstName: "Nasirun",
// 	}

// 	artist.CreateNewArtist(sampleArtist)
// }
