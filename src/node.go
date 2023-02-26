// ======================================================================
// Author: Tobias Meisel (meisto)
// Creation Date: Sun 26 Feb 2023 02:58:26 PM CET
// Description: -
// ======================================================================
package src

type node interface {
   // Getter & Setter
   Length() int
   GetChildren() []node
   GetNodetype() string
   GetContent() string

   // Traversal & Transformation
   applyFirst(func(node) bool, func(node) node) (node, bool)

   // Info & tostring methods
   Info() string
   PrintHierarchy()
   printHierarchy(int) 
}

type activateable interface {
   GetId() string
   GetGroup() string
   Activate()
   IsFocused() bool
   Focus() activateable
   Unfocus() activateable
}
