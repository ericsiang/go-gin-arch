local counter = 0

init = function(args)
    -- 1. 取得一個保證唯一的種子：
    -- os.time() 提供秒級差異
    -- tostring({}): 利用記憶體地址提供執行緒間的差異
    local seed_str = tostring({}):match("0x%x+")
    local seed = os.time() + tonumber(seed_str, 16)
    math.randomseed(seed)
    
    -- 2. 為了更保險，我們先跳過前幾個隨機數（Lua 的隨機特性）
    math.random(); math.random(); math.random()
    
    -- 3. 產生區段，這次保證不一樣了
    counter = math.random(1, 1000) * 100000
    
    print("執行緒啟動，起始帳號區段: users" .. counter)
end

wrk.headers["Content-Type"] = "application/json"

request = function()
    -- 動態產生 JSON 字串
    -- 使用 user%d 確保帳號不重複
    local body = string.format([[{"account":"users%d", "password":"%s"}]], counter, "12345667")
    counter = counter + 1
    -- 回傳 POST 請求
    return wrk.format("POST", "/api/v1/users/createUser", nil, body)
end