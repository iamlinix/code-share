package common

import "go.uber.org/zap/zapcore"


var LogLevelMapping = map[string]zapcore.Level {

	"DEBUG":	zapcore.DebugLevel,
	"INFO":		zapcore.InfoLevel,
	"WARN":		zapcore.WarnLevel,
	"ERROR":	zapcore.ErrorLevel,
	"DPANIC":	zapcore.DPanicLevel,
	"PANIC":	zapcore.PanicLevel,
	"FATAL":	zapcore.FatalLevel,
}


/* org level */
const (
	ORG_LEVEL_HEAD		= 0
	ORG_LEVEL_BRANCH	= 1
	ORG_LEVEL_PLANT		= 2
)

/* class level */
const (
	CLASS_LEVEL_NONE	= -1

	CLASS_LEVEL_MAIN	= 0
	CLASS_LEVEL_MID		= 1
	CLASS_LEVEL_SUB		= 2
	CLASS_LEVEL_MATL	= 3
)

const PRICE_ZONE_MAX = -1
const PRICE_ZONE_DEFAULT = "0000"

