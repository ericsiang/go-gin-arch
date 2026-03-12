# filename - rbac.rego
package rbac

import future.keywords.contains
import future.keywords.if
import future.keywords.in

# user-role assignments
# user_roles := {
#     "admin1": ["admin", "dev"],
#     "user1": ["hr"],
# }

# role-permissions assignments
role_permissions := {
  "admin": [
    {"resource": "user", "action": "read"},
    {"resource": "user", "action": "edit"},
  ]
}

default allow := false

allow if {
  # lookup the list of roles for the user
  # roles := user_roles[input.user]
  # for each role in that list
  # r := roles[_]
  input.role == "admin"
  # lookup the permissions list for role r
  permissions := role_permissions[input.role]
  # for each permission
  p := permissions[_]

  # check if the permission granted to r matches the user's request
  p == {"action": input.action, "resource": input.resource}
}


allow {
  input.role == "user"
  action_groups = ["read","edit"]
  input.action in action_groups
}
