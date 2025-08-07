# Blitra

**An immediate mode rendering framework for elegant terminal UIs in Go**

Blitra gives you precise control over your terminal UI, combining the simplicity of immediate mode rendering with the power of a flexbox-inspired layout system. Build beautiful, responsive terminal applications with a declarative API that feels natural and intuitive.

![Blitra Terminal UI](https://via.placeholder.com/800x400?text=Blitra+Terminal+UI)

## Why Blitra?

Most terminal UI libraries force you to adapt your application to their patterns. Blitra takes a different approach:

- **You control the render loop** - Render frames when and where you want
- **Direct data access** - Use your application data directly, no complex state synchronization
- **Immediate mode rendering** - What you return is what gets rendered, simpler mental model
- **Flexbox-inspired layouts** - Familiar layout concepts that actually work in terminal environments
- **Sophisticated component model** - Build and compose UI elements with a clean, functional API

```go
// This is all you need for a complete, interactive terminal UI
func main() {
  view := blitra.View(viewOpts, func(state blitra.ViewState) any {
    return blitra.Box("container", boxOpts, func(state blitra.BoxState) any {
      if state.Clicked {
        return "You clicked me!"
      }
      return "Click me"
    })
  })
  
  view.Bind()
  defer view.Unbind()
  
  for {
    events, _ := view.RenderFrame()
    // Process events if needed
    time.Sleep(time.Second / 60) // 60 FPS
  }
}
```

## Core Concepts

### Immediate Mode Rendering

Unlike retained mode UI libraries that maintain a persistent widget tree, Blitra rebuilds the UI from your render functions on each frame. This approach:

- Eliminates complex state synchronization
- Makes debugging easier - what you see is exactly what your code returned
- Creates a clear, unidirectional data flow
- Aligns perfectly with Go's simplicity and explicitness

### Layout System

Blitra features a sophisticated four-phase layout algorithm inspired by modern web browsers:

1. **Intrinsic Sizing** (bottom-up) - Elements express their natural size requirements
2. **Available Sizing** (top-down) - Available space is distributed according to constraints
3. **Reflowing** - Adjustments for text wrapping and other bidirectional constraints
4. **Positioning** - Final element position calculation

This allows for expressive layouts with concepts like:

- Horizontal and vertical axes
- Alignment and justification
- Grow and shrink properties
- Margins, padding, and gaps
- Minimum and maximum size constraints

### Component Model

Building UIs in Blitra is based on composable functions that return renderables:

```go
// A custom component is just a function that returns a renderable
func UserCard(user User) any {
  return blitra.Box("user-card", blitra.BoxOpts{
    Border: blitra.RoundBorder(),
    Padding: blitra.P(1),
  }, func(state blitra.BoxState) any {
    return []any{
      blitra.Box("username", blitra.BoxOpts{
        TextColor: blitra.P("#5af"),
      }, func(_ blitra.BoxState) any {
        return user.Name
      }),
      user.Bio,
    }
  })
}
```

## Getting Started

### Installation

```bash
go get github.com/RobertWHurst/blitra
```

Requirements:

- Go 1.23 or later
- Terminal with ANSI escape sequence support
- 256-color or true color terminal for best experience
- Mouse support for interactive elements (most modern terminals support this)
- Compatible terminals: iTerm2, Terminal.app, GNOME Terminal, Konsole, Alacritty, Windows Terminal, etc.

### Basic Example

Here's a complete example showing a centered box with interactive text:

```go
package main

import (
  "github.com/RobertWHurst/blitra"
  "time"
)

func main() {
  // Create and configure the view
  myView := blitra.View(blitra.ViewOpts{
    Align: blitra.P(blitra.CenterAlign),
    Justify: blitra.P(blitra.CenterJustify),
    TextColor: blitra.P("#fff"),
    BackgroundColor: blitra.P("#222"),
    TargetBuffer: blitra.SecondaryBuffer,
  }, func(state blitra.ViewState) any {
    // Return the UI structure
    return blitra.Box("main", blitra.BoxOpts{
      Border: blitra.DoubleBorder(),
      Padding: blitra.P(1),
      BackgroundColor: blitra.P("#333"),
    }, func(boxState blitra.BoxState) any {
      if boxState.Clicked {
        return "You clicked me!"
      }
      return "Click me"
    })
  })

  // Bind to terminal
  myView.Bind()
  defer myView.Unbind()

  // Main loop
  for {
    events, _ := myView.RenderFrame()
    // Process events if needed
    time.Sleep(time.Second / 60)
  }
}
```

## Layout Examples

### Horizontal Layout

```go
blitra.Box("row", blitra.BoxOpts{
  Axis: blitra.P(blitra.HorizontalAxis),
  Gap: blitra.P(1),
}, func(state blitra.BoxState) any {
  return []any{ 
    "Item 1",
    "Item 2",
    "Item 3",
  }
})
```

### Vertical Layout with Alignment

```go
blitra.Box("column", blitra.BoxOpts{
  Axis: blitra.P(blitra.VerticalAxis),
  Align: blitra.P(blitra.CenterAlign),
  Gap: blitra.P(1),
  Border: blitra.LightBorder(),
}, func(state blitra.BoxState) any {
  return []any{
    "Top",
    "Middle",
    "Bottom",
  }
})
```

### Complex Layout with Nested Components

```go
blitra.Box("app", blitra.BoxOpts{
  Axis: blitra.P(blitra.VerticalAxis),
}, func(state blitra.BoxState) any {
  return []any{
    // Header
    blitra.Box("header", blitra.BoxOpts{
      Height: blitra.P(3),
      BackgroundColor: blitra.P("#37c"),
    }, func(state blitra.BoxState) any {
      return "My Application"
    }),
    
    // Main content
    blitra.Box("content", blitra.BoxOpts{
      Axis: blitra.P(blitra.HorizontalAxis),
      Grow: blitra.P(1), // Take remaining space
    }, func(state blitra.BoxState) any {
      return []any{
        // Sidebar
        blitra.Box("sidebar", blitra.BoxOpts{
          Width: blitra.P(20),
          BackgroundColor: blitra.P("#333"),
        }, func(state blitra.BoxState) any {
          return "Navigation"
        }),
        
        // Main panel
        blitra.Box("main", blitra.BoxOpts{
          Grow: blitra.P(1),
          Padding: blitra.P(2),
        }, func(state blitra.BoxState) any {
          return "Content goes here"
        }),
      }
    }),
  }
})
```

## Styling

Blitra offers rich styling options:

| Feature | Options |
|---------|---------|
| Borders | Double, Round, Bold, Light |
| Colors | HEX RGB (`#f00`, `#ff0000`) or named (`red`, `blue`) |
| Margins | Separate values for top, right, bottom, left |
| Padding | Separate values for top, right, bottom, left |
| Alignment | Start, Center, End, Stretch |
| Justification | Start, Center, End, Stretch |
| Text Wrapping | Word, Character, None |

## Comparison with Other Libraries

| Feature | Blitra | Bubble Tea | termui | tview |
|---------|--------|------------|--------|-------|
| **Rendering Model** | Immediate mode | Tea model (Elm-like) | Retained mode | Retained mode |
| **Render Loop Control** | You control it | Framework controlled | Framework controlled | Framework controlled |
| **Data Access** | Direct | Via model updates | Via callbacks | Via callbacks |
| **Layout System** | Flexbox-inspired | Manual positioning | Grid-based | Cell-based |
| **Component Model** | Functional composition | Update/View functions | Widget objects | Widget objects |
| **Learning Curve** | Moderate | Moderate | Moderate | Low-Moderate |

## Core Components

### View

The View is the foundation of every Blitra UI. It serves as the root container that binds to the terminal and manages the rendering lifecycle.

```go
view := blitra.View(blitra.ViewOpts{
  // View configuration options
  Align: blitra.P(blitra.CenterAlign),
  Justify: blitra.P(blitra.CenterJustify),
  BackgroundColor: blitra.P("#222"),
  TargetBuffer: blitra.SecondaryBuffer, // Use alternate screen buffer
}, func(state blitra.ViewState) any {
  // Return content to render
  return "Hello, World"
})
```

**Key Features:**

- **Terminal Binding**: Manages the connection to the terminal, handling capabilities detection and cleanup
- **Buffer Control**: Can render to primary or secondary terminal buffers
- **Event Management**: Captures and returns keyboard and mouse events
- **Layout Root**: Serves as the parent for all other elements
- **Frame Timing**: Provides delta time information for animations
- **Element Querying**: Allows lookups of elements by ID

The `ViewState` object provides:

- Delta time between frames for animation calculations
- The size of the view
- Methods to query elements from the previous frame

### Box

The Box is Blitra's primary layout component, providing structure and organization to your UI. Think of it as similar to a <div> in HTML, but with powerful layout capabilities built in.

```go
blitra.Box("my-box", blitra.BoxOpts{
  // Layout options
  Axis: blitra.P(blitra.VerticalAxis),
  Align: blitra.P(blitra.CenterAlign),
  Justify: blitra.P(blitra.StartJustify),
  Gap: blitra.P(1),
  
  // Sizing options
  Width: blitra.P(40),
  Height: blitra.P(10),
  Grow: blitra.P(1),    // Grow to fill available space
  Shrink: blitra.P(1),  // Shrink if needed
  
  // Style options
  Border: blitra.DoubleBorder(),
  Padding: blitra.P(2),
  BackgroundColor: blitra.P("#333"),
  TextColor: blitra.P("#fff"),
  
  // Text handling
  TextWrap: blitra.P(blitra.WordWrap),
  Ellipsis: blitra.P(true),
}, func(state blitra.BoxState) any {
  // Return content or nested components
  return "Box content"
})
```

**Key Features:**

- **Flexible Layout**: Control how children are arranged using Axis, Align, and Justify
- **Sophisticated Sizing**: Combine fixed sizes with dynamic Grow/Shrink properties
- **Rich Styling**: Borders, colors, padding, and margins
- **Event Awareness**: The BoxState provides information about interaction state
- **Nesting Support**: Boxes can contain other boxes, text, or custom renderables

The `BoxState` object provides:

- Interaction states like Clicked, Hovered (coming soon)
- Access to event information

### Text Handling

Text in Blitra is managed automatically with rich formatting and wrapping capabilities:

```go
// Automatic text wrapping with configurable wrap mode
blitra.Box("text", blitra.BoxOpts{
  Width: blitra.P(40),
  TextWrap: blitra.P(blitra.WordWrap),
  Ellipsis: blitra.P(true),
}, func(state blitra.BoxState) any {
  return "This is a long text that will automatically wrap to fit within the box's width..."
})
```

**Text Features:**

- **Word Wrapping**: Break text at word boundaries when possible
- **Character Wrapping**: Break text at character boundaries when needed
- **No Wrap**: Keep text on a single line with optional ellipsis
- **Truncation**: Automatically truncate with ellipsis when text exceeds available space
- **Styling**: Apply colors and formatting (bold, italic, etc. coming soon)

### Renderable Interface

For advanced cases, you can create custom components by implementing the Renderable interface:

```go
type MyCustomComponent struct {
  id string
  data MyData
}

// Implement the Renderable interface
func (c MyCustomComponent) ID() string {
  return c.id
}

func (c MyCustomComponent) Style() blitra.Style {
  return blitra.Style{
    TextColor: blitra.P("#5af"),
    BackgroundColor: blitra.P("#224"),
  }
}

func (c MyCustomComponent) Render(state blitra.ViewState) any {
  // Return anything that's renderable
  return fmt.Sprintf("Custom component with %v", c.data)
}
```

**Renderable Interface:**

- Gives you complete control over component behavior
- Allows creation of reusable, self-contained components
- Interfaces seamlessly with built-in components

## Advanced Usage

### Querying Elements

Access elements from previous frames:

```go
func(state blitra.ViewState) any {
  counterSize := state.ElementSize("counter")
  // Use size information for layout decisions
}
```

### Event Handling

Process events returned from RenderFrame:

```go
events, _ := view.RenderFrame()
for _, event := range events {
  switch event.Kind {
  case blitra.KeyEvent:
    if event.Key == blitra.KeyEscape {
      // Handle escape key
    }
  case blitra.MouseEvent:
    // Handle mouse event
  }
}
```

## License

Blitra is released under the MIT License. See the [LICENSE](LICENSE) file for details.
