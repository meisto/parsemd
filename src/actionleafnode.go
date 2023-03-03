// ======================================================================
// Author: Tobias Meisel (meisto)
// Creation Date: Sun 26 Feb 2023 06:19:26 PM CET
// Description: -
// ======================================================================
package src



/** Node that can perform some action on activation **/
type actionLeafNode struct {
   nodetype string
   content  func(actionLeafNode) string
   length   int
   children []node

   // Needed for actionLeafNode
   id       string
   group    string
   action   func()
   focused  bool
   evaluate func() string
}

func (self actionLeafNode) Length() int { return self.length }
func (self actionLeafNode) GetChildren() []node { return self.children }
func (self actionLeafNode) GetNodetype() string { return self.nodetype } 
func (self actionLeafNode) GetContent() string { return self.content(self) }
func (self actionLeafNode) GetId() string { return self.id }
func (self actionLeafNode) GetGroup() string { return self.group }

func newActionLeafNode(
   nodetype string,
   content  func(actionLeafNode) string,
   length   int,
   children []node,
   id       string,
   group    string,
   action   func(),
   evaluate func() string,
) actionLeafNode {
   return actionLeafNode {
      nodetype,
      content,
      length,
      children,
      id,
      group,
      action,
      false,
      evaluate,
   }
}

func (self actionLeafNode) applyFirst(
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

func (self actionLeafNode) applyAll(apply func(node) node) node {
   return apply(self)
}


func (self actionLeafNode) Activate() {
   self.action()
}

func (self actionLeafNode) IsFocused() bool { return self.focused}
func (self actionLeafNode) Focus() actionLeafNode {
   self.focused = true
   return self
}
func (self actionLeafNode) Unfocus() actionLeafNode {
   self.focused = false
   return self
}

func (self actionLeafNode) Evaluate() string { return self.evaluate() }
