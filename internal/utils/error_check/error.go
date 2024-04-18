package error_check

import "log/slog"

func CheckError(err error) {
	if err != nil {
		slog.Error("we have problems", err)
		panic(err)
	}
}
