# tracery-golang

This is an implementation of [Kate Compton](https://link.springer.com/chapter/10.1007/978-3-319-27036-4_14)'s tracery, an author-focused generative text tool. The reference implementation is writtin in JavaScript and found [here on github](https://github.com/galaxykate/tracery)

Tracery uses text-expansion: at a high level, it iterates over the string, replacing tokens with strings (that might also contain tokens) until no tokens remain. Additionally, at a high level, it is a pure function; relying only on the initial seed & grammar structure. You can imagine that a single call to it is tracing a ray through the possible probability space: one `trace` is one possible expansion of the grammar.