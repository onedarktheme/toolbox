<h3 align="center">
    <img src="https://raw.githubusercontent.com/onedarktheme/onedark/master/assets/logos/onedark-logo-1544x1544.png" width="100" alt="onedark logo"/></br>
    OneDark Toolbox 
</h3>

<p align="center">
    <a href="https://github.com/onedarktheme/toolbox/stargazers"><img src="https://img.shields.io/github/stars/onedarktheme/toolbox?colorA=282c34&colorB=c678dd&style=for-the-badge"></a>
    <a href="https://github.com/onedarktheme/toolbox/issues"><img src="https://img.shields.io/github/issues/onedarktheme/toolbox?colorA=282c34&colorB=d19a66&style=for-the-badge"></a>
    <a href="https://github.com/onedarktheme/toolbox/contributors"><img src="https://img.shields.io/github/contributors/onedarktheme/toolbox?colorA=282c34&colorB=98c379&style=for-the-badge"></a>
</p>

<p align="center">
    <img src="https://raw.githubusercontent.com/onedarktheme/onedark/master/assets/palette/dark.png" width="400" />
</p>

## OneDark's development tools

A set of tools to work on OneDark.

## Tools

### Palette to Color Swatches

Generate color swatches from a JSON palette. This tool reads a JSON file containing color definitions and outputs 23x23 circle color swatches for each color, saving them as PNG files.
```bash
cat palettes.json | go run cmd/palette-to-color-swatches/main.go
```

You can also specify a file directly:
```bash
go run cmd/palette-to-color-swatches/main.go -file data/palettes.json
```

### Palette to HTML Table

Generate an HTML table showing colors with their hex, RGB, and HSL values from a JSON palette.
```bash
cat data/palettes.json | go run cmd/palette-to-html-table/main.go
```

You can also specify a file directly:
```bash
go run cmd/palette-to-html-table/main.go -file data/palettes.json
```

### 📜 License

OneDark is licensed under the [MIT license](LICENSE).

## 🙏 Acknowledgements

- [Atom One Dark](https://github.com/atom/atom/tree/master/packages/one-dark-ui): the original one dark theme.
- [OneDark.nvim](https://github.com/navarasu/onedark.nvim): the Neovim port that inspired this project.
- [Catppuccin](https://github.com/catppuccin/catppuccin): inspired the structure, design terminology, and documentation style used in this project’s color specification. This project adapts Catppuccin’s design language under the MIT license.