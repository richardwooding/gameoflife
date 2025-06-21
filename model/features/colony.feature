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

