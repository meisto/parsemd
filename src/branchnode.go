// ======================================================================
// Author: Tobias Meisel (meisto)
// Creation Date: Sun 26 Feb 2023 06:18:09 PM CET
// Description: -
// ======================================================================
package src

/** Simple node structure **/
type branchNode struct {
   nodetype string
   content  string
   length   int
   children []node
}

func (self branchNode) Length() int { return self.length }
func (self branchNode) GetChildren() []node { return self.children }
func (self branchNode) GetNodetype() string { return self.nodetype } 
func (self branchNode) GetContent() string { return self.content }


func (self branchNode) applyFirst(
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

/*
   Apply a function too all child nodes and then to this node.
*/
func (self branchNode) applyAll( apply func(node) node) node {
   for i := 0; i < len(self.children); i++ {
      self.children[i] = apply(self.children[i])
   }

   return apply(self)
}


/*
func (n branchNode) Info() string {
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

func (n branchNode) PrintHierarchy() {
   n.printHierarchy(0)
}
func (n branchNode) printHierarchy(level int) {
   fmt.Println(strings.Repeat("   ", level), n.Info())
   for _, x := range(n.children) {
      x.printHierarchy(level + 1)
   }
}
*/
