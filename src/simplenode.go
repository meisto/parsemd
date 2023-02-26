// ======================================================================
// Author: Tobias Meisel (meisto)
// Creation Date: Sun 26 Feb 2023 06:18:09 PM CET
// Description: -
// ======================================================================
package src

import (
   "fmt"
   "strings"
   "unicode/utf8"

)

/** Simple node structure **/
type simplenode struct {
   nodetype string
   content  string
   length   int
   children []node
}

func (self simplenode) Length() int { return self.length }
func (self simplenode) GetChildren() []node { return self.children }
func (self simplenode) GetNodetype() string { return self.nodetype } 
func (self simplenode) GetContent() string { return self.content }

func (n simplenode) Info() string {
   c := n.content
   if utf8.RuneCountInString(c) > 20 {
      c = c[:18] + "..."
   }

   if c != "" {
      return fmt.Sprintf("{%s: %s}", n.nodetype, c)
   } else {
      return fmt.Sprintf("{%s}", n.nodetype)
   }
}

func (n simplenode) PrintHierarchy() {
   n.printHierarchy(0)
}
func (n simplenode) printHierarchy(level int) {
   fmt.Println(strings.Repeat("   ", level), n.Info())
   for _, x := range(n.children) {
      x.printHierarchy(level + 1)
   }
}

func (self simplenode) applyFirst(
   match func(node) bool,
   apply func(node) node,
) (node, bool) {
   if match(self) {
      return apply(self), true
   }

   for i := 0; i < len(self.children); i++ {
      child, match := self.children[i].applyFirst(match, apply)
      if match {
         self.children[i] = child
         return self, true
      }
   }

   return self, false
}
