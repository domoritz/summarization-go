More information at https://blog.golang.org/profiling-go-programs

Run with

```sh
go build
./profile
```

Look at cpu profile

```sh
go tool pprof profile cpu.prof
```
Then type `web` to open visualization in browser.

Look at memory profile

```sh
go tool pprof profile mem.prof
```
