local etcd = require 'etcd'

function trimPostfix(str, postfix)
  if endsWith(str, postfix) then
    return string.sub(str, 0, -string.len(postfix) - 1)
  else
    return str
  end
end

function endsWith(str, postfix)
  return postfix == '' or string.sub(str, -string.len(postfix)) == postfix
end

local backends = etcd.getBackends(ngx.var.host)

if backends == nil then
  backends = etcd.getBackends(trimPostfix(ngx.var.host, os.getenv('WILDCARD_DOMAIN')))
end

if backends == nil then
  ngx.say('Upstream not found')
  ngx.exit(404)
else
  ngx.var.target = backends[1]
end
