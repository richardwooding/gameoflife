Feature: Colony core logic

  Scenario: A cell with fewer than two live neighbours dies (underpopulation)
    Given a 3x3 colony
    And the cell at (1,1) is alive
    When the next generation is computed
    Then the cell at (1,1) should be dead

  Scenario: A cell with two or three live neighbours lives on
    Given a 3x3 colony
    And the cell at (1,1) is alive
    And the cell at (0,1) is alive
    And the cell at (1,0) is alive
    When the next generation is computed
    Then the cell at (1,1) should be alive

  Scenario: A cell with more than three live neighbours dies (overpopulation)
    Given a 3x3 colony
    And the cell at (1,1) is alive
    And the cell at (0,1) is alive
    And the cell at (1,0) is alive
    And the cell at (2,1) is alive
    And the cell at (1,2) is alive
    When the next generation is computed
    Then the cell at (1,1) should be dead

  Scenario: A dead cell with exactly three live neighbours becomes alive (reproduction)
    Given a 3x3 colony
    And the cell at (0,1) is alive
    And the cell at (1,0) is alive
    And the cell at (1,2) is alive
    When the next generation is computed
    Then the cell at (1,1) should be alive

  Scenario: An empty colony remains empty
    Given a 3x3 colony
    When the next generation is computed
    Then all cells should be dead

  Scenario: A cell at the edge with one neighbour dies
    Given a 3x3 colony
    And the cell at (0,0) is alive
    And the cell at (0,1) is alive
    When the next generation is computed
    Then the cell at (0,0) should be dead

  Scenario: A blinker oscillator pattern oscillates
    Given a 3x3 colony
    And the cell at (1,0) is alive
    And the cell at (1,1) is alive
    And the cell at (1,2) is alive
    When the next generation is computed
    Then the cell at (0,1) should be alive
    And the cell at (1,1) should be alive
    And the cell at (2,1) should be alive
    And the cell at (1,0) should be dead
    And the cell at (1,2) should be dead

  Scenario: A block still life remains unchanged
    Given a 4x4 colony
    And the cell at (1,1) is alive
    And the cell at (1,2) is alive
    And the cell at (2,1) is alive
    And the cell at (2,2) is alive
    When the next generation is computed
    Then the cell at (1,1) should be alive
    And the cell at (1,2) should be alive
    And the cell at (2,1) should be alive
    And the cell at (2,2) should be alive

  Scenario: Toggling a cell changes its state
    Given a 3x3 colony
    When I toggle the cell at (1,1)
    Then the cell at (1,1) should be alive
    When I toggle the cell at (1,1)
    Then the cell at (1,1) should be dead
