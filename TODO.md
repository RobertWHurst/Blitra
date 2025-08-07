# Blitra Layout Engine TODO

This file contains prioritized tasks for completing the Blitra layout engine.

## Priority 1: Critical Rendering Components

- [ ] **Complete Text Rendering Implementation**
  - File: `render-text.go`
  - Issue: Function is empty (`renderText` just returns nil)
  - Action: Implement proper text rendering based on element properties
  - Impact: Without this, no text will display at all

- [ ] **Fix Border Rendering**
  - File: `render-container.go` and `render.go`
  - Issue: Border rendering code is either stubbed or commented out
  - Action: Re-enable and complete border rendering implementation
  - Impact: Essential for visual structure of UI components

- [ ] **Fix Element Coordinate System**
  - File: `flow-text.go`
  - Issue: Text positioning is overly simplistic
  - Action: Enhance position calculation to account for alignment
  - Impact: Required for correct visual display of elements

## Priority 2: Core Layout Algorithm

- [ ] **Implement Alignment Logic**
  - File: `flow-container.go`, function: `calcContainerPositionsForChildren`
  - Issue: Marked with TODO comment `// TODO: handle alignment`
  - Action: Complete cross-axis positioning implementation
  - Impact: Critical for flexible layouts (currently elements always positioned at start)

- [ ] **Implement Justification Logic**
  - File: `flow-container.go`, function: `calcContainerPositionsForChildren`
  - Issue: Marked with TODO comment `// TODO: handle justification`
  - Action: Complete main-axis distribution implementation
  - Impact: Essential for proper spacing control between elements

- [ ] **Fix Available Size Calculation**
  - File: `flow-container.go`, function: `calcAvailableContainerSizesForChildren`
  - Issue: Using intrinsic sizes instead of container's available size
  - Action: Properly calculate and propagate available space to children
  - Impact: Critical for correct layout calculations and constraints

## Priority 3: Architectural Consistency

- [ ] **Standardize State Context Pattern**
  - File: `box.go` vs other files
  - Issue: `BoxState` has minimal properties while `ViewState` is more complex
  - Action: Create consistent approach to state propagation
  - Impact: More intuitive API and consistent developer experience

- [ ] **Fix Reflow Logic**
  - File: `flow.go`
  - Issue: Reflow triggered only for text elements
  - Action: Ensure reflow is triggered for all necessary element changes
  - Impact: Improves stability of complex layouts

- [ ] **Implement Style Property Support**
  - File: `flow-container.go`
  - Issue: Incomplete handling of style properties (e.g., Basis not used)
  - Action: Add comprehensive style property support
  - Impact: Enhances flexbox-like functionality and layout control

## Priority 4: Code Quality Improvements

- [ ] **Address Error Handling Inconsistencies**
  - Issue: Mix of error returns and panics
  - Action: Standardize error handling approach
  - Impact: More predictable error handling and better control flow

- [ ] **Improve Documentation**
  - Issue: Inconsistent comment style and documentation
  - Action: Document all public functions and types
  - Impact: Better developer experience and API documentation

- [ ] **Clean Up TODOs and Commented Code**
  - Issue: Numerous TODOs and commented code sections
  - Action: Implement or remove as appropriate
  - Impact: Cleaner, more maintainable codebase

## Additional Consistency Issues

- [ ] **Standardize Visitor Function Naming**
  - Files: `flow.go` and `element.go`
  - Issue: Inconsistent naming between helper functions
  - Action: Adopt consistent naming pattern for visitor functions

- [ ] **Fix Inconsistent Function Capitalization**
  - Issue: Mix of PascalCase and camelCase not following Go conventions
  - Action: Use PascalCase for exported and camelCase for unexported functions

- [ ] **Standardize Pointer Receivers**
  - Files: Throughout codebase
  - Issue: Inconsistent use of pointer vs value receivers
  - Action: Use pointer receivers when methods modify state, value receivers otherwise

- [ ] **Consolidate Helper Functions**
  - Issue: Duplicated null-handling helper functions (`V`, `VOr`, `OrP`)
  - Action: Create a standard set of helper functions used consistently

- [ ] **Extract Common Style Handling**
  - Files: `box.go` and `view.go`
  - Issue: Similar code for style property management
  - Action: Create shared functions for common style operations