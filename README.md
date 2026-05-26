[![GoDoc](https://godoc.org/github.com/eljamo/zxcvbn?status.svg)](https://godoc.org/github.com/eljamo/zxcvbn)

# zxcvbn

This project is a fork of [trustelem/zxcvbn](https://github.com/trustelem/zxcvbn), the Go port of [dropbox/zxcvbn](https://github.com/dropbox/zxcvbn).

It estimates password strength by looking at how real-world password crackers work. Instead of relying only on length or character rules, it checks for common patterns such as frequently used passwords, names and surnames from U.S. Census data, popular English words, dates, repeated characters, sequences, keyboard walks like qwerty, and l33t substitutions. It then uses those matches to give a conservative estimate of how difficult the password would be to guess.

While [trustelem/zxcvbn](https://github.com/trustelem/zxcvbn) focused on being a 1:1 port of [dropbox/zxcvbn](https://github.com/dropbox/zxcvbn), including matching its output, this fork prioritises improvements over strict output parity.

# Modifications

- Improved year detection so recent years are handled dynamically and the regex no longer needs updating for each new decade.
- Improved l33t matching for mixed substitutions, so variants like `p@5$w0rd` are recognized as `password`
