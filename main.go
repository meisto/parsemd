// ======================================================================
// Author: Tobias Meisel (meisto)
// Creation Date: Sat 25 Feb 2023 12:06:56 AM CET
// Description: -
// ======================================================================
package main

import (
   "fmt"
   "log"
   "os"
   "regexp"
   "strings"
   "strconv"
   "unicode/utf8"

	"github.com/muesli/termenv"
)

type MDDocument struct {
   base     string
   renderNodes []parseNode
   isMultiselect bool
   index int

   triggerMap map[string]func()
   output termenv.Output
}

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



func (d MDDocument) Render(width int) string {
   res := ""
   
   // Keep line start
   line := ""
   used := 0
   isBlockquote := false

  //  output := d.output

   // blockquoteColor := "#00FF00"

   flush := func() {
      res += line + "\n"
      line = ""
      used = 0

      if isBlockquote {
         line += ">" // output.String(" ").Foreground(output.Color(blockquoteColor)).String()
         line += " "
         used += 2
      }
   }

   for i := 0; i < len(d.renderNodes); i++ {
      a := d.renderNodes[i]

      switch x := a.(type) {
         case LineBreakparseNode:
            isBlockquote = false
            flush()

         case fullLineparseNode:
            isBlockquote = false

            // Flush remaining line if present
            if used > 0 {
               flush()
            }
            res += x.content + "\n"

         case atomarParseNode:
            l := utf8.RuneCountInString(x.content)
            var a string
            a = x.content

            if l > width {
               a = "[PARSEERROR]"
            }

            if used + l >= width {
               flush()
            } 
            line += a
            used += l

         case toggleParseNode:

            var button, label string


            if x.IsFocused() {
               label = x.label // renderer.GenerateNode(x.label, "MarkdownFocusedElement")
            } else {
               label = x.label // renderer.GenerateNoRenderNode(x.label)
            }

            if x.isActivated {
               button = x.button // renderer.GenerateNode(x.button, "MarkdownToggleElementActive")
            } else {
               button = x.button // renderer.GenerateNode(x.button, "MarkdownToggleElement")
            }

            bl := utf8.RuneCountInString(x.button)
            ll := utf8.RuneCountInString(x.label)
            if bl + ll > width {
               log.Fatal("ERROR")
            }

            if used + bl + ll > width {
               flush()
            }

            line += button
            line += label
            used += bl + ll
         case readInputParseNode:
            var field string

            s := fmt.Sprintf("%" + strconv.Itoa(x.maxLength) + "s", x.input)

            if x.IsFocused() {
               field = s // renderer.GenerateNode(s, "MarkdownInputFieldFocused")
            } else {
               field = s // renderer.GenerateNode(s, "MarkdownInputField")
            }

            l := utf8.RuneCountInString(s)
            if l > width {
               log.Fatal("ERROR")
            }

            if used + l > width {
               flush()
            }

            line += field
            used += l

         case actionParseNode:
            var field string // renderer.Renderable
            if x.IsFocused() {
               field = x.label // renderer.GenerateNode(x.label, "MarkdownActionFieldFocused")
            } else {
               field = x.label // renderer.GenerateNode(x.label, "MarkdownActionField")
            }

            l := utf8.RuneCountInString(x.label)
            if used + l > width { 
               flush() }

            line += field
            used += l


         case blockquoteparseNode:
            isBlockquote = true
            if used > 0 { 
               flush() 
            } else {
               line += " "// renderer.GenerateNode(" ", "MarkdownBlockquote"))
               line += " "// renderer.GenerateNoRenderNode(" "))
               used += 2
            }

         case bareTextparseNode:
            s := x.GetContent() 

            if utf8.RuneCountInString(s) > width - used {
               line += s[:width-used]
               s = s[width - used:]
               flush()
            }

            for utf8.RuneCountInString(s) > width {
               s1 := s[0:width - used - 1]
               s  = s[width - used - 1:]

               line += s1
               flush()

            }
            line += s
            used += utf8.RuneCountInString(s)
      }
   }
   res += line
   return res
}



type parseNode interface {
   GetContent()         string
}

type fullLineparseNode struct {
   content  string
   style    string
}
func (n fullLineparseNode) GetContent() string {return n.content}

type atomarParseNode struct {
content  string
style    string
}
func (n atomarParseNode) GetContent() string {return n.content}


type bareTextparseNode struct { content  string }
func (n bareTextparseNode) GetContent() string {return n.content}


type LineBreakparseNode struct {}
func (n LineBreakparseNode) GetContent() string {return ""}

type blockquoteparseNode struct {}
func (n blockquoteparseNode) GetContent() string {return ""}

type focusableParseNode interface {
   GetContent() string
   IsFocused() bool
   Focus() focusableParseNode
   Unfocus() focusableParseNode
}

/** Node to read user input **/
type readInputParseNode struct {
   input string
   id string
   isFocused bool
   maxLength int
}
func (ri readInputParseNode) GetContent() string {return ""}
func (ri readInputParseNode) IsFocused() bool { return ri.isFocused }
func (ri readInputParseNode) Focus() focusableParseNode {
   ri.isFocused = true
   return ri
}
func (ri readInputParseNode) Unfocus() focusableParseNode {
   ri.isFocused = false
   return ri
}
func (ri readInputParseNode) readInput() focusableParseNode {
   ri.input = "123" // ReadLine(true, ri.maxLength)
   return ri
}

// Ensure that the type implements the interface during compiletime
var _ focusableParseNode = (*readInputParseNode)(nil)

/** Button that can be toggled on and off **/
type toggleParseNode struct {
   button   string
   label    string
   id       string
   isActivated bool
   isFocused bool
}
func (tp toggleParseNode) GetContent() string { return tp.button + " " + tp.label }
func (tp toggleParseNode) GetId() string { return tp.id }
func (tp toggleParseNode) IsActive() bool { return tp.isActivated }
func (tp toggleParseNode) Activate() toggleParseNode { 
   tp.isActivated = true
   return tp
}
func (tp toggleParseNode) Deactivate() toggleParseNode { 
   tp.isActivated = false
   return tp
}
func (tp toggleParseNode) Toggle() toggleParseNode { 
   tp.isActivated = !tp.isActivated
   return tp
}
func (tp toggleParseNode) IsFocused() bool { return tp.isFocused }
func (tp toggleParseNode) Focus() focusableParseNode { 
   tp.isFocused = true
   return tp
}
func (tp toggleParseNode) Unfocus() focusableParseNode {
   tp.isFocused = false
   return tp
}

/** Action Parse node **/
type actionParseNode struct {
   label       string
   id          string
   isFocused   bool
}
func (a actionParseNode) GetContent() string { return a.label }
func (a actionParseNode) IsFocused() bool { return a.isFocused }
func (a actionParseNode) Focus() focusableParseNode {
   a.isFocused = true
   return a
}
func (a actionParseNode) Unfocus() focusableParseNode {
   a.isFocused = false
   return a
}
func (a actionParseNode) GetId() string { return a.id }



func parse(line string) []parseNode {
   nodes := []parseNode{}
   lines := strings.Split(line, "\n")

   // Regular expressions used for matching
   headerRegex := regexp.MustCompile(`^\s*#.*$`)
   linebreakRegex := regexp.MustCompile(`^\s*$`)
   unnumberedlistRegex := regexp.MustCompile(`^\s*(-|\*).*$`)
   numberedlistRegex := regexp.MustCompile(`^\s*[0-9]+\..*$`)
   blockquoteRegex := regexp.MustCompile(`^\s*>.*$`)

   line = ""
   for _, element := range(lines) {

      switch {
         case headerRegex.MatchString(element):
               nodes = append(nodes, fullLineparseNode{element, ""})

         case linebreakRegex.MatchString(element):
            _, ok1 := nodes[len(nodes) - 1].(fullLineparseNode) 
            _, ok2 := nodes[len(nodes) - 1].(LineBreakparseNode) 
            _, ok3 := nodes[len(nodes) - 1].(blockquoteparseNode) 
            if !(ok1 || ok2 || ok3) {
               nodes = append(nodes, LineBreakparseNode{})
            }
            nodes = append(nodes, LineBreakparseNode{})

         case unnumberedlistRegex.MatchString(element):
            nodes = append(nodes, fullLineparseNode{element, ""})

         case numberedlistRegex.MatchString(element):
            nodes = append(nodes, fullLineparseNode{element, ""})

         case blockquoteRegex.MatchString(element):
            nodes = append(nodes, blockquoteparseNode{})
            element = regexp.MustCompile(`\s*>\s*`).ReplaceAllString(element, "")
            nodes = append(nodes, bareTextparseNode{element})
            nodes = append(nodes, LineBreakparseNode{})

         default:
            nodes = append(nodes, bareTextparseNode{element})
      }
   }

   nodes2 := []parseNode{}
   wasBareText := false
   for i := 0; i < len(nodes); i++ {
      textNode, isBareText := nodes[i].(bareTextparseNode)

      if i > 0 && wasBareText && isBareText {
         newText := nodes2[len(nodes2) - 1].GetContent() + " " + textNode.GetContent()
         nodes2 = nodes2[:len(nodes2) - 2]
         nodes2 = append(nodes2, bareTextparseNode{newText})



      } else {
         nodes2 = append(nodes2, nodes[i])
      }

      // Update prev
      wasBareText = isBareText

   }


   return nodes2
}


func GetMDDocument(base string) MDDocument {

   parsed := parse(base)

   index := -1
   for i, x := range(parsed) {
      a, ok := x.(focusableParseNode)

      if ok {
         index = i
         parsed[index] = a.Focus()
         break
      }
   }

   return MDDocument{
      base:          base,
      renderNodes:   parsed,
      isMultiselect: false,
      index:         index,
   }
}

func main() {
   content, err := os.ReadFile("test.md")
   if err != nil {
      return
   }

   d := GetMDDocument(string(content))

//   for i := 0; i < len(d.renderNodes); i ++ {
//      fmt.Printf("%T\n", d.renderNodes[i])
//   }

   fmt.Println(d.Render(50))

}
