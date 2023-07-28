# errors

It provides same functionality as standard "errors" with some extra sugar:

- ```Wrap(description string, err error) error```:
It wraps the errors similarly tothe standard fmt.Errorf("... %w", ..).
- ```TWrap(description string, err error) error```:
Same as ```Wrap``` but it locates the position where the error is generated.
- ```TNew(description string) error```:
Same as the standard ```errors.New``` but it locates the position where the error is generated.
- ```Error() string``` "instance method"(not really correct wording for golang) returns the string representation only the outermost(the one generated last in the chain) error including its location if not empty.
- ```String() string``` returns the string represetation of the whole chain of errors.
- ```Stack() interface{}``` returns the error stack trace and it is compatible with **zerolog**
- All other functions cover the standard ```errors``` functionality and behave virtually the same.
