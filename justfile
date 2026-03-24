run size="two":
  @go build -o go-brc .
  @./go-brc {{size}}

flame:
  @go tool pprof -http 127.0.0.1:8080 cpu_profile.prof
