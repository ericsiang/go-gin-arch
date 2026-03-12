package rbac


test_allow_admin_read {
  allow with input as {"role": "admin", "action": "read","resource":"user"}
}

test_allow_admin_edit {
  allow with input as {"role": "admin", "action": "edit","resource":"user"}
}

test_allow_user_read {
  allow with input as {"role": "user", "action": "read"}
}

test_allow_user_edit {
  allow with input as {"role": "user", "action": "edit"}
}

test_deny_user_read {
  not allow with input as {"role": "user", "action": "add"}
}

test_deny_role {
  not allow with input as {"role": "test", "action": "add"}
}