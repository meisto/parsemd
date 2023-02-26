// ======================================================================
// Author: Tobias Meisel (meisto)
// Creation Date: Sat 25 Feb 2023 12:06:56 AM CET
// Description: -
// ======================================================================
package src

import (
   "log"
   "regexp"
   "strings"
   "unicode/utf8"

	"github.com/muesli/termenv"
   
)


type MDDocument struct {
   base     string
   renderNodes node
   isMultiselect bool
   index int

   triggerMap map[string]func()
   output termenv.Output
}

func (d MDDocument) PrintParseTree() { d.renderNodes.PrintHierarchy() }

/*
func (d *MDDocument) SelectNext() {
   // No focusable elements present
   if d.index == -1 {return }

   // Unfocus current element
   a, ok := d.renderNodes[d.index].(focusableParseNode)
   if ok  {
      d.renderNodes[d.index] = a.Unfocus()
   }

   for i := 1; i < len(d.renderNodes) + 1; i++ {
      a, ok := d.renderNodes[(d.index + i) % len(d.renderNodes)].(focusableParseNode)

      // Check if this element can be focused, if yes focus it and return
      if ok {
         d.index = (d.index + i) % len(d.renderNodes)
         d.renderNodes[d.index] = a.Focus()

         return 
      }
   }
}

func (d *MDDocument) SelectPrev() {
   // No focusable elements present
   if d.index == -1 {return }

   // Unfocus current element
   a, ok := d.renderNodes[d.index].(focusableParseNode)
   if ok  {
      d.renderNodes[d.index] = a.Unfocus()
   }

   for i := 1; i < len(d.renderNodes) + 1; i++ {
      a, ok := d.renderNodes[(d.index - i) % len(d.renderNodes)].(focusableParseNode)

      // Check if this element can be focused, if yes focus it and return
      if ok {
         d.index = (d.index - i) % len(d.renderNodes)
         d.renderNodes[d.index] = a.Focus()

         return 
      }
   }
}

func (d *MDDocument) ActivateElement() {
   switch a := d.renderNodes[d.index].(type) {
      case toggleParseNode:
         if !d.isMultiselect {
            for i := 0; i < len(d.renderNodes); i++ {
               b, ok := d.renderNodes[i].(toggleParseNode)
               if ok {
                  d.renderNodes[i] = b.Deactivate()
               }
            }
         }
         d.renderNodes[d.index] = a.Toggle()

      case readInputParseNode:
         d.renderNodes[d.index] = a.readInput()
         d.SelectNext()

      case actionParseNode:
         id := a.GetId()
         action, ok := d.triggerMap[id]
         if !ok {
            log.Fatal("Could not retrieve action of id '", id, "'.")
         }
         action()
   }
}

func (d *MDDocument) ActivateMultiselect()   { d.isMultiselect = true }
func (d *MDDocument) DeactivateMultiselect() { d.isMultiselect = false }
func (d MDDocument) GetActiveElements() []string {
   res := []string{}
   for i := 0; i < len(d.renderNodes); i++ {
      a, ok := d.renderNodes[i].(toggleParseNode)
      if ok && a.IsActive() {
         res = append(res, a.GetId())

         // Found all elements in singleselect mode 
         if !d.isMultiselect { break }
      }
   }
   return res
}
*/


func (d MDDocument) Render(width int) []string {

   root := d.renderNodes
   res := []string{}

   id := func(x termenv.Style) termenv.Style { return x }

   f := func(
      children []node,
      prefix string,
      prefixStyle func(termenv.Style) termenv.Style,
      followupPrefix string,
      followupPrefixStyle func(termenv.Style) termenv.Style,
      style func(termenv.Style) termenv.Style,
   ) {
      
      res = append( 
         res,
         formatLine(
            d.output,
            children,
            width, 
            prefix,
            prefixStyle,
            followupPrefix,
            followupPrefixStyle,
            style,
         )...
      )
   }

   for _, x := range(root.GetChildren()) {
      switch x.GetNodetype() {
         case "header":

            s := regexp.MustCompile(`^\s*`).ReplaceAllString(x.GetContent(), "") + " "

            f(
               x.GetChildren(),
               s,
               func(x termenv.Style) termenv.Style { return x.Bold().Background(d.output.Color("#009900")).Foreground(d.output.Color("#000000"))},
               strings.Repeat(" ", utf8.RuneCountInString(s) + 3),
               id,
               func(x termenv.Style) termenv.Style { return x.Foreground(d.output.Color("#009900")).Bold()},
            )
         
         case "unnumberedlist", "numberedlist":
            s := regexp.MustCompile(`^\s*`).ReplaceAllString(x.GetContent(), " ")
            
            f(
               x.GetChildren(),
               s,
               func(x termenv.Style) termenv.Style { return x.Bold().Foreground(d.output.Color("#009900"))},
               strings.Repeat(" ", utf8.RuneCountInString(s)+ 3),
               id,
               id,
            )
         case "blockquote":
            s := regexp.MustCompile(`^\s*`).ReplaceAllString(x.GetContent(), " ")

            f(
               x.GetChildren(),
               s,
               func(x termenv.Style) termenv.Style { return x.Bold().Foreground(d.output.Color("#009900"))},
               s + "    ",
               func(x termenv.Style) termenv.Style { return x.Bold().Foreground(d.output.Color("#009900"))},
               id,
            )

         case "line":
            f(
               x.GetChildren(),
               "",
               id,
               "",
               id,
               id,
            )

         case "paragraphbreaker":
            res = append(res, "")

         
         default:
            res = append(res, x.Info())
      }
   }

   return res
}

func formatLine(
   op    termenv.Output,
   nodes []node,
   width int,
   prefix string,
   prefixStyle func(termenv.Style) termenv.Style,
   followUpPrefix string,
   followUpPrefixStyle func(termenv.Style) termenv.Style,
   addStyle func(termenv.Style) termenv.Style,
) []string {
   res := []string{}

   prefixL := utf8.RuneCountInString(prefix)
   followUpPrefixL := utf8.RuneCountInString(followUpPrefix)

   // This would result in an endless loop of recursive calls
   if prefixL >= width || followUpPrefixL >= width { log.Fatal() }

   // Set up first prefix
   l := prefixStyle(op.String(prefix)).String()
   used := prefixL

   // Function to add line to result
   flush := func() {
      res = append(res, l)

      // Set up following prefixes
      l = followUpPrefixStyle(op.String(followUpPrefix)).String()
      used = followUpPrefixL
   }

   for i := 0; i < len(nodes); i++ {
      node := nodes[i]
      t := node.GetNodetype()

      // Catch illegal values
      if t != "text" && t != "bold" && t != "italic" && t != "bolditalic" &&
         t != "linebreak" {
         log.Fatal("Got illegal value:", t)
      }

      if t != "text" {
         // Catch too long elements
         if (node.Length() + prefixL + 3 > width || node.Length() + followUpPrefixL + 3 > width) {
            log.Fatal("Element too large:", node.GetContent())
         }

         // Flush if non-text element wouldnt fit on line
         if used + node.Length() > width {
            flush()
         }
      }

      // Catch trivial elements 
      if t == "linebreak" {
         flush()
         continue
      }


      if t == "text" && used + node.Length() > width {
         
         s := node.GetContent()
         for true {
            i := utf8.RuneCountInString(s)

            ub := width - used

            if width - used > i {
               l += s
               used += i
               break
            } else {
               l += s[:ub]
               s = s[ub:]

               flush()
            }
         }
      } else {
         l += addStyle(styleTextNode(op, node)).String()
         used += node.Length()
      }
   }

   flush()

   return res
}


func styleTextNode(op termenv.Output, n node) termenv.Style {

   switch n.GetNodetype() {
      case "bolditalic":
         return op.String(n.GetContent()).Bold().Italic()
      case "bold":
         return op.String(n.GetContent()).Bold()
      case "italic":
         return op.String(n.GetContent()).Italic()

      default:
         return op.String(n.GetContent())
   }
}





func parse(line string) node {
   lines := strings.Split(line, "\n")

   children := []node{}
   // Parse text line
   for _, element := range(lines) {
      children = append(children, parseLine(element)...)
   }
   rootNode := simplenode{"rootnode", "", 0, sanitize(children)}


   return rootNode
}

func sanitize(nodes []node) []node {
   res := []node{}

   // Merge sequential linenodes into single line nodes
   for i := 0; i < len(nodes); i++ {
      lastEl := len(res) - 1

      if i > 0 && res[lastEl].GetNodetype() == "line" && nodes[i].GetNodetype() == "line" {
         a := append(res[lastEl].GetChildren(), nodes[i].GetChildren()...)
         a = sanitize(a)
         res[lastEl] = simplenode {"line", "", 0, a}

      } else {
         res = append(res, nodes[i])
      }
   }

   // Merge sequential textnodes into single textnodes in each linenode
   nodes = res

   for i := 0; i < len(nodes); i ++ {
      // Skip irrelevant nodes
      if nodes[i].GetNodetype() != "line" {
         continue
      }

      n := nodes[i].GetChildren()
      nNew := []node{}

      for j := 0; j < len(n); j++ {
         lastEl := len(nNew) - 1
         if j > 0 && nNew[lastEl].GetNodetype() == "text" && n[j].GetNodetype() == "text" {
            a := nNew[lastEl]
            b := n[j]

            nNew[lastEl] = simplenode{
               a.GetNodetype(),
               a.GetContent() + " " + b.GetContent(),
               a.Length() + 1 + b.Length(),
               []node{},
            }

         } else {
            nNew = append(nNew, n[j])
         }

      }

      { // Replace node with an updated version
         n := nodes[i]
         nodes[i] = simplenode{n.GetNodetype(), n.GetContent(), n.Length(), nNew}
      }
   }
   return nodes
}

/**
   Parse a complete line of text. This should correspond ony to complete
   line in the source text not in the target text.
**/
func parseLine(element string) []node {
   // Regular expressions used for matching
   headerRegex := regexp.MustCompile(`^(?P<prefix>\s*#+)(?P<content>.*)$`)
   paragraphbreakRegex := regexp.MustCompile(`^\s*$`)
   unnumberedlistRegex := regexp.MustCompile(`^(\s*-)(.*)$`)
   numberedlistRegex := regexp.MustCompile(`^(\s*[0-9]+\.)(.*)$`)
   blockquoteRegex := regexp.MustCompile(`^(\s*>)(.*)$`)

   f := func(reg *regexp.Regexp, nodetype string) []node {
      a := reg.FindStringSubmatch(element)

      x := []node{}
      if a[2] != "" {
         x = parseText(a[2])
      }
      x = sanitize(x)

      return []node{
         simplenode {
            nodetype,
            a[1],
            utf8.RuneCountInString(a[1]),
            x,
         },
      }

   }

   switch {
      case paragraphbreakRegex.MatchString(element):
         return []node{simplenode{"paragraphbreaker", "", 0, []node{}}}

      case headerRegex.MatchString(element):
         return f(headerRegex, "header")

      case unnumberedlistRegex.MatchString(element):
         return f(unnumberedlistRegex, "unnumberedlist")

      case numberedlistRegex.MatchString(element):
         return f(numberedlistRegex, "numberedlist")

      case blockquoteRegex.MatchString(element):
         return f(blockquoteRegex, "blockquote")
         
      default:
         //TODO: here capture individual elements
         return []node{
            simplenode {
               "line",
               "",
               0,
               parseText(element),
            },
         }
   }
}

/**
   Elements whose meaning depends on the whole line context cant be recognized
   here. (E.g. numbered lists)
**/
func parseText(text string) []node {
   // Regexes with a single capture group
   linebreakRegex := regexp.MustCompile(`^(.*)(  )$`)

   regexes := []*regexp.Regexp{linebreakRegex}
   nodetypes := []string{"linebreak"}
   for i := 0; i < len(regexes); i++ {
      regex := regexes[i]
      nodetype := nodetypes[i]

      if regex.MatchString(text) {
         a := regex.FindStringSubmatch(text)
      
         return append(parseText(a[1]), simplenode{nodetype, "", 0, []node{}})
      }
   }

   // Regexes with three capture groups (left, matched, right)
   bolditalicRegex := regexp.MustCompile(`(.*)\*\*\*(.*)\*\*\*(.*)`)
   boldRegex := regexp.MustCompile(`(.*)\*\*(.*)\*\*(.*)`)
   italicRegex := regexp.MustCompile(`(.*)\*(.*)\*(.*)`)

   regexes     = []*regexp.Regexp{bolditalicRegex, boldRegex, italicRegex}
   nodetypes   = []string{"bolditalic", "bold", "italic"}
   for i := 0; i < len(regexes); i ++ { 
      regex := regexes[i]
      nodetype := nodetypes[i]

      if regex.MatchString(text) {
         a := regex.FindStringSubmatch(text)

         res := []node{}
         if a[1] != "" {
            res = append(res, parseText(a[1])...)
         }
         res = append(
            res,
            simplenode{
               nodetype,
               a[2],
               utf8.RuneCountInString(a[2]),
               []node{},
            },
         )
         if a[3] != "" {
            res = append(res, parseText(a[3])...)
         }

         return res
      }
   }

   
   return []node{
      simplenode {
         "text",
         text,
         utf8.RuneCountInString(text),
         []node{},
      },
   }
}

func GetMDDocument(base string) MDDocument {

   parsed := parse(base)

   /*
   index := -1
   for i, x := range(parsed) {
      a, ok := x.(focusableParseNode)

      if ok {
         index = i
         parsed[index] = a.Focus()
         break
      }
   }
   */

   doc := MDDocument{
      base:          base,
      renderNodes:   parsed,
      isMultiselect: false,
      index:         0,
   }

   return doc

}

