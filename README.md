Description
===================
[![GoDoc](https://godoc.org/github.com/bronze1man/kmgStringMutexMap?status.svg)](https://godoc.org/github.com/bronze1man/kmgStringMutexMap)
[![GitHub issues](https://img.shields.io/github/issues/bronze1man/kmgStringMutexMap.svg)](https://github.com/bronze1man/kmgStringMutexMap/issues)
[![GitHub stars](https://img.shields.io/github/stars/bronze1man/kmgStringMutexMap.svg)](https://github.com/bronze1man/kmgStringMutexMap/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/bronze1man/kmgStringMutexMap.svg)](https://github.com/bronze1man/kmgStringMutexMap/network)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/bronze1man/kmgStringMutexMap/blob/master/LICENSE)

a golang data race free map[string]*sync.mutex implement.

Problem it solved
===================
* For example, you write a website with mysql database, If your user finish one step of a three step task, you can use StringMutexMap to lock this user to protect him to finish this step in two different request, and get two reward.
* You should only have one process of that website, or you need some way to distribute the request base on userId.
* Use transaction and Row lock in mysql is much more complex task then this simple locker.

Example
===================
```golang
package main

var gUserIdLockerMap kmgStringMutexMap.StringMutexMap
func processOne(userId string){
    gUserIdLockerMap.LockByString(userId)
    defer gUserIdLockerMap.UnlockByString(userId)
    // SELECT ...
    // REPLACE INTO ...
}
func main(){
    processOne("1")
    processOne("2")
}


```

Notice
===================
* it will free memory after lock and unlock a lot of strings.
* It can not replace mysql transaction
entirely, You should use transaction to avoid your request run half and your process being killed.