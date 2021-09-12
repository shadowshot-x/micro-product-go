# Golang Project Best Practices

## Dependency Injection :-
In simple words, we want our functions and packages to receive the objects they depend on ie. We dont want to declare new instances inside packages to have control over them. For Eg :- Using Structs to declare the methods and variables and the methods are passed on the value of logger.
