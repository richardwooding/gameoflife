Feature: Game of Life Patterns

  Scenario: Getting the name of a pattern
    Given a pattern named "Glider"
    When I get the name of the pattern
    Then the result should be "Glider"

  Scenario: Stamping a pattern onto a grid
    Given a 3x3 grid
    And a pattern with live cells at (1,1) and (2,2)
    When I stamp the pattern at offset (0,0)
    Then the grid should have live cells at (1,1) and (2,2)

  Scenario: Stamping a pattern with an offset
    Given a 4x4 grid
    And a pattern with live cells at (0,0) and (1,1)
    When I stamp the pattern at offset (1,2)
    Then the grid should have live cells at (1,2) and (2,3)

  Scenario: All predefined patterns are valid
    Given the predefined patterns
    Then each pattern should have a non-empty name
    And each pattern should have valid coordinates

