package main

import "C"

//export my_sum
func my_sum(a, b int) int {
    return (a + b)
}
