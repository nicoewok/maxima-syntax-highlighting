## Generate language.json

Since the list of built-in functions and variables is so long, I wrote a script to automate the generation of the ```mac.tmLanguage.json```.

You have to provide the ```functions.txt```, ```variables.txt``` and ```constants.txt``` files and execute the ```main.go``` by using:

```bash
go run main.go
```