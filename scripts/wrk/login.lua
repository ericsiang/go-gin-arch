local accounts = {}

-- 在 init 階段一次性讀入記憶體 
init = function(args)
    local f = io.open("scripts/wrk/login_account.txt", "r")
    if f then
        for line in f:lines() do
            table.insert(accounts, line)
        end
        f:close()
    end
    
    -- 每個 thread 從不同位置開始，避免重複
    local id_num = tonumber(thread and thread.id or 0) or 0
    counter = (id_num % #accounts) + 1
    print(string.format("執行緒 %d 啟動，讀入 %d 筆帳號", id_num, #accounts))
end

wrk.headers["Content-Type"] = "application/json"

request = function()
    local account = accounts[counter]
    local body = string.format([[{"account":"%s", "password":"%s"}]], account, "12345667")
    counter = (counter % #accounts) + 1 -- 循環使用帳號
    
    return wrk.format("POST", "/api/v1/users/login", nil, body)
end