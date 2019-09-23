# Confident

Go configuration with confident!

## Install
```console
go get -u github.com/croz-ltd/confident
```
or as a module:
```
require (
	github.com/croz-ltd/confident v0.0.1
)
```

## What is Confident?

Confident is the configuration solution for short living Go applications like commander application (`kubectl`, `oc`, ...).

It supports JSON and YAML file types for now.

## Why Confident?

Confident is heavily inspired by [github.com/spf13/viper](https://github.com/spf13/viper).

While Viper is designed with long-running Go process execution in mind (like web servers) 
his approach to handling configuration is not suitable for short living commander Go processes.

Confident is developed with short-living Go processes in mind, meaning read at the beginning and persisting 
at the end of execution.

Another main difference with Viper is that Confident unmarshal configuration data into the provided structure. 
All changes to the configuration are performed by modifying structure values and those changes will be persisted
in the configuration file. With this approach, you achieve compile-time verification that the configuration parameter
path exists.

## Usage

First, you need to define the configuration structure:
```go
package config

type Configuration struct {
	Core       Core        `json:"core" yaml:"core"`
	Servers    []Server    `json:"servers,omitempty" yaml:"servers,omitempty"`
}

type Core struct {
	Editor string `json:"editor" yaml:"editor"`
}

type Server struct {
	Name string `json:"name" yaml:"name"`
	Url  string `json:"url" yaml:"url"`
}
```

Next create configuration instance:
> NOTE: Provide configuration default values when creating configuration instance
```go
var configuration = Configuration{
	Core: Core{
		Editor: "vi",
	},
}
```

Next create Confident instance and reference configuration instance for usage:
```go
var k *confident.Confident

func Bootstrap() {
    k = confident.New()
    k.WithConfiguration(&configuration)
     // <Optional>
    k.Name = "config"
    k.Type = "json"
    k.Path = "."
    k.Permission = os.FileMode(0644)
     // </Optional>
    k.Read()
}
```
Create helper method that will cast configuration instance to structure type:
```go
func Get() *Configuration {
	configuration := k.GetConfig()
	if configuration != nil {
		return configuration.(*Configuration)
	}
	return nil
}
```
Modify configuration attributes:
```go
configuration := Get()
configuration.Core.Editor = "vim"
```

Persist changes to file by calling:
```go
func Persist(){
    k.Persist()
}
```

### Typical use case for commander applications

If you write commander application like `kubectl`, `oc` or similar, you can "wrap" your main 
code with Confident initialization before, and Confident Persist at the end:
```go
func main() {
	config.Bootstrap()
	<...magic...>
	config.Persist()
}
```
and in your magic code just call `config.Get()` to obtain configuration and read and/or modify it.
If there is any changes to the configuration, Confident Persist will save it. 