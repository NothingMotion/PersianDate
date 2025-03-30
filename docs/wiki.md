# Welcome
Welcome to PersianDate converter wiki! in this wiki you will learn how to use PersianDate library.


## Table of content
- [Initializing](#initializing)


### Initializing

in order to initialize library and make instance of it you have 2 way.

- using persiandate.New() which creates new instance of it
- using persiandate.Instance() which gives you singleton version of it

1- using persiandate.New()

```go
    import persiandate "github.com/NothingMotion/PersianDate"

    func main(){
        pd:= persiandate.New("FORMAT")
    }
```

2- using persiandate.Instance()

```go
import persiandate "github.com/NothingMotion/PersianDate"

func main(){

    pd := persiandate.Instance("FORMAT")
}
```



both methods works. peak based on project


