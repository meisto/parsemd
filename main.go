// ======================================================================
// Author: Tobias Meisel (meisto)
// Creation Date: Sun 26 Feb 2023 03:01:45 PM CET
// Description: -
// ======================================================================
package main

import (
   "fmt"
   "log"
   "os"
   "os/exec"
   "strings"

   "parsemd/src"
)

func display(d src.MDDocument, width int) {
   fmt.Println(strings.Repeat("-", width))

   b := d.Render(width)
   for _, x := range(b) {
      fmt.Println(x)
   }
   fmt.Println(strings.Repeat("-", width))
}

func main() {
       // disable input buffering
    exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
    // do not display entered characters on the screen
    exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
    // restore the echoing state when exiting
    defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()


   content, err := os.ReadFile("test.md")
   if err != nil {
      return
   }

   width := 90

   a := src.GetMDDocument(string(content))

   running := true
   for running {
      display(a, width)

      var x = make([]byte, 3)
      numRead, err := os.Stdin.Read(x)
      if err != nil {
           log.Fatal(err)
      }

      if numRead == 3 { continue }

      switch rune(x[0]) {
         case 'j':
            a.SelectNext()
         case ' ':
            a.ActivateElement()
         case 'q':
            running = false
      }
   }
}

