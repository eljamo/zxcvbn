# Changelog

Entries follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

## Unreleased

- feat(feedback): add feedback package
  - Adds `GetFeedback`, which turns a score and a match sequence into a warning and a list of suggestions. Follows https://github.com/dropbox/zxcvbn/blob/master/src/feedback.coffee with one deliberate divergence: `AllUpper` excludes digits (`^[^a-z\d]+$` vs upstream `^[^a-z]+$`), so digit-containing all-caps tokens such as `PASSWORD123` skip the all-uppercase suggestion. Feedback phrasing may differ from upstream; scoring math keeps upstream semantics.
- feat(zxcvbn): expose feedback on Result
  - `PasswordStrength` now populates a `Feedback` field on `Result`. Additive: existing fields and JSON tags are unchanged.
- fix(matching): correct repeat match indices for multibyte passwords
  - `J` was derived from the byte offset of the first byte of the final rune, so for any non-ASCII repeat `Token` disagreed with `password[I:J+1]` and `lastIndex` could resume mid-rune. `RepeatCount` divided a rune length by a byte length, under-counting repeats of multibyte base tokens.
- fix(matching): pass user inputs into repeat base analysis
  - The repeat matcher's recursive `Omnimatch` call was given `nil` user inputs, so a password built from a repeated user input (`kwyjibokwyjibo` with the user input `kwyjibo`) had its base token scored as bruteforce rather than a rank-1 dictionary hit, materially overestimating its strength.
- fix(scoring): clamp UppercaseVariations to a minimum of 1
  - Defensive only, and unreachable under the current regex gates: any word reaching the loop has at least one ASCII uppercase and one ASCII lowercase letter, so the i=1 term is already >= 2. The clamp stops a future change to those gates from silently collapsing a guess > estimate to 0.
- fix(matching): emit byte indices from the l33t matcher
  - `I`/`J` were rune offsets while the scorer and every other matcher use byte offsets, so a multibyte character before or inside a l33t token landed the match in the wrong DP cell. `PasswordStrength("äp@ssword")` overestimated guesses ~4,000x versus `"xp@ssword"`; both now score comparably.
- fix(matching): match sequences over runes instead of bytes
  - Deltas were computed between successive bytes, so non-ASCII sequences of consecutive code points (e.g. Cyrillic `стуфхц`) were scored as bruteforce instead of cheap sequences, and an emitted `Token` could start or end mid-rune, producing invalid UTF-8. Deltas are now between code points, with `Token`/`I`/`J` still byte-aligned per the library convention.
- fix(scoring): count sequence guesses over runes
  - `SequenceGuesses` read the first byte and multiplied by byte length, overcounting multibyte sequences (`стуфхц` was 26x12 instead of 26x6) and misclassifying the first character of multibyte tokens.
- fix(feedback): choose the repeat warning by base-token length
  - The `Repeats like "aaa" are easy to guess` warning branched on `RepeatCount == 1`, which the repeat matcher can never produce (its regex requires at least one repetition, so the count is always >= 2), making the warning unreachable. Now branches on the base token being a single rune, matching upstream's `base_token.length == 1`.
- perf(matching): lowercase the password once per dictionary-matcher call
  - `strings.ToLower` ran on every candidate substring — O(n) work for each of O(n^2) substrings. The password is now lowered once via an offset table (lowering can change a rune's UTF-8 byte length, so original indices cannot be reused directly).
- perf(scoring): iterate only the sequence lengths present in the DP state
  - The search looped `l` over all n possible lengths with a map lookup each, instead of the handful of lengths actually present (upstream iterates object keys, rarely more than ~5). Present keys are now snapshotted in ascending order, keeping output deterministic. Together with the dictionary-matcher fix this takes a 1,600-character repeated-character password from 79s to ~0.6s; scaling is now the designed O(n^2).
- style(feedback): rewrite the password-rank warning if-else chain as a switch
  - Clears the repo's only golangci-lint (gocritic ifElseChain) issue; no behavior change.
- test(matching): cover multibyte repeats and user inputs
  - Covers a single multibyte base, a multi-rune multibyte base, and adjacent ASCII/multibyte repeats, asserting `Token == password[I:J+1]` in each case. Also covers user inputs reaching the repeat matcher's base analysis.
- test(scoring): cover digit-containing all-caps in UppercaseVariations
- test(zxcvbn): cover user inputs through PasswordStrength, standalone and repeated
- test(feedback): add tests for the feedback package
- test(matching): cover multibyte l33t and sequence matches
  - Multibyte prefixes/suffixes around l33t tokens, Cyrillic and mixed ASCII/multibyte sequences, plus `Token == password[I:J+1]` and UTF-8-validity invariant assertions.
- test(zxcvbn): fuzz match invariants in PasswordStrength
  - Native Go fuzz target asserting, for every emitted match, byte-aligned indices, `Token == password[I:J+1]`, valid UTF-8 tokens, and sane `Guesses`/`Score`. The seed corpus doubles as multibyte regression coverage and runs on every `go test ./...`. Inputs are bounded at 1024 bytes because the match search is O(n^2) by design.
- test(matching): build the trie once in the l33t determinism test
  - The test constructed the matcher without the cached trie, rebuilding the full-dictionary trie on each of its 100 iterations (~5.7s of test time for ~0.07s of value); it now uses `newl33tMatch`. Full matching-package test time drops significantly
- ci(test): run go test across all packages
  - `go test` ran only against the root package, so `matching`, `scoring`, `feedback` and the rest were neither tested nor reported to Codecov. Now runs with `./...`.
- docs(readme): show feedback in the example output
