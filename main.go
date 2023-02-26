// ======================================================================
// Author: Tobias Meisel (meisto)
// Creation Date: Sun 26 Feb 2023 03:01:45 PM CET
// Description: -
// ======================================================================
package main

import (
   "fmt"
   "os"
   "strings"

   "parsemd/src"
)

func main() {
   content, err := os.ReadFile("test.md")
   if err != nil {
      return
   }

   width := 90

   a := src.GetMDDocument(string(content))
   b := a.Render(width)
   
   if true {
      fmt.Println(strings.Repeat("-", width))
      for _, x := range(b) {
         fmt.Println(x)
      }
      fmt.Println(strings.Repeat("-", width))
   } else {
      a.PrintParseTree()
   }


//   for i := 0; i < len(d.renderNodes); i ++ {
//      fmt.Printf("%T\n", d.renderNodes[i])
//   }

   // fmt.Println(d.Render(50))

}

