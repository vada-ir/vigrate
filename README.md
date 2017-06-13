# Vigrate

migration tool and wrapper for sql-migrate written in golang.

## Install
```bash
go get github.com/vada-ir/vigrate
```

## Create Migration
```bash
vigrate create --name=migration1
```

## Migrate Up
```bash
vigrate up --schema=schema1
```

## Migrate Rollback
```bash
vigrate rollback --schema=schema1 --step=1
```

## Reset Database
```bash
vigrate reset --schema=schema1
```

## Migrate Refresh
```bash
vigrate refresh --schema=schema1 --step=1
```

## More help
```bash
vigrate --help
```
