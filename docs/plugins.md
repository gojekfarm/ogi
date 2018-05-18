
### What is a Golang plug-in?

Go has supported [plugins](https://golang.org/pkg/plugin/) for a while [now since v1.8](https://golang.org/doc/go1.8#plugin).

* These are like compiled Go objects which can be loaded in runtime from a file, path of which can inferred dynamically during runtime as well.

* So they are like dynamically linked libraries but better, these are dynamically loaded libraries. No linking required during compile time.

* They fundamentally work like shared objects in linux, but Go main program don't need to know of all plug-ins it can use during compile time.

> above mentioned aspects are iteration of same feature, it helps your logic flow become truly generic and highly extensible


* So anyone can write a plug-in adhering to it's exposed functionality in use and use it with just making plugin's compiled file available to main at runtime
> * people can manage plug-in with sensitive domain logic internally and use it safely in higher context
> * people can contribute to capabilities of main workflow without worrying of it being merged first

To understand Go plugins separately, Francis Campoy's [simple demo](https://github.com/campoy/golang-plugins) can be looked at. [KrakenD, api gateway](https://github.com/devopsfaith/krakend-ce) uses Go plugins as [here](http://www.krakend.io/blog/krakend-golang-plugins/), if need to look at a real world project.

---

### How can I write a plug-in for Ogi

If you want to write a flow which needs all custom components where it can read from a file, check if each line is an existing path and produce that path mapped to its existence to stdout.

* Consumer plugin would read a file line by line, and pass to Transformer

```

```

* Transformer will check existence of passed path and send this mapping to Producer

```

```

* Producer will print this map of passed path with its existence

```

```

---
