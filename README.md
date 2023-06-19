# errors

It provides same exact functionality as standard "errors" plus 3 exported functions:

- ```Wrap(description string, err error) error```:
It wraps the errors with the standard fmt.Errorf("... %w", ..) adding a default description if not provided.
- ```WrapLocate(description string, err error) error```:
Same as ```Wrap``` but it locates the position where the error is generated.
- ```NewLocate(description string) error```:
Same as ```New``` but it locates the position where the error is generated.
- In addition a default description is enforced for the standard ```errors.New```

I don't want a fancy package. It can be a drop-in replacement for the standard errors package.
