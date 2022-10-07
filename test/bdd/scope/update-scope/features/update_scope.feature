Feature: Update Scope

    As an Admin,
    I want to update the scope's information
    So that  scope will have up-to-date information.

    Background:
        Given I am logged in as admin user
            | email           | password      | role         |
            | admin@gmail.com | adminPassword | update_scope |

    @success
    Scenario: Successful scope update
        Given there is scope with the following details
            | name   | description       | resource_server_name |
            | openid | your profile info | sso                  |
        When I fill the following details
            | description     |
            | new description |

        And I update scope
        Then The scope should be updated

    @failure
    Scenario: Failure scope update
        Given there is scope with the following details
            | name   | description       | resource_server_name |
            | openid | your profile info | sso                  |
        When I fill the following details
            | description   |
            | <description> |

        And I update scope
        Then The update should fail with field error description "<message>"

        Examples:
            | description | message                 |
            |             | description is required |

    @failure
    Scenario: Scope not found
        Given there is scope "<name>"
        And I fill the following details
            | description   |
            | <description> |

        When I update scope
        Then The update should fail with message "<message>"

        Examples:
            | name     | description     | message         |
            | emailing | new description | scope not found |
