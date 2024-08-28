[![Go version](https://img.shields.io/badge/Go-v1.23-blue)](https://img.shields.io/)
[![CI Status](https://github.com/mdanialr/apilog/workflows/CI/badge.svg)](https://github.com/mdanialr/apilog/actions/workflows/on_push_pr.yml)

# APILog
Write log and 

## Getting Started
```go
// set console/terminal writer to listen to all logs level DEBUG, INFO, WARNING, ERROR
//  console writer, write logs to local terminal/console
cns := apilog.NewConsoleWriter(apilog.DebugLevel)	

// use predefined zap as the backend
wr := apilog.NewZapLogger(cns)
//  or use this if you want to use slog under the hood instead
//    wr := logger.NewSlogLogger(cns)

// call Init before using any other API
wr.Init(3 * time.Second) // you may give longer or shorter timeout/deadline

// info level log message that include contextual data 'hello':'world'
wr.Inf("INFO message", apilog.String("hello", "world"))
//  terminal: 2024-08-28T07:57:14.259+0700    INFO    INFO message    {"hello": "world"}
//  json: {"level":"INFO","time":"2024-08-28T07:59:13+07:00","msg":"INFO message","hello":"world"}

// debug level log message
wr.Dbg("DEBUG message")
//  terminal: 2024-08-28T07:57:14.259+0700    DEBUG   DEBUG message
//  json: {"level":"DEBUG","time":"2024-08-28T07:59:13+07:00","msg":"DEBUG message"}

// SHOULD be called before program exit to make sure any pending logs in buffer properly flushed by each Writer
wr.Flush(2 * time.Second) // you may give longer or shorter timeout/deadline
```

## Contextual Data
```go
// give contextual data that will be passed down to subsequent call
wr = wr.With(log.String("app_env", "local"))
wr.Wrn("warning log")
//  terminal: 2024-08-28T08:04:26.599+0700    WARN    WARNING log     {"app_env": "local"}
//  json: {"level":"WARN","time":"2024-08-28T08:05:13+07:00","msg":"WARNING log","app_env":"local"}

wr = wr.With(log.Num("ram", 2)) // this will also accumulate previous contextual data
wr.Inf("look how many ram i have")
//  terminal: 2024-08-28T08:04:26.599+0700    INFO    look how many ram i have        {"app_env": "local", "ram": 2}
//  json: {"level":"INFO","time":"2024-08-28T08:05:13+07:00","msg":"look how many ram i have","app_env":"local","ram":2}
```

## Leveled Log
```go
cns := log.NewConsoleWriter(log.ErrorLevel)
wr := log.NewZapLogger(cns)
wr.Init(3 * time.Second)

// won't print, less than error
wr.Dbg("DEBUG message")

// won't print, less than error
wr.Inf("INFO message")

// won't print, less than error
wr.Wrn("warning log")

// printed, higher or equal than error
wr.Err("oops!!")
//  terminal: 2023-09-22T13:38:39.784+0700    ERROR   oops!!
//  json: {"level":"ERROR","time":"2024-08-28T08:35:43+07:00","msg":"oops!!"}

// never forget to flush before exiting
wr.Flush(1 * time.Second)
```
Log is prioritized in these order:
1. Error `Err`: (Error) print only in log level Error
2. Warning `Wrn`: (Warning, Error) print in log level Warning, Error
3. Info `Inf`: (Info, Warning, Error) print in log level Info, Warning, Error
4. Debug `Dbg`: (Debug, Info, Warning, Error) print in all log level

## Logger with Context
```go
// put the logger wr to context with 'log.WithCtx'
ctx := log.WithCtx(context.Background(), wr)
printFromCtx(ctx) // pass logger contained context

func printFromCtx(ctx context.Context) {
    // grab logger from context
    wr := log.FromCtx(ctx)
    // note that even if there is no logger inside context
    // this won't cause any panic and just return nop-logger
    // that will never write any logs
	
    wr.Inf("my information")
    //  terminal: 2024-08-28T08:41:56.087+0700    INFO    my information
    //  json: {"level":"INFO","time":"2024-08-28T08:41:24+07:00","msg":"my information"}
}
```
