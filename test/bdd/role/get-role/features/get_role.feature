Feature: Get Role By Name
    As an admin
    I want to fetch particular role
    So that I can see the details and manage that particular role

    Background:
        Given The following roles are registered on the system
            | name  | permissions               |
            | role1 | create_user,create_client |

        And I am logged in as admin user
            | email           | password      | role     |
            | admin@gmail.com | adminPassword | get_role |

    @success
    Scenario Outline: Successful get role by name
        When I request to get a role by "<filled_name>"
        Then I should get the following role
            | name   | permissions   |
            | <name> | <permissions> |

        Examples:
            | filled_name | name  | permissions               |
            | role1       | role1 | create_user,create_client |

    @failure
    Scenario Outline: Failed get role by name
        When I request to get a role by "<filled_name>"
        Then my request should fail with "<message>"

        Examples:
            | filled_name | message        |
            | role2       | role not found |
