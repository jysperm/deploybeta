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

local upstreams = etcd.getUpstreams(trimPostfix(ngx.var.host, '.deploybeta.site'))

ngx.var.target = upstreams[1]
