# simple text edito

Written in GoLang

## TODO

### system behavior

- [x] To improve multibyte characters
- [ ] To improve display text when the cursor position is over the display size.
- [ ] To show line & column number.
- [ ] To save the text to the file.

## system design

- Separate virtual editor and virtual display.
    - Virtual Editor: containing whole text and cursor position
    - virtual display: having cursor position depending on display size and adjust text to show
