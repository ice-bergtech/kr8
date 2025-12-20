// Package util contains various utility functions for directories and files.
// It includes functions for
// formatting JSON,
// writing to files,
// directory management,
// and go control-flow helpers
package util

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	types "github.com/ice-bergtech/kr8/pkg/types"
)

// Configure zerolog with some defaults and cleanup error formatting.
func SetupLogger(enableColor bool) zerolog.Logger {
	//nolint:exhaustruct
	consoleWriter := zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: !enableColor,
		FormatErrFieldValue: func(err interface{}) string {
			// https://github.com/rs/zerolog/blob/a21d6107dcda23e36bc5cfd00ce8fdbe8f3ddc23/console.go#L21
			colorRed := 31
			colorBold := 1
			s := strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(err.(string), "\\t", " "), "\\n", " | ",
				), "|  |", "|")

			return Colorize(Colorize(s, colorBold, !enableColor), colorRed, !enableColor)
		},
		// Other fields:
		// TimeFormat, TimeLocation, PartsOrder, PartsExclude,
		// FieldsOrder, FieldsExclude, FormatTimestamp, FormatLevel,
		// FormatCaller, FormatMessage, FormatFieldName, FormatFieldValue,
		// FormatErrFieldName, FormatExtra, FormatPrepare
	}

	return log.Output(consoleWriter)
}

// Logs an error and exits the program if the error is not nil.
func FatalErrorCheck(message string, err error, logger zerolog.Logger) {
	if err != nil {
		logger.Fatal().Err(err).Msg(message)
	}
}

// If err != nil, wraps it in a Kr8Error with the message.
func ErrorIfCheck(message string, err error) error {
	if err != nil {
		return types.Kr8Error{Message: message, Value: err}
	}

	return nil
}

// If the error is not nil, log an error and wrap the error in a Kr8Error.
func LogErrorIfCheck(message string, err error, logger zerolog.Logger) error {
	if err != nil {
		logger.Error().Err(err).Msg(message)

		return types.Kr8Error{Message: message, Value: err}
	}

	return nil
}
