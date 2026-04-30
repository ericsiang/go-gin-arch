# filename - rbac.rego
package rbac

import future.keywords.contains
import future.keywords.if
import future.keywords.in


# role-permissions assignments
role_permissions := {
  "admins": [
    {"resource": "all", "action": "create"},
    {"resource": "all", "action": "read"},
    {"resource": "all", "action": "edit"},
    {"resource": "all", "action": "delete"},
  ],
}

default allow := false

allow if {
  # lookup the list of roles for the user
  # roles := user_roles[input.user]
  # for each role in that list
  # r := roles[_]
  input.role == "admins"
  # lookup the permissions list for role r
  permissions := role_permissions[input.role]
  # for each permission
  p := permissions[_]

  # check matches 
  p == {"action": input.action, "resource": input.resource}
}

allow {
  input.role == "users"
  input.resource == "user" 
  input.action in ["read","edit"]
}
