// ======================================================================
// Author: meisto
// Creation Date: Fri 03 Mar 2023 01:49:59 AM CET
// Description: -
// ======================================================================
package src



/** Node that can perform some action on activation **/
type leafNode struct {
   content  string
   length   int
}

func (self leafNode) Length() int { return self.length }
func (self leafNode) GetContent() string { return self.content }

func (self leafNode) applyFirst(
   match func(node) bool,
   apply func(node) node,
) (node, bool) {
   if match(self) {
      return apply(self), true
   }

   return self, false
}

func (self leafNode) applyAll( apply func(node) node) node {
   return apply(self)
}
