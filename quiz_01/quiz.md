# How to build

```shell
cd {project root}/quiz_01
go build -o quiz .
./quiz --help
```

# Flags

| Flag  | Description                                   | default      | type   |
|-------|-----------------------------------------------|--------------|--------|
| csv   | a csv file in the format of 'question,answer' | problems.csv | string |
| limit | the time limit for the quiz in seconds        | 30           | int64  |

# Execution

```shell
./quiz  # default flag
./quiz -csv=test.csv  # default limit 
./quiz -limit=10  # default csv
./quiz -csv=test.csv -limit=120  # No Default start
```


