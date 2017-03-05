local http = require 'resty.http'
local cjson = require 'cjson.safe'

local etcd = {}

function etcd.getUpstreams(appName)
  local httpc = http.new()

  local res, err = httpc:request_uri('http://127.0.0.1:2379/v3alpha/kv/range', {
    method = 'POST',
    headers = {
      ["Content-Type"]= 'application/json',
    },
    body = cjson.encode({
      key = ngx.encode_base64('/upstreams/' .. appName)
    })
  })

  local upstreams = {}

  for index, key in ipairs(cjson.decode(ngx.decode_base64(cjson.decode(res.body).kvs[1].value))) do
    upstreams[index] = '127.0.0.1:' .. key.port
  end

  return upstreams
end

return etcd
