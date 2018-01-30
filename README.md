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

Copyright © 2017-2018 wzshiming<[https://github.com/wzshiming](https://github.com/wzshiming)>.

MIT is open-sourced software licensed under the [MIT License](https://opensource.org/licenses/MIT).
