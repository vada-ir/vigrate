# Migrate

[![build status](https://git.vada.ir/sdp/migrate/badges/master/build.svg)](https://git.vada.ir/sdp/migrate/commits/master) [![coverage report](https://git.vada.ir/sdp/migrate/badges/master/coverage.svg)](https://git.vada.ir/sdp/migrate/commits/master)

migrate will do migration for databases

## Install
```bash
go install git.vada.ir/sdp/migrate
```

## Create Migration
```bash
migrate create --name=migration1
```

## Migrate Up
```bash
migrate up --schema=schema1
```

## Migrate Rollback
```bash
migrate rollback --schema=schema1 --step=1
```

## Reset Database
```bash
migrate reset --schema=schema1
```

## Migrate Refresh
```bash
migrate refresh --schema=schema1 --step=1
```

## More help
```bash
migrate --help
```
