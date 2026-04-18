// Package util contains various utility functions for directories and files.
// It includes functions for
// formatting JSON,
// writing to files,
// directory management,
// and go control-flow helpers
package util

import (
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	types "github.com/ice-bergtech/kr8/pkg/types"
)

// Configure zerolog with:
//   - Sensible defaults
//   - Cleanup error formatting
//   - Flatten <anonymous> frames.
func SetupLogger(enableColor bool) zerolog.Logger {
	//nolint:exhaustruct
	consoleWriter := zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: !enableColor,
		FormatErrFieldValue: func(err any) string {
			// https://github.com/rs/zerolog/blob/a21d6107dcda23e36bc5cfd00ce8fdbe8f3ddc23/console.go#L21
			colorRed := 31
			colorBold := 1

			// Replace escaped tabs/newlines with spaces/pipes for easier splitting
			pipedErrString := strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(err.(string), "\\t", " "), "\\n", " | ",
				), "|  |", "|")

			// Split the stack trace into frames to format and flatten anonymous frames
			frames := strings.Split(pipedErrString, "|")
			var resultFrames []string
			var anonFrameCount int
			anonFrameCount = 0

			for _, frame := range frames {
				frame = strings.TrimSpace(frame)
				if frame == "" {
					continue
				}
				// Flatten consecutive or duplicate <anonymous> function frames
				if strings.Contains(frame, "<function <anonymous>>") {
					anonFrameCount++

					continue
				} else {
					if anonFrameCount > 0 {
						resultFrames = append(resultFrames, "[anonymous func x"+strconv.Itoa(anonFrameCount)+"]")
					}
					anonFrameCount = 0
				}
				resultFrames = append(resultFrames, frame)
			}
			// Log if last frame is an anonymous one
			if anonFrameCount > 0 {
				resultFrames = append(resultFrames, "[anonymous func x"+strconv.Itoa(anonFrameCount)+"]")
			}
			// Join frames with newlines and indentation
			formatted := strings.Join(resultFrames, "\n  ")

			// Colorize as before
			return Colorize(Colorize(formatted, colorBold, !enableColor), colorRed, !enableColor)
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
