# Backend Engineering Challenge

Only requisite: Docker

### Input

Place your test input in the "`events.json`" file. If you wish to use another file name, use the appropriate flag to do so (explained next).

### Usage

Edit the "`Dockerfile`" to call the application with your desired parameters. The application usage can be seen by calling "`./unbabel -h`" and the output of that command is as follows:

```
A simple command line application that parses a stream of events and produces an aggregated output. In this case, we're interested in calculating, for every minute, a moving average of the translation delivery time for the last X 
minutes.

Usage:
  unbabel [flags]

Flags:
      --client_name string       Specify a specific client name (default "all")
  -d, --debug                    Display debugging output in the console (default: false)
  -h, --help                     help for unbabel
      --input_file string        Input file (default "events.json")
      --source_language string   Specify a specific source language (default "all")
      --target_language string   Specify a specific target language (default "all")
      --window_size int          Specify a window size (default 10)
```

Finally, you can call a script to build and run the application in a Docker environment:

```
sh build_and_run.sh
```
