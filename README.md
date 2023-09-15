# xlsxconv

Convert xlsx to csv, json or yaml

## Install

```
go install github.com/johejo/xlsxconv@latest
```

## Usage


```
xlsxconv -sheet=$SHEET_NAME -format=$FORMAT < $SOURCE > $DESTINATION
```

```
Usage of xlsxconv:
  -format string
        Output format ("csv", "json", "yaml") (default "csv")
  -header-row-index int
        Row index to use as header for json and yaml format (defaults to first row)
  -sheet string
        Sheet name (defaults to first sheet)
  -sheet-index int
        Sheet index (defaults to first sheet)
```

## License

MIT
