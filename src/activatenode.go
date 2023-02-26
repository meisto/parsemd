// ======================================================================
// Author: Tobias Meisel (meisto)
// Creation Date: Sun 26 Feb 2023 06:19:26 PM CET
// Description: -
// ======================================================================
package src

import (
   "fmt"
   "strings"
   "unicode/utf8"
)

/** Node that can perform some action on activation **/
type activatenode struct {
   nodetype string
   content  string
   length   int
   children []node

   // Needed for activatable
   id       string
   group    string
   action   func()
   focused  bool
}

func (self activatenode) Length() int { return self.length }
func (self activatenode) GetChildren() []node { return self.children }
func (self activatenode) GetNodetype() string { return self.nodetype } 
func (self activatenode) GetContent() string { return self.content }
func (self activatenode) GetId() string { return self.id }
func (self activatenode) GetGroup() string { return self.group }

func (n activatenode) Info() string {
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

func (n activatenode) PrintHierarchy() {
   n.printHierarchy(0)
}
func (n activatenode) printHierarchy(level int) {
   fmt.Println(strings.Repeat("   ", level), n.Info())
   for _, x := range(n.children) {
      x.printHierarchy(level + 1)
   }
}

func (self activatenode) applyFirst(
   match func(node) bool,
   apply func(node) node,
) (node, bool) {
   if match(self) {
      return apply(self), true
   }

   for _, x := range(self.children) {
      return x.applyFirst(match, apply)
   }

   return self, false
}

func (self activatenode) Activate() {
   self.action()
}

func (self activatenode) IsFocused() bool { return self.focused}
func (self activatenode) Focus() activateable {
   self.focused = true
   return self
}
func (self activatenode) Unfocus() activateable {
   self.focused = false
   return self
}
