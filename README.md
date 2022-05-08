
**Build & Environment preparation instructions**
```
make  # will produce srv binary

make env       # to start clickhouse in docker
make env_stop  # to shutdown clickhouse
./test.sh      # perform testing
```

**Usage**
```
./srv
```

**Notices**
I was short of time, so got a lot of hardcoded values, like usernames, passwords, some string constants.
In normal development process I try to move them to cfg files or some constant lists, but I was told
that I need to spent only 2-3 hours for this task, so it was a kind of tradeoff. If there were any questions, please feel free to ask them.