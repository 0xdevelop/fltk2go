# API Documentation for fltk2go

## Overview
The `fltk2go` package provides an interface to the FLTK (Fast, Light Toolkit) library, wrapping its functionalities for ease of use in Go applications. This documentation describes the main types, interfaces, and functions available in the package.

## Main Types

### `Element`
- Represents a GUI element in the FLTK framework. 
- Fields:
  - `Label string`: The text label displayed.
  - `Width int`: The width of the element.
  - `Height int`: The height of the element.

### `Window`
- Represents a top-level window in FLTK. 
- Fields:
  - `Title string`: The title of the window.
  - `Elements []Element`: The elements contained in the window.

## Interfaces

### `Widget`
- An interface that all widgets implement.
- Methods:
  - `Draw()`: Method to draw the widget on the screen.
  - `Handle(event Event)`: Method to handle events.

## Functions

### `NewWindow(title string, width int, height int) *Window`
- Creates a new window.
- Parameters:
  - `title`: The title of the window.
  - `width`: The width of the window.
  - `height`: The height of the window.

### `Run() int`
- Starts the FLTK main loop.
- Returns the exit code.

## Example Usage
```go
package main

import "github.com/0xdevelop/fltk2go"

func main() {
    window := fltk2go.NewWindow("My Window", 800, 600)
    window.Run()
}
```

## Conclusion
This documentation provides a brief overview of the main types, interfaces, and functions in the `fltk2go` package. For more information, please refer to the source code and additional documentation available within the repository.