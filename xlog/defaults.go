// Copyright (c) 2019, 2021, Geert JM Vanderkelen

package xlog

// defaultLogLevel defines the default log level.
var defaultLogLevel = InfoLevel

// defaultTextFormat
var defaultTextTemplate = "[{{ .Timestamp }}] {{ .Scope }} {{ .Level }} {{ .Message }}"
