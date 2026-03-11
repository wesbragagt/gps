# Research: tree-cli

**Date:** 2026-03-10T00:00:00Z
**Mode:** answer

## Answer

tree-cli is an npm package that provides a recursive directory listing program producing depth-indented file listings, similar to the Linux `tree` command. Key features:

**Installation**: `npm install -g tree-cli`

**Usage**:
- Command: `tree` or `treee` (to avoid system command conflicts on Windows/Linux)
- Set depth: `tree -l 2` (limits to 2 levels deep)
- Output to file: `tree -o output.txt`
- Can be used as Node module: `require('tree-cli')({...})`

**Popular Alternatives**:
1. **tree-node-cli** - More actively maintained, offers both CLI and Node.js API with options like `-I` to ignore directories (e.g., `tree -I "node_modules"`)
2. **dir-tree-cli** - Supports internationalization and hiding specific directories
3. **treei** - Lightweight alternative with similar functionality
4. **tree-console** - Simple CLI tool usable in both terminal and browser

All tools generate visual tree structures of directories with customizable depth, ignore patterns, and output options.

## Citations

- [tree-cli - npm](https://www.npmjs.com/package/tree-cli)
- [tree-node-cli - GitHub](https://github.com/yangshun/tree-node-cli)
- [tree-node-cli Documentation](http://tree-cli.js.org/)
- [dir-tree-cli - npm](https://www.npmjs.com/package/dir-tree-cli)
- [treei - GitHub](https://github.com/w2xi/treei)
- [tree-console - GitHub](https://github.com/yswrepos/tree-console)
