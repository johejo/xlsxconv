# xlsxconv

Convert xlsx to csv, json or yaml

## Install

```
go install github.com/johejo/xlsxconv@latest
```

## Usage


```
xlsx2csv -sheet=$SHEET_NAME -format=$FORMAT < $SOURCE > $DESTINATION
```

```
Usage of xlsxconv:
  -format string
        Output format ("csv", "json", "yaml") (default "csv")
  -sheet string
        Sheet name (defaults to first sheet)
```

## License

MIT
