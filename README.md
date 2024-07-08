# `tabp` - A Lisp like language that operates on Table instead of list.

Tabp is a programming language inspired by Lisp and Lua. It borrows its syntax
from Lisp but uses table (similar to Lua tables) as it's only datastructure.
It is called Tabp because unlike Lisp (List Processor), it operates on tables.

Table is a datastructure that acts as a map and a vector/slice at the same time.
All entries are stored in map except those that are part of the **sequence**.
Table sequence define entries with an integer key `i` in range `0 to n` (exclusive)
where `Table.Get(i) != nil` and `Table.Get(n) == nil`.

## Contributing

If you want to contribute to `tabp` to add a feature or improve the code contact
me at [alexandre@negrel.dev](mailto:alexandre@negrel.dev), open an
[issue](https://github.com/negrel/tabp/issues) or make a
[pull request](https://github.com/negrel/tabp/pulls).

## :stars: Show your support

Please give a :star: if this project helped you!

[![buy me a coffee](https://github.com/negrel/.github/blob/master/.github/images/bmc-button.png?raw=true)](https://www.buymeacoffee.com/negrel)

## :scroll: License

MIT © [Alexandre Negrel](https://www.negrel.dev/)
