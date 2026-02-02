# **[golang-private-composer-check](https://github.com/KaiserSohn/golang-private-composer-check)**

## Description

This application allows you to quickly check for the presence of packages in a private packagist repository.

## Requirements
Requirement: Go 1.24+
```bash  
go build
```
## Usage
```bash  
./PrivateComposerChecker --lock-file=composer.lock --auth=auth.json
```
## Output

```
symfony/form:v6.4.31 Status: 404 Available versions (v6.4.30, v6.4.24, v6.4.21, v6.4.20, )
symfony/translation:v6.4.31 Status: 404 Available versions (v6.4.24, v6.4.23, v6.4.21, v6.4.19, )
symfony/twig-bridge:v6.4.31 Status: 404 Available versions (v6.4.30, v6.4.25, v6.4.24, v6.4.22, )
symfony/twig-bundle:v6.4.31 Status: 404 Available versions (v6.4.24, v6.4.23, v6.4.13, v6.4.12, )
```
OR
```
✅ All packages is available
```