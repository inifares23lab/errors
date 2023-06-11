# errors

It provides same exact functionality as standard "errors" plus 2 functions:

- Wrap(description string, err error) error:
It wraps the errors with the standard fmt.Errorf("... %w", ..) in an opinionated way.
It also locates the error if it is the first in the chain or in the absence of a descriptive message.
- Locate(err error) error:
Locates the errors more flexibly.

The assumptions are that I want to know what and when it happened but without adding much overhead hence the conditional location functionality and the "sligtly enforced" description.

- The Join functionality is also unchanged but the output is more readable with the formatting in use here
