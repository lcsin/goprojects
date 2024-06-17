local key = KEYS[1]
local cntKey = key .. ":cnt"
local inputCode = ARGV[1]

local cnt = tonumber(redis.call("get", cntKey))
local code = redis.call("get", key)



-- 验证次数耗尽
if cnt == nil or cnt <= 0 then
    return -1
end

-- 验证通过
if code == inputCode then
    redis.call("del", cntKey)
    return 0
else
    -- 验证失败
    redis.call("decr", cntKey)
    return -2
end