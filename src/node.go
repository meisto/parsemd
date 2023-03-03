// ======================================================================
// Author: Tobias Meisel (meisto)
// Creation Date: Sun 26 Feb 2023 02:58:26 PM CET
// Description: -
// ======================================================================
package src

type node interface {
   // Getter & Setter
   Length() int
   GetContent() string

   // Traversal & Transformation
   applyFirst(func(node) bool, func(node) node) (node, bool)
   applyAll(func(node) node) node
}
