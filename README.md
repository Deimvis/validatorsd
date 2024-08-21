# Validator Simplified package

Validator Simplified is an extension for [validator v10](https://pkg.go.dev/github.com/go-playground/validator/v10) package that allows you to avoid the registration step when adding custom validation to your struct. Instead of writing and registering a separate function, you need to implement a `ValidateSelf() error` method in your struct and call the generic `Validate` function which will validate the struct using both validation tags and `ValidateSelf` methods.

## Installation

```bash
go get github.com/Deimvis/validatorsd
```

## Quick Start

1. Define your models (structs)

    ```go
    // models.go
    package hi

    import "errors"

    type MyStruct struct {
        Answer int
    }

    func (s *MyStruct) ValidateSelf() error {
        if s.Answer != 42 {
            return errors.New("answer is wrong")
        }
        return nil
    }
    ```

2. (optionally) Create a validation shortcut

    ```go
    // validation.go
    package hi

    import (
        "github.com/Deimvis/validatorsd"
        "github.com/go-playground/validator/v10"
    )

    var val = validator.New(validator.WithRequiredStructEnabled())

    func Validate(obj interface{}) error {
        return validatorsd.Validate(val, obj)
    }
    ```

3. Validate your model

    ```go
    // main.go
    package hi
    
    import "fmt"
    
    func main() {
        s := MyStruct{
            Answer: 99,
        }
        err := Validate(s)
        fmt.Println(err) // answer is wrong
    }
    ```

## Details

* `Validate` function validates using both validation tags and `ValidateSelf` methods on structs and substructs
* `Validate` panics when the given object is nil
* `Validate` goes recursively through substructs, including embedded structs
* If a struct/substruct doesn't implement a `ValidateSelf` method, nothing will happen and it will still be recursively traversed by `Validate` function
* It doesn't matter whether `ValidateSelf` method has a value or a pointer receiver — `Validate` function will find it either way
