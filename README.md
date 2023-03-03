# parsemd
A simple library to write small interfaces using a markdown-like language

## Markdown-like language

**Supported:**
- Italic (\*sth\*)
- Bold (\*\*sth\*\*)
- Bold and italic (\*\*\*sth\*\*\*)
- Headers (\# sth, \#\# sth, ...)
- Ordered List (1. sth, 2.sth, ...)
- Unordered List (- sth, -sth, ...)
- Linebreak ("  \\n")
- Blockquote

**Halfway there:**
- Toggleable node (\<(button|label|id|group\>

button = [a-zA-Z0-9 ]+  
label = [a-zA-Z0-9 ]+  
id = \[a-z\]+\([a-z_]*\[a-z\]\)?  
group = [a-z]\([a-z\.]*[a-z]\)?  

**To be added/Not tested:**
- Singleselect node (\<label|id|length\>)
- Multiselect node
- Text input node
- Action node (\<!label|id\>)

**Maybe?**
- Tables

## Known bugs
- Markdown: double space is interpreted as linebreak
