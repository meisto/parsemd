// ======================================================================
// Author: meisto
// Creation Date: Sun 26 Feb 2023 11:40:51 PM CET
// Description: -
// ======================================================================
package src

import (
   "log"
   "strings"


	"github.com/muesli/termenv"
)

// Takes two strings representing styles and returns a string representing 
// a combination. On a clash the first string will take precedence. For style
// information such as bold or italic, all info of both elements is combined
func combineStyle(s1 string, s2 string) string {
   if !styleDescRegex.MatchString(s1) || !styleDescRegex.MatchString(s2) {
      log.Fatal("Could not combine styles, one is not a legal representation.")
   }

   var fg, bg, style string

   a := styleDescRegex.FindStringSubmatch(s1)
   b := styleDescRegex.FindStringSubmatch(s2)
   c := styleDescRegex.SubexpNames()

   for i := 0; i < len(c); i++ {
      switch c[i] {
      case "fg":
         if a[i] != "" {
            fg = a[i]
         } else {
            fg = b[i]
         }
      case "bg":
         if a[i] != "" {
            bg = a[i]
         } else {
            bg = b[i]
         }
      case "style":

         // This is a hack because there is no set in go
         styleSet := make(map[string]bool)
         
         for _, x := range(strings.Split(a[i], ",")) {
            styleSet[x] = true
         }
         for _, x := range(strings.Split(b[i], ",")) {
            styleSet[x] = true
         }

         style = ""
         for x := range(styleSet) {style += x + ","}
         style = style[:len(style)-1]
      }
   }

   return "[" + fg + "" + bg + "" + style + "]"
}


func getStyle(s string) func(termenv.Output, termenv.Style) termenv.Style {

   if !styleDescRegex.MatchString(s) {
      log.Fatal("Could not parse string '", s, "' to a valid style.")
   }

   var fg, bg string
   var style []string

   a := styleDescRegex.FindStringSubmatch(s)
   b := styleDescRegex.SubexpNames()
   for i := 0; i < len(a); i++ {
      switch b[i] {
      case "fg":
         fg = a[i]
      case "bg":
         bg = a[i]
      case "style":
         style = strings.Split(a[i], ",")
      }
   }

   return func(op termenv.Output, x termenv.Style) termenv.Style {
      if fg != "" {
         x = x.Foreground(op.Color(fg))
      }

      if bg != "" {
         x = x.Background(op.Color(bg))
      }

      for _, st := range(style) {
         switch st {
            case "bold": x = x.Bold()
            case "italic": x = x.Italic()
            case "": 
            default: log.Print("Captured a default style '", st, "'.")
         }
      }

      return x
   }
}

