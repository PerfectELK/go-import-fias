# import-fias

Import-fias is a golang application for importing fias data into databases (pgsql/mysql) or json

## Installation

Use the git to clone project.

Then
```bash
go build ./cmd/app/main.go
```
Or you may download compiled binary file in latest release


## Usage

You may use .env file or command line arguments

### Cmd args

```
-db-driver=<mysql|pgsql)>
-db-host=<127.0.0.1>
-db-port=<3306>
-db-name=<fias>
-db-user=<fias>
-db-password=<123>
-objects-table=<fias_objects>
-objects-hierarchy-table=<fias_objects_hierarchy>
-object-kladr-table=<fias_object_kladr>
-threads=<1,2,3...>
-archive-path=<if you have fias archive /path/to/fias/archive.zip>
-archive-source=<local if have atchive in ./storage dir, link if you have download link>
-archive-link=<required if -archive-source=link http[s]://link/to/archive.zip>
```

### Env

You may see it in .env.example file

### Usage example

```bash
./main -threads=3 -download=true
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)