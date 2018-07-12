# Grm 

Grm is a tool like mybatis for the Go programming language.

## Requirements

Go version >= 1.9.

## Download and install

``` shell
go get -u -v gopkg.in/grm.v1/cmd/grm
```

## Commands

``` shell
grm --hl fmt -f ./models/sql/
grm --hl gen go -p models -f ./models/sql/ -o ./models/sql.go
```

or

``` shell
go generate
```


## MIT License

Pouch is licensed under the MIT License. See [LICENSE](https://github.com/go-grm/grm/blob/master/LICENSE) for the full license text.

