package config

import (
	"io"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func InitLogger() {
	zerolog.TimeFieldFormat = time.RFC3339

	var logWriters []io.Writer
	logWriters = append(logWriters, os.Stdout)

	conn, err := net.Dial("tcp", "localhost:5000")
	if err == nil {
		logWriters = append(logWriters, conn)
	}

	multiWriter := zerolog.MultiLevelWriter(logWriters...)

	Log = zerolog.New(multiWriter).
		With().
		Timestamp().
		Caller().
		Logger()

	if err != nil {
		Log.Warn().Err(err).Msg("Logstash tidak terdeteksi di localhost:5000. Log hanya akan tampil di terminal.")
	} else {
		Log.Info().Msg("Sistem Logging Zerolog berhasil tersambung dengan Logstash/ELK Stack!")
	}
}
