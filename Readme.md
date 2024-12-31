# Blitra

A immediate mode renderer for the terminal - Provides a framework for creating
beautiful terminal applications with simplicity and ease.

## Overview

A point of emphasis in Blitra's design is allowing the developer to control when
frames are drawn, and in which go routine to do so. A lot of other TUI
frameworks impose a render loop, which adds complexity around synchronization;
just getting data in becomes a hassle. The idea with Blitra is that you can use
your data directly, without needing to pass it off to a render loop.

Below is an example of a blitra view that renders a bordered box in the center
of the terminal with the text "Hello World!" in the center. If clicked the text
will change to "You're clicking the message!" until released.

```go
package main

import (
  "github.com/RobertWHurst/blitra"
)

func main() {
  helloView := blitra.View(helloViewOptions, func(viewState) {

    messageDiv := blitra.Box(messageBoxOptions, func(messageBoxState) {
      if (messageBoxState.Clicked) {
        return "You're clicking the message!"
      }
      return "Hello World!"
    })

    return messageDiv
  })

  helloView.Bind()
  defer helloView.Release()

  // in your own render loop...
  helloView.RenderFrame()
}

var helloViewOptions = blitra.ViewOpts{
  Align: blitra.CenterAlign,
  Justify: blitra.CenterJustify,
  TextColor: '#fff',
  BackgroundColor: '#000',
}

var messageBoxOptions = blitra.DivisionOptions{
  Border: blitra.BorderDouble,
  Padding: 1,
}

```

An example of the output:

```
╭──────────────────────────────────────────────────────────────────────────────╮
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                             ╔══════════════╗                                 │
│                             ║              ║                                 │
│                             ║ Hello World! ║                                 │
│                             ║              ║                                 │
│                             ╚══════════════╝                                 │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
╰══════════════════════════════════════════════════════════════════════════════╯
```

And when clicked:

```
╭──────────────────────────────────────────────────────────────────────────────╮
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                     ╔══════════════════════════════╗                         │
│                     ║                              ║                         │
│                     ║ You're clicking the message! ║                         │
│                     ║                              ║                         │
│                     ╚══════════════════════════════╝                         │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
╰══════════════════════════════════════════════════════════════════════════════╯
```

## How Blitra Works

Using Blitra begins with creating a view. To do so the View function is called,
with view options - things like layout, and styling. Views are the root of all
UIs created with Blitra, and are intented to be used on demand.

- View - Takes a function which is expected to return:
  - string
  - struct implementing Node
  - nil
  - a slice of any containing a mix of the above

- Renderable - Any struct that has a render method which returns:
  - one or a slice of:
    - string
    - struct implementing Node
    - nil
    - a slice of any containing a mix of the above
  - This is most often a built in component like a box.

- Node - Any struct which implements the needed methods to be wrapped by an
  element. These methods are:
  - ???layout related???

- Element - A wrapper for a node; created by the view to wrap each node created
  during the render process. Contains values calculated while traversing up and
  down the tree during the render process. Element contains the following info:
  - Node - The node the element wraps
  - Children - A slice of child elements
  - ???layout related???
