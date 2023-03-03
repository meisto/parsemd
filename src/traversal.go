// ======================================================================
// Author: meisto
// Creation Date: Fri 03 Mar 2023 01:30:24 AM CET
// Description: This file contains documentation for the mddocument type
// ======================================================================
package src

import (
   "log"
)


func (d *MDDocument) SelectNext() {
   // No focusable elements present
   if !d.hasFocuseable {return }

   var a string

   // Unfocus old node and safe its id
   d.root, _ = d.root.applyFirst(
      func(n node) bool {
         n2, ok := n.(actionLeafNode)
         return ok && n2.IsFocused()
      },
      func (n node) node { 
         n2, ok := n.(actionLeafNode)
         if !ok {
            log.Fatal("Error during passing")
         }

         a = n2.GetId()

         return n2.Unfocus()
      },
   )

   // Focus next node
   passedLast := false
   var changed bool
   d.root, changed = d.root.applyFirst(
      func(n node) bool {
         n2, ok := n.(actionLeafNode)
         if passedLast && ok { return true }

         // We are now past the previously focused element
         if (n2.GetId() == a) { passedLast = true }

         return false
      },
      func (n node) node { 
         n2, ok := n.(actionLeafNode)

         // This condition is fullfilled by the filter function
         if !ok { log.Fatal("Error during passing") }

         return n2.Focus()
      },
   )
   if !changed {
      d.root, changed = d.root.applyFirst(
         func(n node) bool {
            _, ok := n.(actionLeafNode)
            return ok
         },
         func (n node) node { 
            n2, ok := n.(actionLeafNode)
            if !ok { log.Fatal("Error during passing") }
            return n2.Focus()
         },
      )
   }
}

/*
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

*/
