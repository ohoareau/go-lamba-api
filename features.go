package gola

import (
	"github.com/ohoareau/gola/features"
)

func applyFeatures(r Router, f Features) {
	if f["logger"] {
		features.Logger(r)
	}
	if f["recoverer"] {
		features.Recoverer(r)
	}
	if f["cors"] {
		features.Cors(r)
	}
	if f["ping"] {
		features.Ping(r)
	}
	if f["root"] {
		features.Root(r)
	}
	if f["profiler"] {
		features.Profiler(r)
	}
	if f["notfound"] {
		features.Notfound(r)
	}
	if f["jwtauth"] {
		features.JwtAuth(r)
	}
}
