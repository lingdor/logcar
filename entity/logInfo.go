package entity

const LogLeveLOff  = '#'
const LogLeveLTrace = '$'
const LogLeveLDebug = '%'
const LogLeveLInfo  = '&'
const LogLeveLWarn  = '\''
const LogLeveLError = '('
const LogLeveLFatal = ')'
const LogLevelAll   = '*'

/*
OFF、FATAL、ERROR、WARN、INFO、DEBUG、TRACE、 ALL
35 # ALL
36 $ TRACE
37 % DEUBG
38 & INFO
39 ' WARN
40 ( ERROR
41 ) FATAL
42 * OFF
*/

type LogLine struct {
	Level  rune
	Line   []byte
	Prefix bool
}
