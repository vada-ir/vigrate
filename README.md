# Migrate
migrate will do migration for databases

## Migrate Up
```bash
make migrate-up FOLDER=path/to/directory SCHEMA=my_schema
```

## Migrate Down
```bash
make migrate-down FOLDER=path/to/directory SCHEMA=my_schema
```

## Migrate Down All
```bash
make migrate-down-all FOLDER=path/to/directory SCHEMA=my_schema
```

## Migrate Redo
```bash
make migrate-redo FOLDER=path/to/directory SCHEMA=my_schema
```

## Migrate List
```bash
make migrate-list FOLDER=path/to/directory SCHEMA=my_schema
```