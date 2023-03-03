// ======================================================================
// Author: meisto
// Creation Date: Mon 27 Feb 2023 12:04:56 AM CET
// Description: -
// ======================================================================
package src

import "regexp"

var hexaRegString string = `(#[0-9A-Fa-f]{6,6})`
var styleRegString = `((italic)|(bold))`
var styleDescRegex *regexp.Regexp = regexp.MustCompile(
   `\[(?P<fg>` + hexaRegString + `?)\|(?P<bg>` + hexaRegString + 
   `?)\|(?P<style>(` + styleRegString + `(,` + styleRegString + `)*)?)\]`)

// var styleDescRegex *regexp.Regexp = regexp.MustCompile(
//    `^\[(?P<fg>` + hexaRegString + `?)\|(?P<bg>` + hexaRegString + 
//    `?)\|(?P<style>(` + styleRegString + `(,` + styleRegString + `)*)?)$` )

