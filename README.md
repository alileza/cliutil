# cliutil

github.com/urfave/cli utilities function, to simplify cli app development

# Example

## Generate Flags

```go
var config struct {
		NodeID           string `default:"ulala"`
		LogLevel         string `usage:"just to see logs"`
		ListenAddress    string
		AdvertiseAddress string
		MigrationPath    string
		DatabaseDSN      string
		Test             int64
		SkipMigration    time.Duration
}

app := cli.NewApp()

flags, err := cliutil.GenerateFlags(&config)
if err != nil {
    t.Fatalf("cliutil: %v", err)
}
app.Flags = flags

app.Action = func(ctx *cli.Context) error {
    return cli.ShowAppHelp(ctx)
}
app.Run(os.Args)
```

and you simply got all flags set-up from just passing struct to generate flags

```sh
GLOBAL OPTIONS:
   --node-id value            cliutil: tag usage is not set (default: "ulala") [$NODE_ID]
   --log-level value          just to see logs [$LOG_LEVEL]
   --listen-address value     cliutil: tag usage is not set [$LISTEN_ADDRESS]
   --advertise-address value  cliutil: tag usage is not set [$ADVERTISE_ADDRESS]
   --migration-path value     cliutil: tag usage is not set [$MIGRATION_PATH]
   --database-dsn value       cliutil: tag usage is not set [$DATABASE_DSN]
   --test value               cliutil: tag usage is not set (default: 0) [$TEST]
   --skip-migration value     cliutil: tag usage is not set (default: 1s) [$SKIP_MIGRATION]
   --help, -h                 show help
   --version, -v              print the version
```

### Available struct tag

| Tag      |      Description          | 
|----------|:-------------------------:|
| flag     | Flag name                 |
| env      | Environment variable name |
| usage    | Set usage description     |
| default  | Set default flag value    |
| hidden   | Hide flag from help       |