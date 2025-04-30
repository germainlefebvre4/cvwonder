# Theme functions

The theme functions are a way to customize the look and add some logic to your CV Theme. You can use them in your theme template to help you display the content of your CV.

## Functions

### `inc` - Increment a number

The `inc` function is used to increment a number. It can also be used to get the next item of a list.

| Code | Output |
|-------|--------|
| `{{ inc 5 }}` | 6 |
| `{{ inc 0 }}` | 1 |
| `{{ inc $length }}` | `$length + 1` |

### `dec` - Decrement a number

The `dec` function is used to decrement a number. It can also be used to get the last item of a list.

| Code | Output |
|-------|--------|
| `{{ dec 5 }}` | 4 |
| `{{ dec 0 }}` | -1 |
| `{{ dec $length }}` | `$length - 1` |

!!! tip "Last item of a list"
    ```html
    {{ $list := (list "a" "b" "c") }}
    {{ index $list (dec (len $list)) }}
    ```

    This code will output `c`.

### `list` - Create a list

The `list` function is used to create a list of items.

| Code | Output |
|-------|--------|
| `{{ list 1 2 3 }}` | [1 2 3] |
| `{{ list "a" "b" "c" }}` | [a b c] |
| `{{ list "CV" "Wonder" }}` | [CV Wonder] |

### `join` - Concatenate strings

The `join` function is used to concatenate a list of strings without or with a separator.

| Code | Output |
|-------|--------|
| `{{ join (list "CV" "Wonder") "" }}` | CVWonder |
| `{{ join (list "CV" "Wonder") " " }}` | CV Wonder |
| `{{ join (list "a" "b" "c") " " }}` | a b c |
| `{{ join (list "1" "2" "3") "-" }}` | 1-2-3 |

### `split` - Split a string

The `split` function is used to split a string into a list of substrings.

| Code | Output |
|-------|--------|
| `{{ split "CV Wonder" " " }}` | [CV Wonder] |
| `{{ split "CV-Wonder" "-" }}` | [CV Wonder] |
| `{{ split "1-2-3" "-" }}` | [1 2 3] |

### `trim` - Remove leading/trailing whitespaces

The `trim` function is used to remove leading and trailing whitespace from a string.

| Code | Output |
|-------|--------|
| `{{ trim " CV Wonder " }}` | CV Wonder |
| `{{ trim "  CV Wonder  " }}` | CV Wonder |
| `{{ trim "   CV Wonder   " }}` | CV Wonder |

### `lower` - String to lowercase

The `lower` function is used to convert a string to lowercase.

| Code | Output |
|-------|--------|
| `{{ lower "CV Wonder" }}` | cv wonder |
| `{{ lower "CVWONDER" }}` | cvwonder |
| `{{ lower "cv wonder" }}` | cv wonder |

### `upper` - String to uppercase

The `upper` function is used to convert a string to uppercase.

| Code | Output |
|-------|--------|
| `{{ upper "CV Wonder" }}` | CV WONDER |
| `{{ upper "cvwonder" }}` | CVWONDER |
| `{{ upper "cv wonder" }}` | CV WONDER |
| `{{ upper "CV WONDER" }}` | CV WONDER |

### `replace` - Replace a substring

The `replace` function is used to replace a substring in a string.

| Code | Output |
|-------|--------|
| `{{ replace "Hello World" "World" "CV Wonder" }}` | Hello CV Wonder |
| `{{ replace "Hello World" "World" "" }}` | Hello  |

### `odd` - Check if a number is odd

The `odd` function is used to check if a number is odd.

| Code | Output |
|-------|--------|
| `{{ odd 1 }}` | true |
| `{{ odd 2 }}` | false |
| `{{ odd 3 }}` | true |

### `even` - Check if a number is even

The `even` function is used to check if a number is even.

| Code | Output |
|-------|--------|
| `{{ even 1 }}` | false |
| `{{ even 2 }}` | true |
| `{{ even 3 }}` | false |

### `add` - Add two numbers

The `add` function is used to add two numbers.

| Code | Output |
|-------|--------|
| `{{ add 1 2 }}` | 3 |
| `{{ add 2 3 }}` | 5 |
| `{{ add 10 2 }}` | 12 |

### `sub` - Subtract two numbers

The `sub` function is used to subtract two numbers.

| Code | Output |
|-------|--------|
| `{{ sub 1 2 }}` | -1 |
| `{{ sub 2 3 }}` | -1 |
| `{{ sub 10 2 }}` | 8 |

## Missing functions

If you need a function that is not available in CV Wonder, you can ask us to add it on [Github Issues](https://github.com/germainlefebvre4/cvwonder/issues/new?template=feature_request.md&title=Theme%20function%20-%20What%20should%20be%20done). We will be happy to help you.
