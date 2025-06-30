```
             .o8                            
            "888                            
ooo. .oo.    888oooo.   .oooooooo  .ooooo.  
`888P"Y88b   d88' `88b 888' `88b  d88' `88b 
 888   888   888   888 888   888  888   888 
 888   888   888   888 `88bod8P'  888   888 
o888o o888o  `Y8bod8P' `8oooooo.  `Y8bod8P' 
                       d"     YD            
                       "Y88888P'            
```

# nbgo - A Retro Note-Taking System

nbgo is a lightweight, terminal-based note-taking system with a 1980s-inspired TUI, designed for simplicity and speed. It stores notes and bookmarks as Markdown files in a filesystem-based structure, using 84 as the default editor and Glow for viewing.
Features

## Retro TUI

Navigate with up/down keys or mouse, with single-key actions (a, b, e, v, q).
Notes and Bookmarks: Create notes or bookmark URLs as .md files.
Markdown Editing: Uses 84 for editing with Markdown-specific shortcuts.
Viewing: Renders Markdown with Glow.
Notebook Switching: Organize notes in separate notebooks (e.g., work, personal).
Filesystem-Based: Stores notes in ~/.nbgo/<notebook>.

## Usage

Switch to a Notebook

Create or switch to a notebook:

```bash
nbgo use work
```

## Run the TUI

Launch the note-taking interface:

```bash
nbgo
```

Navigate with up/down or mouse.
a: Add a note (opens 84 to edit).
b: Add a bookmark (opens 84 to edit).
v: View selected note/bookmark in Glow.
e: Edit selected note/bookmark in 84.
q or Ctrl+C: Quit.

## Keybindings in 84 (Editor)

F1: Show help modal (lists keybindings).
F2: Save and exit.
F3: Search text (Enter to find, Esc to cancel).
F4: Toggle function key menu.
F5: Toggle Markdown preview.
F10/Esc: Quit without saving.
Ctrl+H: Insert header (# ) or increment level.
Ctrl+L: Insert list item (- ).
Ctrl+B: Insert bold markers (****, cursor between).
Ctrl+I: Insert italic markers (**, cursor between).
Ctrl+K: Insert link template ([text](url), cursor on text).
Ctrl+M: Insert inline code (`code`, cursor inside).
Mouse Click: Set cursor position.

## Notes

### Storage

Notes and bookmarks are saved as .md or .bookmark.md files in ~/.nbgo/<notebook>.

### Terminal Support

Requires a modern terminal (e.g., iTerm2, Terminal.app, Kitty) for mouse and function keys. Test on compact keyboards (e.g., Fn+5 for F5 in 84).

### Limitations

No mouse-based text selection or syntax highlighting in 84 (by design). Search in 84 moves to the first match only.
