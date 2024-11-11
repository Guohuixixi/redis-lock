local val  = redis.call("get",KEYS[1])
--在加锁重试的时候，判断自己上一次成功or失败
if val==false then
    --key不存在
    return redis.call("set",KEYS[1],ARGV[1],"EX",ARGV[2])
elseif val==ARGV[1] then
    --刷新过期时间
    redis.call("expire",KEYS[1],ARGV[2])
    return "OK"
else
    --此时其他人持有锁
    return ""
end