local token = "Bearer your_token_here"

setup = function(thread)
   thread:set("token", token)
end

request = function()
   -- 將 Token 加入 Header
   wrk.headers["Authorization"] = token
   -- 你也可以在這裡動態更換 URL 或 Body
   return wrk.format("GET", "/api/v1/users/(user_id)")
end