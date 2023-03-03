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
   root     node
   isMultiselect bool

   hasFocuseable bool

   triggerMap map[string]func()
   output termenv.Output
}

func (d *MDDocument) ActivateElement() {
   d.root, _ = d.root.applyFirst(
      func(n node) bool {
         n2, ok := n.(actionLeafNode)
         return ok && n2.IsFocused()
      },
      func (n node) node { 
         n2, ok := n.(actionLeafNode)
         if !ok { log.Fatal("Error during passing") }
         n2.Activate()
         return n2
      },
   )
}
/*

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

   root, ok := d.root.(branchNode)
   if !ok { return []string{d.root.GetContent()} }

   res := []string{}

   f := func(
      children []node,
      prefix string,
      prefixStyle string,
      followupPrefix string,
      followupPrefixStyle string,
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
         )...
      )
   }




   for _, x := range(root.GetChildren()) {

      // TODO: This is shoddy work and should be fixed later to enable true
      // Tree structures
      x, ok := x.(branchNode)
      if !ok { log.Fatal("Shoddy work crashed.")}


      switch x.GetNodetype() {
         case "header":

            s := regexp.MustCompile(`^\s*`).ReplaceAllString(x.GetContent(), "") + " "

            f(
               x.GetChildren(),
               s,
               "[#000000|#009900|bold]",
               strings.Repeat(" ", utf8.RuneCountInString(s) + 3),
               "[#009900||bold]",
            )
         
         case "unnumberedlist", "numberedlist":
            s := regexp.MustCompile(`^\s*`).ReplaceAllString(x.GetContent(), " ")
            
            f(
               x.GetChildren(),
               s,
               "[#009900||bold]",
               strings.Repeat(" ", utf8.RuneCountInString(s)+ 3),
               "[||]",
            )
         case "blockquote":
            s := regexp.MustCompile(`^\s*`).ReplaceAllString(x.GetContent(), " ")

            f(
               x.GetChildren(),
               s,
               "[#009900||bold]",
               s + "    ",
               "[#009900||bold]",
            )

         case "line":
            f(
               x.GetChildren(),
               "",
               "[||]",
               "",
               "[||]",
            )

         case "emptyline":
            res = append(res, "")
      }
   }

   if true { // Replace color markers with actual color
      op := d.output
      style := getStyle("[||]")
      for i := 0; i < len(res); i++ {
         line := res[i]
         indices := styleDescRegex.FindAllStringIndex(line, -1)

         var newLine string
         if len(indices) == 0 { // No style information
            newLine = style(op, op.String(line)).String()
         } else {
            newLine = style(op, op.String(line[:indices[0][0]])).String()

            for j := 0; j < len(indices); j++ {
               ind := indices[j]
               
               style = getStyle(line[ind[0]:ind[1]])
               var content string
               if j < len(indices) - 1 {
                  content = line[ind[1]:indices[j+1][0]]
               } else {
                  content = line[ind[1]:]
               }
               newLine += style(op, op.String(content)).String()
            }
         }

         res[i] = newLine
      }
   }

   return res
}

func formatLine(
   op    termenv.Output,
   nodes []node,
   width int,
   prefix string,
   prefixStyle string,
   followUpPrefix string,
   followUpPrefixStyle string,
) []string {
   res := []string{}

   { // Organize general structure
      style := "[||]"


      prefixL := utf8.RuneCountInString(prefix)
      followUpPrefixL := utf8.RuneCountInString(followUpPrefix)

      // This would result in an endless loop of recursive calls
      if prefixL >= width || followUpPrefixL >= width { log.Fatal("Too long") }

      // Set up first prefix
      l := prefixStyle + prefix + "[||]"
      used := prefixL

      // Function to add line to result
      flush := func() {
         res = append(res, l)

         // Set up following prefixes
         l = followUpPrefixStyle + followUpPrefix + style
         used = followUpPrefixL
      }

      for i := 0; i < len(nodes); i++ {
         node := nodes[i]

         // Process invisible nodes (style info and linebreaks)
         if node.Length() == 0 {
            if node.GetContent() == "\n" { 
               flush()
            } else {
               style = node.GetContent()
               l += node.GetContent()
            }
            continue
         }


         if used + node.Length() > width {
            
            s := node.GetContent()
            for true {
               i := utf8.RuneCountInString(s)

               ub := width - used

               if width - used > i {
                  l += s // style(op, op.String(s)).String()
                  used += i
                  break
               } else {
                  l += s[:ub] // style(op, op.String(s[:ub])).String()
                  s = s[ub:]

                  flush()
               }
            }
         } else {
            l += node.GetContent()
            used += node.Length()

            // log.Print(">>" + node.GetContent()+ "<< ", node.Length())
         }
      }
      
      flush()
   }
   return res
}




func parse(line string) node {
   lines := strings.Split(line, "\n")

   children := []node{}
   // Parse text line
   for _, element := range(lines) {
      children = append(children, parseLine(element)...)
   }
   rootNode := branchNode{"rootnode", "", 0, sanitize(children)}


   return rootNode
}

func sanitize(nodes []node) []node {
   res := []node{}

   // Merge sequential linenodes into single line nodes
   for i := 0; i < len(nodes); i++ {
      lastEl := len(res) - 1

      predicate := func(x node) bool {
         bn, ok := x.(branchNode)
         return ok && bn.GetNodetype() == "line"
      }

      if i > 0 && predicate(res[lastEl]) && predicate(nodes[i]){
         c1, _ := res[lastEl].(branchNode)
         c2, _ := nodes[i].(branchNode)

         a := append(c1.GetChildren(), leafNode{" ", 1})
         a = append(a, c2.GetChildren()...)
         a = sanitize(a)
         res[lastEl] = branchNode {"line", "", 0, a}

      } else {
         res = append(res, nodes[i])
      }
   }

   return res
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
         branchNode {
            nodetype,
            a[1],
            utf8.RuneCountInString(a[1]),
            x,
         },
      }

   }

   switch {
      case paragraphbreakRegex.MatchString(element):
         return []node{branchNode{"emptyline", "", 0, []node{}}}

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
            branchNode {
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
   for i := 0; i < len(regexes); i++ {
      regex := regexes[i]

      if regex.MatchString(text) {
         a := regex.FindStringSubmatch(text)
      
         return append(parseText(a[1]), leafNode{"\n", 0})
      }
   }

   // Regexes with two capture groups
   styleinfoRegex := regexp.MustCompile(`^(?P<left>.*)(?P<styledesc>` + styleDescRegex.String() + `)(?P<right>.*)$`)
   if styleinfoRegex.MatchString(text) {
      var left, right, content string
      a := styleinfoRegex.FindStringSubmatch(text)
      b := styleinfoRegex.SubexpNames()

      for i := 0; i < len(b); i++ {
         if b[i] == "styledesc" { content = a[i] }
         if b[i] == "left" {left = a[i] }
         if b[i] == "right" {right = a[i] }
      }

      res := []node{}
      if left != "" { res = append(res, parseText(left)...) }
      res = append(res, leafNode{content, 0})
      if right != "" { res = append(res, parseText(right)...) }
   
      return res
   }

   // Regexes with three capture groups (left, matched, right)
   bolditalicRegex := regexp.MustCompile(`(.*)\*\*\*(.*)\*\*\*(.*)`)
   boldRegex := regexp.MustCompile(`(.*)\*\*(.*)\*\*(.*)`)
   italicRegex := regexp.MustCompile(`(.*)\*(.*)\*(.*)`)

   regexes     = []*regexp.Regexp{bolditalicRegex, boldRegex, italicRegex}
   styles      := []string{"[||bold,italic]", "[||bold]", "[||italic]"}

   for i := 0; i < len(regexes); i ++ { 
      regex := regexes[i]

      if regex.MatchString(text) {
         a := regex.FindStringSubmatch(text)

         res := []node{}
         if a[1] != "" {
            res = append(res, parseText(a[1])...)
         }
         res = append(res, leafNode{styles[i], 0})
         res = append(res, leafNode{a[2], utf8.RuneCountInString(a[2])})
         if a[3] != "" {
            res = append(res, parseText(a[3])...)
         }

         return res
      }
   }

   // Special cases
   buttonRegexStr := `(?P<button>[a-zA-Z0-9 ]+)`
   labelRegexStr := `(?P<label>[a-zA-Z0-9 ]+)`
   idRegexStr := `(?P<id>[a-z]([a-z\.]*[a-z])?)`
   groupRegexStr := `(?P<group>[a-z]([a-z\.]*[a-z])?)`

   toggleableRegex := regexp.MustCompile(`(?P<last_block>.*)\<` + 
      buttonRegexStr + `\|` + labelRegexStr + `\|` + idRegexStr + `\|` + 
      groupRegexStr + `\>(?P<next_block>.*)`)

   if toggleableRegex.MatchString(text) {
      a := toggleableRegex.FindStringSubmatch(text)
      b := toggleableRegex.SubexpNames()

      var lastblock, button, label, id, group, nextblock string

      for i := 0; i < len(a); i ++ {
         switch b[i] {
            case "last_block":
               lastblock = a[i]
            case "button":
               button = a[i]
            case "label":
               label = a[i]
            case "id":
               id = a[i]
            case "group":
               group = a[i]
            case "next_block":
               nextblock = a[i]
            default:
               // Don't match
         }
      }

      res := []node{}

      isOn := false

      newNode := newActionLeafNode(
         "toggle", // nodetype
         func(a actionLeafNode) string { // function generating representation
            var buttonstyle, labelstyle string
            if a.IsFocused() {
               labelstyle = "[|#DDFFDD|]"
            } else {
               labelstyle = "[||]"
            }
            if isOn {
               buttonstyle = "[#000000|#FFFFFF|bold]"
            } else {
               buttonstyle = "[#FFFFFF|#000000|]"
            }

            return buttonstyle + button + labelstyle + label
         },
         utf8.RuneCountInString(button) + utf8.RuneCountInString(label), // length
         []node{}, // children
         id,
         group,
         func() {isOn = !isOn}, // An action
         func() string { if isOn {return "on"} else {return "off"}}, // Evaluation function
      )

      // Assemble nodes and return
      if lastblock != "" { res = append(res, parseText(lastblock)...) }
      res = append(res, newNode)
      if nextblock != "" { res = append(res, parseText(nextblock)...) }
      return res
   }


   
   return []node{
      leafNode {
         text,
         utf8.RuneCountInString(text),
      },
   }
}

func GetMDDocument(base string) MDDocument {

   parsed := parse(base)


   var hasActivatableNode bool
   parsed, hasActivatableNode = parsed.applyFirst(
      func(n node) bool {
         _, ok := n.(actionLeafNode)
         return ok
      },
      func (n node) node { 
         n2, ok := n.(actionLeafNode)
         if !ok {
            log.Fatal("Error during passing")
         }

         return n2.Focus()
      },
   )

   doc := MDDocument{
      base:          base,
      root:          parsed,
      isMultiselect: false,
      hasFocuseable: hasActivatableNode,
   }

   return doc

}

